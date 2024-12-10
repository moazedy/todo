package dto

import (
	"context"
	"time"

	"github.com/moazedy/todo/pkg/cerror"
	"github.com/moazedy/todo/pkg/validator"
)

type CreateTodoItemRequest struct {
	Description string    `json:"description" validate:"required,lte=100,gte=3"`
	DueDate     time.Time `json:"dueDate" validate:"required"`
	FileID      string    `json:"fileId" validate:"required,uuid"`
}

func (ctr CreateTodoItemRequest) Validate(ctx context.Context) error {
	if err := validator.Validate(ctx, ctr); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	return nil
}

type CreateTodoItemResponse struct {
	ID string `json:"id"`
}

type UpdateTodoItemRequest struct {
	ID          string    `json:"id" validate:"required,uuid"`
	Description string    `json:"description" validate:"required,lte=100,gte=3"`
	DueDate     time.Time `json:"dueDate" validate:"required"`
	FileID      string    `json:"fileId" validate:"required,uuid"`
}

func (utr UpdateTodoItemRequest) Validate(ctx context.Context) error {
	if err := validator.Validate(ctx, utr); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	return nil
}

type DeleteTodoItemRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (dtr DeleteTodoItemRequest) Validate(ctx context.Context) error {
	if err := validator.Validate(ctx, dtr); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	return nil
}

type GetTodoItemByIDRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (gtr GetTodoItemByIDRequest) Validate(ctx context.Context) error {
	if err := validator.Validate(ctx, gtr); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	return nil
}

type GetTodoItemByIDResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId"`
}
