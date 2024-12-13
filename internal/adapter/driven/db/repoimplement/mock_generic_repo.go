package repoimplement

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
)

func (grf genericRepoFactory[E]) newMockGenericRepo() repository.GenericRepo[E] {
	return newMockGenericRepo[E]()
}

type mockGenericRepo[E model.Entity] struct {
	db map[string]E
}

func newMockGenericRepo[E model.Entity]() repository.GenericRepo[E] {
	return &mockGenericRepo[E]{
		db: make(map[string]E),
	}
}

func (gr *mockGenericRepo[E]) Create(ctx context.Context, entityData E) (*E, error) {
	id := uuid.New().String()
	gr.db[id] = (entityData.WithIDSet(id)).(E)

	return &entityData, nil
}

func (gr *mockGenericRepo[E]) Update(ctx context.Context, entityData E) (*E, error) {
	_, exist := gr.db[entityData.GetID()]
	if !exist {
		return nil, errors.New("entity not found")
	}

	gr.db[entityData.GetID()] = entityData

	return &entityData, nil
}

func (gr *mockGenericRepo[E]) Delete(ctx context.Context, entityID string) error {
	_, exist := gr.db[entityID]
	if !exist {
		return errors.New("entity not found")
	}

	delete(gr.db, entityID)

	return nil
}

func (gr mockGenericRepo[E]) GetByStringField(ctx context.Context, fieldName, fieldValue string) (*E, error) {
	entityInstance, exist := gr.db[fieldValue]
	if !exist {
		return nil, errors.New("not found")
	}

	return &entityInstance, nil
}

func (gr mockGenericRepo[E]) List(context.Context) ([]E, error) {
	entityList := make([]E, 0)
	for _, val := range gr.db {
		entityList = append(entityList, val)
	}

	return entityList, nil
}

func (gr mockGenericRepo[E]) Exist(ctx context.Context, id string) (bool, error) {
	_, exist := gr.db[id]

	return exist, nil
}
