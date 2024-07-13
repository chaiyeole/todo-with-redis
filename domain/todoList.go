package domain

import (
	"context"

	"github.com/chaiyeole/todo/models"
	"github.com/google/uuid"
)

type ITask interface {
	GetAllTasks(ctx context.Context) (*GetAllTasksRes, *CustomError)
	GetTask(ctx context.Context, getReq *GetTaskReq) (*GetTaskRes, *CustomError)
	SetTask(ctx context.Context, setReq *SetTaskReq) (*SetTaskRes, *CustomError)
}

type GetTaskReq struct {
	Id uuid.UUID
}

type GetTaskRes struct {
	Task models.Task
}

type SetTaskReq struct {
	Task models.Task
}

type SetTaskRes struct {
	Task models.Task
}

type GetAllTasksRes struct {
	Task []models.Task
}



