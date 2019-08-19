package app

import (
	"bytes"
	"context"
	"encoding/json"
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
}

type appUsecase struct {
	appRepo        app.IRepository
	contextTimeout time.Duration
	botURL         string
}

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

func (a *appUsecase) Webhook(c context.Context, m *models.Webhook) (string, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	log.Println(ctx)
	test := models.SendMessage{
		ChatID: m.Message.Chat.ID,
		Text:   m.Message.Text,
	}
	bytesBody, err := json.Marshal(test)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(bytes.NewBuffer(bytesBody))
	req, err := http.Post(a.botURL+"/sendMessage", "application/json", bytes.NewBuffer(bytesBody))
	if err != nil {
		log.Fatalln(err)
		return "failed", err
	}
	log.Println(req)
	return "success", nil
}

// func sendMessage(chatId, text) {

// }
