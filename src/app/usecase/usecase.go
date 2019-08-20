package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	app "work-relic-bot/src/app/repository"
	"work-relic-bot/src/models"
)

type IUsecase interface {
	Check(ctx context.Context) (string, error)
	Webhook(ctx context.Context, m *models.Webhook) (string, error)
	handleInsert(text string, username string) (string, error)
}

type appUsecase struct {
	appRepo        app.IRepository
	contextTimeout time.Duration
	botURL         string
}

const timeFormat = "2006-01-02 15:04"

func NewUsecase(a app.IRepository, timeout time.Duration) IUsecase {
	credential, err := os.Open("./src/config/credentials.json")
	if err != nil {
		panic(fmt.Sprintf("Error Opening Credentials %s", err))
	}
	log.Println("Success to open credentials.json")
	defer credential.Close()
	byteVal, _ := ioutil.ReadAll(credential)
	var res = new(models.Bot)
	json.Unmarshal([]byte(byteVal), res)
	log.Println(res.BaseAPI + res.Token)
	if string(res.BaseAPI) != "https://api.telegram.org/" {
		panic(fmt.Sprintf("wrong baseAPI"))
	}
	if strings.HasPrefix(res.Token, "bot") == false {
		panic(fmt.Sprintf("token should start with bot"))
	}
	return &appUsecase{
		appRepo:        a,
		contextTimeout: timeout,
		botURL:         res.BaseAPI + res.Token,
	}
}

func (a *appUsecase) Check(c context.Context) (string, error) {
	return "goods", nil
}

func (a *appUsecase) handleInsert(text string, username string) (string, error) {
	/*
		-----FORMAT-----
		insert\n
		task={}\n
		due_date={}\n --YYYY-MM-DD
		description={}
	*/
	input := strings.Split(text, "\n")
	if len(input) != 4 {
		return "failed", errors.New("Wrong Format")
	}
	var format [4]string
	format[1] = "task"
	format[2] = "due_date"
	format[3] = "description"
	for i := 1; i < 4; i++ {
		key := strings.Split(input[i], "=")
		if key[0] != format[i] {
			return "failed", errors.New("Wrong Format")
		}
		format[i] = key[1]
	}
	t, err := time.Parse(timeFormat, format[2])
	if err != nil {
		return "failed", errors.New("Wrong Date Format")
	}
	var tasks []models.Task
	task, err := a.appRepo.GetUndone(username)
	log.Println("ASDSADAS")
	log.Println(task, err)
	var flag bool
	for i := 0; i < len(task); i++ {
		if task[i].Task == format[1] {
			// if exist then just update
			flag = true
			tasks = append(tasks, models.Task{
				Task:        task[i].Task,
				DueDate:     t.Format(timeFormat),
				Description: format[3],
			})
		} else {
			tasks = append(tasks, models.Task{
				Task:        task[i].Task,
				DueDate:     task[i].DueDate,
				Description: task[i].Description,
			})
		}
	}
	if err == nil && flag == false {
		tasks = append(tasks, models.Task{
			Task:        format[1],
			DueDate:     t.Format(timeFormat),
			Description: format[3],
		})
	}
	byteJson, _ := json.Marshal(tasks)
	a.appRepo.Insert(byteJson, username)

	return "inserted", nil
}

func (a *appUsecase) Webhook(c context.Context, m *models.Webhook) (string, error) {
	_, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	var replyText string
	username := m.Message.From.Username
	text := m.Message.Text
	if text[:6] == "insert" {
		result, err := a.handleInsert(text, username)
		if err != nil {
			fmt.Errorf("error: %s", err)
			return "failed", err
		}
		replyText = result
	}
	message := models.SendMessage{
		ChatID:    m.Message.Chat.ID,
		Text:      replyText,
		MessageID: m.Message.MessageID,
	}
	bytesBody, err := json.Marshal(message)
	if err != nil {
		fmt.Errorf("error: %s", err)
	}
	log.Println(bytes.NewBuffer(bytesBody))
	req, err := http.Post(a.botURL+"/sendMessage", "application/json", bytes.NewBuffer(bytesBody))
	if err != nil {
		fmt.Errorf("error: %s", err)
		return "failed", err
	}
	fmt.Println(req)
	return "success", nil
}
