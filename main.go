package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	_handler "work-relic-bot/src/app/delivery/http"
	_appRepo "work-relic-bot/src/app/repository"
	_appUsecase "work-relic-bot/src/app/usecase"
	"work-relic-bot/src/middleware"
	"work-relic-bot/src/models"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		panic(fmt.Sprintf("LevelDB open failed %s", err))
	}
	defer db.Close()
	byteVal, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(fmt.Sprintf("Error Opening Credentials %s", err))
	}
	log.Println("Success to open credentials.json")
	var res = new(models.Bot)
	json.Unmarshal([]byte(byteVal), res)
	if string(res.BaseAPI) != "https://api.telegram.org/" {
		panic(fmt.Sprintf("wrong baseAPI"))
	}
	if strings.HasPrefix(res.Token, "bot") == false {
		panic(fmt.Sprintf("token should start with bot"))
	}
	port := flag.String("port", ":3000", "Port for the server listen on")
	botURL := res.BaseAPI + res.Token
	flag.Parse()
	r := gin.Default()
	middL := middleware.InitMiddleware()
	r.Use(middL.CORS())
	timeoutContext := time.Duration(2) * time.Second
	ar := _appRepo.NewLevelDBRepository(db)
	au := _appUsecase.NewUsecase(ar, timeoutContext, botURL)
	_handler.NewHttpHandler(r, au)
	r.Run(*port)
}
