package todoListService

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/chaiyeole/todo/domain"
	repository "github.com/chaiyeole/todo/repository/redis"
)

type store struct {
	repo repository.IFileRepo
}

func New(repo repository.IFileRepo) (domain.ITask, error) {
	store := store{
		repo: repo,
	}

	return &store, nil
}

// GetTask takes context and pointer to id ; returns pointer to task and pointer to customError ; implements domain.ITask.
func (s *store) GetTask(ctx context.Context, getTaskReq *domain.GetTaskReq) (*domain.GetTaskRes, *domain.CustomError) {
	task, err := s.repo.Get(ctx, getTaskReq.Id.String())
	if err != nil {
		return nil, err
	}

	return &domain.GetTaskRes{
		Task: *task,
	}, nil
}

// GetAllTasks takes context ; returns pointer to a struct of all tasks and pointer to customError ; implements domain.ITask.
func (s *store) GetAllTasks(ctx context.Context) (*domain.GetAllTasksRes, *domain.CustomError) {
	taskList, err := s.repo.Load(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.GetAllTasksRes{
		Task: taskList,
	}, nil
}

// Set takes context and task ; returns pointer to the same task and pointer to customeError ; implements domain.ITask.
func (s *store) SetTask(ctx context.Context, setTaskReq *domain.SetTaskReq) (*domain.SetTaskRes, *domain.CustomError) {
	idString := setTaskReq.Task.Id.String()
	if (idString) == "" {
		err := errors.New("error UUID")

		slog.Error("Error while parsing UUID in API", "err", err)

		return nil, &domain.CustomError{
			StatusCode: http.StatusBadRequest,
			ErrMsg:     "error while parsing UUID in API",
		}
	}

	err := s.repo.Set(ctx, setTaskReq.Task.Id.String(), setTaskReq.Task)
	if err != nil {
		return nil, err
	}

	return &domain.SetTaskRes{
		Task: setTaskReq.Task,
	}, nil
}

// Issue : getting empty values in response. I am storing map[uuid]task, but request and response format is json {key:uuid,value:task}
// Need to think over how to convert my map to json.
//
