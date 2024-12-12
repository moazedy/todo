package repository

import (
	"context"

	"github.com/moazedy/todo/internal/domain/model"
)

type Queue interface {
	SendMessage(ctx context.Context, messageBody string) error
	ReceiveMessage(ctx context.Context, maxNumberOfMessages int32) ([]model.Message, error)
	DeleteMessage(ctx context.Context, messageId string) error
}
