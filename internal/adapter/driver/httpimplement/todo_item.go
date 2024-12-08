package httpimplement

import (
	"net/http"

	"github.com/labstack/echo"
	httpdriver "github.com/moazedy/todo/internal/port/driver/http"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
	"github.com/moazedy/todo/pkg/cerror"
)

type todoItem struct {
	todoItemService service.TodoItem
}

func NewTodoItem(todoItemService service.TodoItem) httpdriver.TodoItem {
	return todoItem{
		todoItemService: todoItemService,
	}
}

func (ti todoItem) Create(ctx echo.Context) error {
	var req dto.CreateTodoItemRequest
	if err := ctx.Bind(&req); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	res, err := ti.todoItemService.Create(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (ti todoItem) Update(ctx echo.Context) error {
	var req dto.UpdateTodoItemRequest
	if err := ctx.Bind(&req); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	err := ti.todoItemService.Update(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (ti todoItem) Delete(ctx echo.Context) error {
	var req dto.DeleteTodoItemRequest
	if err := ctx.Bind(&req); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	err := ti.todoItemService.Delete(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (ti todoItem) GetByID(ctx echo.Context) error {
	id := ctx.Param("id")
	req := dto.GetTodoItemByIDRequest{ID: id}

	res, err := ti.todoItemService.GetByID(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}
