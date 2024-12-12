package srvimplement

import (
	"context"
	"encoding/json"
	"log"

	"github.com/moazedy/todo/internal/domain/errors"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/domain/srvimplement/cast"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
	"github.com/moazedy/todo/pkg/infra/tx"
)

const (
	fieldNameTodoItemID = "id"
)

type todoItem struct {
	txFactory           tx.TXFactory
	todoItemRepoFactory repository.GenericRepoFactory[model.TodoItem]
	queueRepo           repository.Queue
}

func NewTodoItem(txFactory tx.TXFactory, todoItemRepoFactory repository.GenericRepoFactory[model.TodoItem], queueRepo repository.Queue) service.TodoItem {
	return todoItem{
		txFactory:           txFactory,
		todoItemRepoFactory: todoItemRepoFactory,
		queueRepo:           queueRepo,
	}
}

func (ti todoItem) Create(ctx context.Context, req dto.CreateTodoItemRequest) (out dto.CreateTodoItemResponse, err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	tx := ti.txFactory.NewTX()
	out, err = ti.create(ctx, tx, req)
	err = tx.AutoCR(err)
	return
}

func (ti todoItem) create(ctx context.Context, tx tx.TX, req dto.CreateTodoItemRequest) (out dto.CreateTodoItemResponse, err error) {
	// TODO : if any kind of similarity between todo items is been considered as a constraint for creation process,
	// it needs to be implemented here.
	todoItemRepo := ti.todoItemRepoFactory.NewGenericRepo(tx)
	todoServiceModel, err := todoItemRepo.Create(ctx, cast.CreateTodoItemRequestToServiceModel(req))
	if err != nil {
		return
	}

	if todoServiceModel != nil {
		out = cast.ToCreateTodoItemResponse(*todoServiceModel)
	}

	go func() {
		byteMessage, _ := json.Marshal(*todoServiceModel)
		err := ti.queueRepo.SendMessage(ctx, string(byteMessage))
		if err != nil {
			log.Printf("error while pushing message to queue : %v \n", string(byteMessage))
		}
	}()

	return
}

func (ti todoItem) Update(ctx context.Context, req dto.UpdateTodoItemRequest) (err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	tx := ti.txFactory.NewTX()
	err = ti.update(ctx, tx, req)
	err = tx.AutoCR(err)
	return
}

func (ti todoItem) update(ctx context.Context, tx tx.TX, req dto.UpdateTodoItemRequest) (err error) {
	todoItemRepo := ti.todoItemRepoFactory.NewGenericRepo(tx)

	if err = ti.checkExistence(ctx, todoItemRepo, req.ID); err != nil {
		return
	}

	_, err = todoItemRepo.Update(ctx, cast.UpdateTodoItemRequestToServiceModel(req))

	return err
}

func (ti todoItem) Delete(ctx context.Context, req dto.DeleteTodoItemRequest) (err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	tx := ti.txFactory.NewTX()
	err = ti.delete(ctx, tx, req)
	err = tx.AutoCR(err)
	return
}

func (ti todoItem) delete(ctx context.Context, tx tx.TX, req dto.DeleteTodoItemRequest) (err error) {
	todoItemRepo := ti.todoItemRepoFactory.NewGenericRepo(tx)

	if err = ti.checkExistence(ctx, todoItemRepo, req.ID); err != nil {
		return
	}

	return todoItemRepo.Delete(ctx, req.ID)
}

func (ti todoItem) GetByID(ctx context.Context, req dto.GetTodoItemByIDRequest) (out dto.GetTodoItemByIDResponse, err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	tx := ti.txFactory.NewTX()
	out, err = ti.getByID(ctx, tx, req)
	err = tx.AutoCR(err)
	return
}

func (ti todoItem) getByID(ctx context.Context, tx tx.TX, req dto.GetTodoItemByIDRequest) (out dto.GetTodoItemByIDResponse, err error) {
	todoItemRepo := ti.todoItemRepoFactory.NewGenericRepo(tx)

	theItem, err := todoItemRepo.GetByStringField(ctx, fieldNameTodoItemID, req.ID)
	if err != nil {
		return
	}

	if theItem == nil {
		err = errors.ErrTodoItemNotFound
		return
	}

	out = cast.ToGetTodoItemByIDResponse(*theItem)
	return
}

func (ti todoItem) checkExistence(ctx context.Context, repo repository.GenericRepo[model.TodoItem], id string) error {
	itemExists, err := repo.Exist(ctx, id)
	if err != nil {
		return err
	}

	if !itemExists {
		return errors.ErrTodoItemNotFound
	}

	return nil
}
