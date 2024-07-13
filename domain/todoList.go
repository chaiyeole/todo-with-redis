package domain

import (
	"context"

	"github.com/chaiyeole/todo/models"
	"github.com/google/uuid"
)

type ITask interface {
	GetAll(ctx context.Context) (*GetAllRes, *CustomError)
	Get(ctx context.Context, getReq *GetReq) (*GetRes, *CustomError)
	Set(ctx context.Context, setReq *SetReq) (*SetRes, *CustomError)
}

type GetReq struct {
	Id uuid.UUID
}

type GetRes struct {
	Task models.Task
}

type SetReq struct {
	Task models.Task
}

type SetRes struct {
	Task models.Task
}

type GetAllRes struct {
	Task []models.Task
}



