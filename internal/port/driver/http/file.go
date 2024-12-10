package http

import "github.com/labstack/echo"

type File interface {
	Upload(echo.Context) error
	Download(echo.Context) error
}
