package redisRepository

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/chaiyeole/todo/config"
	"github.com/chaiyeole/todo/domain"
	"github.com/chaiyeole/todo/models"
	"github.com/redis/go-redis/v9"
)

type IFileRepo interface {
	Load(ctx context.Context) ([]models.Task, *domain.CustomError)
	Get(ctx context.Context, id string) (*models.Task, *domain.CustomError)
	Set(ctx context.Context, id string, task models.Task) *domain.CustomError
}

type repo struct {
	rdb *redis.Client
}

// New takes config ; returns fileRepo (a redis client) and error
func New(config config.Config) (IFileRepo, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ConfigRedis.Addr,
		Password: config.ConfigRedis.Password, // no password set
		DB:       config.ConfigRedis.DB,       // use default DB
	})

	// check if connection is success or not by ping function
	redisPingResponse := rdb.Ping(context.Background())
	if redisPingResponse.Err() != nil {
		slog.Error("Error while connecting to redis", "err", redisPingResponse.Err())

		return nil, &domain.CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrMsg:     redisPingResponse.Err().Error(),
		}
	}

	return &repo{
		rdb: rdb,
	}, nil
}

// Load takes context ; returns a list of tasks and pointer to customError ; implements IFileRepo interface
func (r *repo) Load(ctx context.Context) ([]models.Task, *domain.CustomError) {
	redisResponse := r.rdb.Keys(ctx, "*")
	keys, err := redisResponse.Result()
	if err != nil {
		slog.Error("Error while loading all keys from redis", "err", err)

		return nil, &domain.CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrMsg:     "Error while loading all keys from redis",
		}
	}

	var taskList []models.Task

	for _, v := range keys {
		var task models.Task
		jsonTask := r.rdb.Get(ctx, v)
		if jsonTask.Err() != nil {
			slog.Error("Error while getting from redis in load function", "err", jsonTask.Err())

			return nil, &domain.CustomError{
				StatusCode: http.StatusBadRequest,
				ErrMsg:     jsonTask.Err().Error(),
			}
		}

		err = json.Unmarshal([]byte(jsonTask.Val()), &task)
		if err != nil {
			slog.Error("Error while unmarshalling from redis in load function", "err", err)

			return nil, &domain.CustomError{
				StatusCode: http.StatusInternalServerError,
				ErrMsg:     err.Error(),
			}
		}

		taskList = append(taskList, task)
	}

	return taskList, nil
}

// Set takes context, id, and task ; returns pointer to customeError ; implements IFileRepo interface
func (r *repo) Set(ctx context.Context, id string, task models.Task) *domain.CustomError {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		slog.Error("Error while marshalling task during set in redis", "err", err)

		return &domain.CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrMsg:     err.Error(),
		}
	}

	redisResponse := r.rdb.Set(ctx, id, jsonTask, 0)
	if redisResponse.Err() != nil {
		slog.Error("Error while setting a value in redis", "err", redisResponse.Err())

		return &domain.CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrMsg:     redisResponse.Err().Error(),
		}
	}

	return nil
}

// Get takes context and id ; returns pointer to task and pointer to customError
func (r *repo) Get(ctx context.Context, id string) (*models.Task, *domain.CustomError) {
	redisResponse := r.rdb.Get(ctx, id)

	var task models.Task

	err := json.Unmarshal([]byte(redisResponse.Val()), &task)
	if err != nil {
		slog.Error("Error while getting a task from redis", "err", err)

		return nil, &domain.CustomError{
			StatusCode: http.StatusInternalServerError,
			ErrMsg:     "Error while getting a task from redis",
		}
	}

	return &task, nil
}
