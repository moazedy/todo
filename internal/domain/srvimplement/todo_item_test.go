package srvimplement

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement"
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
	"github.com/moazedy/todo/pkg/infra/queue"
	"github.com/moazedy/todo/pkg/infra/tx"
	"github.com/stretchr/testify/assert"
)

var (
	mockTXFactory    = tx.NewTXFactory(true, nil)
	mockRepoFactory  = repoimplement.NewGenericRepoFactory[model.TodoItem](true)
	sqsClientFactory = queue.NewSQSClientFactory(true, nil)
	mockQueueRepo    = repoimplement.NewQueue(sqsClientFactory, "testqueue")
	todoService      = NewTodoItem(mockTXFactory, mockRepoFactory, mockQueueRepo)
)

func TestCreateTodoItem_Success(t *testing.T) {
	req := dto.CreateTodoItemRequest{
		Description: "some test desc",
		DueDate:     time.Now(),
		FileName:    uuid.NewString() + ".txt",
	}

	resp, err := todoService.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEqual(t, "", resp.ID)

	_, err = uuid.Parse(resp.ID)
	assert.Nil(t, err)
}

func TestCreateTodoItem_FileNameFailure(t *testing.T) {
	req := dto.CreateTodoItemRequest{
		Description: "some test desc",
		DueDate:     time.Now(),
	}

	_, err := todoService.Create(context.Background(), req)
	assert.ErrorContains(t, err, "Error:Field validation for 'FileName' failed on the 'required' tag")
}

func TestCreateTodoItem_DescriptionFailure(t *testing.T) {
	req := dto.CreateTodoItemRequest{
		DueDate:  time.Now(),
		FileName: uuid.NewString() + ".txt",
	}

	_, err := todoService.Create(context.Background(), req)
	assert.ErrorContains(t, err, "Error:Field validation for 'Description' failed on the 'required' tag")
}

func TestCreateTodoItem_DueDateFailure(t *testing.T) {
	req := dto.CreateTodoItemRequest{
		Description: "some test desc",
		FileName:    uuid.NewString() + ".txt",
	}

	_, err := todoService.Create(context.Background(), req)
	assert.ErrorContains(t, err, "Error:Field validation for 'DueDate' failed on the 'required' tag")
}

func TestUpdateTodoItem_Success(t *testing.T) {
	updateReq := dto.UpdateTodoItemRequest{
		ID:          uuid.NewString(),
		Description: "new desc",
		DueDate:     time.Now(),
		FileName:    uuid.NewString() + "pdf",
	}

	err := todoService.Update(context.Background(), updateReq)
	assert.Nil(t, err)
}

func TestUpdateTodoItem_IDFailure(t *testing.T) {
	updateReq := dto.UpdateTodoItemRequest{
		Description: "new desc",
		DueDate:     time.Now(),
		FileName:    uuid.NewString() + "pdf",
	}

	err := todoService.Update(context.Background(), updateReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'ID' failed on the 'required' tag")
}

func TestUpdateTodoItem_DescriptionFailure(t *testing.T) {
	updateReq := dto.UpdateTodoItemRequest{
		ID:       uuid.NewString(),
		DueDate:  time.Now(),
		FileName: uuid.NewString() + "pdf",
	}

	err := todoService.Update(context.Background(), updateReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'Description' failed on the 'required' tag")
}

func TestUpdateTodoItem_DueDateFailure(t *testing.T) {
	updateReq := dto.UpdateTodoItemRequest{
		ID:          uuid.NewString(),
		Description: "new desc",
		FileName:    uuid.NewString() + "jpg",
	}

	err := todoService.Update(context.Background(), updateReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'DueDate' failed on the 'required' tag")
}

func TestUpdateTodoItem_FileNameFailure(t *testing.T) {
	updateReq := dto.UpdateTodoItemRequest{
		ID:          uuid.NewString(),
		Description: "new desc",
		DueDate:     time.Now(),
	}

	err := todoService.Update(context.Background(), updateReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'FileName' failed on the 'required' tag")
}

func TestDeleteTodoItem_Success(t *testing.T) {
	delReq := dto.DeleteTodoItemRequest{
		ID: uuid.NewString(),
	}

	err := todoService.Delete(context.Background(), delReq)
	assert.Nil(t, err)
}

func TestDeleteTodoItem_Failure(t *testing.T) {
	delReq := dto.DeleteTodoItemRequest{}

	err := todoService.Delete(context.Background(), delReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'ID' failed on the 'required' tag")
}

func TestGetByIDTodoItem_Success(t *testing.T) {
	getReq := dto.GetTodoItemByIDRequest{
		ID: uuid.NewString(),
	}

	_, err := todoService.GetByID(context.Background(), getReq)
	assert.Nil(t, err)
}

func TestGetByIDTodoItem_Failure(t *testing.T) {
	getReq := dto.GetTodoItemByIDRequest{}

	_, err := todoService.GetByID(context.Background(), getReq)
	assert.ErrorContains(t, err, "Error:Field validation for 'ID' failed on the 'required' tag")
}
