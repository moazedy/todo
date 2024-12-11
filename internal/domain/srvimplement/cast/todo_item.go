package cast

import (
	"github.com/moazedy/todo/internal/domain/model"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

func CreateTodoItemRequestToServiceModel(in dto.CreateTodoItemRequest) model.TodoItem {
	return model.TodoItem{
		Description: in.Description,
		DueDate:     in.DueDate,
		FileName:    in.FileName,
	}
}

func ToCreateTodoItemResponse(in model.TodoItem) dto.CreateTodoItemResponse {
	return dto.CreateTodoItemResponse{
		ID: in.ID,
	}
}

func UpdateTodoItemRequestToServiceModel(in dto.UpdateTodoItemRequest) model.TodoItem {
	return model.TodoItem{
		ID:          in.ID,
		Description: in.Description,
		DueDate:     in.DueDate,
		FileName:    in.FileName,
	}
}

func ToGetTodoItemByIDResponse(in model.TodoItem) dto.GetTodoItemByIDResponse {
	return dto.GetTodoItemByIDResponse{
		ID:          in.ID,
		Description: in.Description,
		DueDate:     in.DueDate,
		FileName:    in.FileName,
	}
}
