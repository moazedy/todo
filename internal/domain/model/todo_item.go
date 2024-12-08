package model

import (
	"time"
)

type TodoItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId"`
}

func (ti TodoItem) IsEntity() {}
