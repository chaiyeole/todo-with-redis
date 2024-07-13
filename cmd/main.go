package main

import (
	"github.com/chaiyeole/todo/config"
	"github.com/chaiyeole/todo/repository"
	"github.com/chaiyeole/todo/service/todoListService"
	"github.com/chaiyeole/todo/transport/http"
)

func main() {
	config, err := config.ConfigReader()
	if err != nil {
		return
	}

	newRepo, err := repository.New(config.ConfigRedis.Addr, config.ConfigRedis.Password, config.ConfigRedis.DB)
	if err != nil {
		return
	}

	tasks, err := todoListService.New(newRepo)
	if err != nil {
		return
	}

	svr := http.NewServer()

	http.InitHandlers(svr, tasks)

	err = http.RunServer(svr, config.ConfigHTTP.Localhost)
	if err != nil {
		return
	}

}
