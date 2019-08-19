package main

import (
	"flag"
	"fmt"
	"time"

	_handler "work-relic-bot/src/app/delivery/http"
	_appRepo "work-relic-bot/src/app/repository"
	_appUsecase "work-relic-bot/src/app/usecase"
	"work-relic-bot/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	listenAddress string
)

func main() {
	db, err := leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		panic(fmt.Sprintf("LevelDB open failed %s", err))
	}
	defer db.Close()
	port := flag.String("port", ":3000", "Port for the server listen on")
	flag.Parse()
	r := gin.Default()
	middL := middleware.InitMiddleware()
	r.Use(middL.CORS())
	timeoutContext := time.Duration(2) * time.Second
	ar := _appRepo.NewLevelDBRepository(db)
	au := _appUsecase.NewUsecase(ar, timeoutContext)
	_handler.NewHttpHandler(r, au)
	r.Run(*port)
}
