package main

import (
	"github.com/chaiyeole/todo/config"
	redisRepository "github.com/chaiyeole/todo/repository/redis"
	"github.com/chaiyeole/todo/service/todoListService"
	"github.com/chaiyeole/todo/transport/http"
)

func main() {
	// function name
	config, err := config.New()
	if err != nil {
		return
	}
	// padh ke samajh aa jana chahiye
	newRedisRepo, err := redisRepository.New(*config)
	if err != nil {
		return
	}

	tasks, err := todoListService.New(newRedisRepo)
	if err != nil {
		return
	}

	svr := http.NewServer()

	http.InitHandlers(svr, tasks)
	// there is always only one host
	err = http.RunServer(svr, config.ConfigHTTP.Host)
	if err != nil {
		return
	}

}

// git remote add origin URL
