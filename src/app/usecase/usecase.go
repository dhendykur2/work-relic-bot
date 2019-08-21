package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	app "work-relic-bot/src/app/repository"
	"work-relic-bot/src/models"
)

// IUsecase is interface to handle what usecase function
type IUsecase interface {
	Check(ctx context.Context) (string, error)
	Webhook(ctx context.Context, m *models.Webhook) (string, error)
}

type appUsecase struct {
	appRepo        app.IRepository
	contextTimeout time.Duration
	botURL         string
}

const (
	timeFormat   = "2006-01-02 15:04"
	statusDone   = "_done"
	statusUndone = "_undone"
)

// NewUsecase is to Register the App Usecase
func NewUsecase(a app.IRepository, timeout time.Duration, url string) IUsecase {
	return &appUsecase{
		appRepo:        a,
		contextTimeout: timeout,
		botURL:         url,
	}
}

func (a *appUsecase) Check(c context.Context) (string, error) {
	return "goods", nil
}

func (a *appUsecase) insert(text string, username string) string {
	/*
		-----FORMAT-----
		insert\n
		task={}\n
		due_date={}\n --YYYY-MM-DD
		description={}
	*/
	input := strings.Split(text, "\n")
	if len(input) != 4 {
		return "Wrong Format"
	}
	var format [4]string
	format[1] = "task"
	format[2] = "due_date"
	format[3] = "description"
	for i := 1; i < 4; i++ {
		key := strings.Split(input[i], "=")
		if key[0] != format[i] {
			return "Wrong Format"
		}
		format[i] = key[1]
	}
	if len(format[1]) > 20 {
		return "task title cannot more than 20 characters"
	}
	t, err := time.Parse(timeFormat, format[2])
	if err != nil {
		return "Wrong Date Format"
	} else if t.Before(time.Now()) {
		return "Due Date should after " + time.Now().Format(timeFormat)
	}
	var tasks []models.Task
	key := username + statusUndone
	task, err := a.appRepo.Get(key)
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
			// just map into tasks
			tasks = append(tasks, models.Task{
				Task:        task[i].Task,
				DueDate:     task[i].DueDate,
				Description: task[i].Description,
			})
		}
	}
	// when there is no error when get the data and the task is unique, then just append it
	if err == nil && flag == false {
		tasks = append(tasks, models.Task{
			Task:        format[1],
			DueDate:     t.Format(timeFormat),
			Description: format[3],
		})
	}
	byteJSON, _ := json.Marshal(tasks)
	a.appRepo.Put(key, byteJSON)

	return "inserted"
}

func (a *appUsecase) getTodos(text string, username string) string {
	/*
		-----FORMAT-----
		todos
	*/
	key := username + statusUndone
	tasks, err := a.appRepo.Get(key)
	sort.Slice(tasks, func(i, j int) bool {
		a, _ := time.Parse(timeFormat, tasks[i].DueDate)
		b, _ := time.Parse(timeFormat, tasks[j].DueDate)
		return a.Before(b)
	})
	var top int
	top = len(tasks)
	if top > 10 {
		top = 10
	}
	if top == 0 {
		return "No Work's for now"
	}
	var reply string
	reply += "----TOP 10 Work ToDo's----\n"
	if err == nil {
		for i := 0; i < top; i++ {
			reply += "----" + strconv.Itoa(i+1) + "----\n"
			reply += "Task\t: " + tasks[i].Task + "\n"
			reply += "Due Date\t: " + tasks[i].DueDate + "\n"
			reply += "Description\t: " + tasks[i].Description + "\n"
		}
	}
	return reply
}

func (a *appUsecase) getDoneTaks(text string, username string) string {
	/*
		-----FORMAT-----
		task_done
	*/
	key := username + statusDone
	tasks, err := a.appRepo.Get(key)
	var top int
	top = len(tasks)
	if top == 0 {
		return "No Work Done"
	}
	if top > 10 {
		top = 10
	}
	var reply string
	reply += "----TOP 10 Work Done's----\n"
	if err == nil {
		for i := 0; i < top; i++ {
			reply += "----" + strconv.Itoa(i+1) + "----\n"
			reply += "Task\t: " + tasks[i].Task + "\n"
			reply += "Due Date\t: " + tasks[i].DueDate + "\n"
			reply += "Description\t: " + tasks[i].Description + "\n"
			reply += "Done At\t: " + tasks[i].DoneAt + "\n"
		}
	}
	return reply
}

