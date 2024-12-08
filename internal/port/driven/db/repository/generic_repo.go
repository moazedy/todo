package repository

import (
	"context"

	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/pkg/infra/tx"
)

type GenericRepo[E model.Entity] interface {
	Create(context.Context, E) (*E, error)
	Update(context.Context, E) (*E, error)
	Delete(context.Context, string) error
	GetByStringField(ctx context.Context, fieldName, fieldValue string) (*E, error)
	List(context.Context) ([]E, error)
}

type GenericRepoFactory[E model.Entity] interface {
	NewGenericRepo(tx tx.TX) GenericRepo[E]
}
