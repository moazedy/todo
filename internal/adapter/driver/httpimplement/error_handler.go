package httpimplement

import (
	"github.com/labstack/echo"
	"github.com/moazedy/todo/pkg/cerror"
)

func CustomErrorHandlerMiddleWare(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		handlerErr := handler(ctx)
		err := cerror.GetCustomError(handlerErr)
		if err != nil {
			return ctx.JSON(
				err.GetHTTPCode(),
				err.Error(),
			)
		}

		return nil
	}
}