func RemoveIndex(s []models.Task, index int) []models.Task {
	return append(s[:index], s[index+1:]...)
}

func (a *appUsecase) insertDone(text string, username string) string {
	/*
		-----FORMAT-----
		done={task name}
	*/
	task := strings.Split(text, "=")
	keyUndone := username + statusUndone
	keyDone := username + statusDone
	undoneTasks, err := a.appRepo.Get(keyUndone)
	var tasksD []models.Task
	now := time.Now()
	doneAt, _ := time.Parse(timeFormat, now.Format(timeFormat))
	fmt.Println(task)
	if err == nil && len(undoneTasks) > 0 {
		for i := 0; i < len(undoneTasks); i++ {
			if undoneTasks[i].Task == task[1] {
				tasksD = append(tasksD, models.Task{
					Task:        undoneTasks[i].Task,
					DueDate:     undoneTasks[i].DueDate,
					Description: undoneTasks[i].Description,
					DoneAt:      doneAt.Format(timeFormat),
				})
				undoneTasks = RemoveIndex(undoneTasks, i)
				break
			}
		}
		if len(tasksD) < 1 {
			return "Work Not Found"
		}
		doneTasks, err := a.appRepo.Get(keyDone)
		if err == nil {
			for i := 0; i < len(doneTasks); i++ {
				tasksD = append(tasksD, models.Task{
					Task:        doneTasks[i].Task,
					DueDate:     doneTasks[i].DueDate,
					Description: doneTasks[i].Description,
					DoneAt:      doneTasks[i].DoneAt,
				})
			}
		}
		byteJSONDone, _ := json.Marshal(tasksD)
		a.appRepo.Put(keyDone, byteJSONDone)
		byteJSONUndone, _ := json.Marshal(undoneTasks)
		a.appRepo.Put(keyUndone, byteJSONUndone)
		return fmt.Sprintf("Task %s Done at %s", tasksD[0].Task, tasksD[0].DoneAt)
	}
	fmt.Println(len(undoneTasks), err)
	return "Work is empty"
}

func (a *appUsecase) Webhook(c context.Context, m *models.Webhook) (string, error) {
	_, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	var replyText string
	replyText = "wrong text command, checkout /help for documentation"
	username := m.Message.From.Username
	text := m.Message.Text
	if len(text) >= 5 && text[:5] == "/help" {
		replyText = "*Work Relic Bot*\n" +
			"1. to insert new Work type with this format:\n" +
			"insert\ntask=example\ndue_date=2019-05-05 09:00\ndescription=descriptionexample\n" +
			"2. to look all your undone work just type \"todo\"\n" +
			"3. to look all your done work just type \"work_done\"\n"
	} else if len(text) >= 6 && text[:6] == "insert" {
		result := a.insert(text, username)
		replyText = result
	} else if len(text) >= 4 && text[:4] == "todo" {
		result := a.getTodos(text, username)
		replyText = result
	} else if len(text) >= 5 && text[:5] == "done=" {
		result := a.insertDone(text, username)
		replyText = result
	} else if len(text) >= 9 && text[:9] == "work_done" {
		result := a.getDoneTaks(text, username)
		replyText = result
	}
	message := models.SendMessage{
		ChatID:    m.Message.Chat.ID,
		Text:      replyText,
		MessageID: m.Message.MessageID,
	}
	bytesBody, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	req, err := http.Post(a.botURL+"/sendMessage", "application/json", bytes.NewBuffer(bytesBody))
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return "failed", err
	}
	fmt.Println("send message to", m.Message.From.FirstName, "| statuscode:", req.StatusCode)
	return "success", nil
}
