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
	return cerror.NewBadRequestError(cerror.ErrorMessage(validator.Validate(ctx, ctr).Error()))
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
	return cerror.NewBadRequestError(cerror.ErrorMessage(validator.Validate(ctx, utr).Error()))
}

type DeleteTodoItemRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (dtr DeleteTodoItemRequest) Validate(ctx context.Context) error {
	return cerror.NewBadRequestError(cerror.ErrorMessage(validator.Validate(ctx, dtr).Error()))
}

type GetTodoItemByIDRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

func (gtr GetTodoItemByIDRequest) Validate(ctx context.Context) error {
	return cerror.NewBadRequestError(cerror.ErrorMessage(validator.Validate(ctx, gtr).Error()))
}

type GetTodoItemByIDResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId"`
}
