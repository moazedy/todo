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
)

func TestCreateTodoItem_Success(t *testing.T) {
	todoService := NewTodoItem(mockTXFactory, mockRepoFactory, mockQueueRepo)

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
