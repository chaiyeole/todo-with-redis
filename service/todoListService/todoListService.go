package todoListService

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/chaiyeole/todo/domain"
	"github.com/chaiyeole/todo/repository"
)

type store struct {
	repo  repository.IFileRepo
	mutex sync.Mutex
}

func New(repo repository.IFileRepo) (domain.ITask, error) {

	store := store{
		repo: repo,
		mutex: sync.Mutex{},
	}

	return &store, nil
}

func (s *store) Get(ctx context.Context, getReq *domain.GetReq) (*domain.GetRes, *domain.CustomError) {
	idString := getReq.Id.String()
	if (idString) == "" {
		err := errors.New("error UUID")

		slog.Error("Error while parsing UUID in API", "err", err)

		return nil, &domain.CustomError{
			StatusCode: http.StatusBadRequest,
			ErrMsg:     "Error while parsing UUID in API",
		}
	}

	task, err := s.repo.Get(ctx, getReq.Id.String())
	if err != nil {
		return nil, err
	}

	return &domain.GetRes{
		Task: *task,
	}, nil
}

// Get implements domain.ITask.
func (s *store) GetAll(ctx context.Context) (*domain.GetAllRes, *domain.CustomError) {
	taskList, err := s.repo.Load(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.GetAllRes{
		Task: taskList,
	}, nil
}

// Set implements domain.ITask.
func (s *store) Set(ctx context.Context, setReq *domain.SetReq) (*domain.SetRes, *domain.CustomError) {
	idString := setReq.Task.Id.String()
	if (idString) == "" {
		err := errors.New("error UUID")

		slog.Error("Error while parsing UUID in API", "err", err)

		return nil, &domain.CustomError{
			StatusCode: http.StatusBadRequest,
			ErrMsg:     "error while parsing UUID in API",
		}
	}

	err := s.repo.Set(ctx, setReq.Task.Id.String(), setReq.Task)
	if err != nil {
		return nil, err
	}

	return &domain.SetRes{
		Task: setReq.Task,
	}, nil
}

// Issue : getting empty values in response. I am storing map[uuid]task, but request and response format is json {key:uuid,value:task}
// Need to think over how to convert my map to json.
//
