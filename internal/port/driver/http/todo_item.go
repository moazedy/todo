package http

import "github.com/labstack/echo/v4"

type TodoItem interface {
	Create(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
	GetByID(echo.Context) error
}
