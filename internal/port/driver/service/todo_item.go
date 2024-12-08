package service

import (
	"context"

	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

type TodoItem interface {
	Create(context.Context, dto.CreateTodoItemRequest) (dto.CreateTodoItemResponse, error)
	Update(context.Context, dto.UpdateTodoItemRequest) error
	Delete(context.Context, dto.DeleteTodoItemRequest) error
	GetByID(context.Context, dto.GetTodoItemByIDRequest) (dto.GetTodoItemByIDResponse, error)
}
