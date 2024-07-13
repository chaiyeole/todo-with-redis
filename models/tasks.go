package models

import "github.com/google/uuid"

type Task struct {
	Id   uuid.UUID
	Name string
	Desc string
	Done bool
}
