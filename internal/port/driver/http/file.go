package http

import "github.com/labstack/echo/v4"

type File interface {
	Upload(echo.Context) error
	Download(echo.Context) error
}
