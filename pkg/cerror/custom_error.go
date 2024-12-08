package cerror

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type ErrorMessage string

type CustomError struct {
	httpCode int
	message  ErrorMessage
}

func (ce CustomError) Error() string {
	return string(ce.message)
}

func (ce CustomError) GetHTTPCode() int {
	return ce.httpCode
}

func NewInternalError(message ErrorMessage) CustomError {
	return CustomError{
		httpCode: http.StatusInternalServerError,
		message:  message,
	}
}

func NewNotFoundError(message ErrorMessage) CustomError {
	return CustomError{
		httpCode: http.StatusNotFound,
		message:  message,
	}
}

func NewBadRequestError(message ErrorMessage) CustomError {
	return CustomError{
		httpCode: http.StatusBadRequest,
		message:  message,
	}
}

func NewForbiddenError(message ErrorMessage) CustomError {
	return CustomError{
		httpCode: http.StatusForbidden,
		message:  message,
	}
}

func NewUnauthorizedError(message ErrorMessage) CustomError {
	return CustomError{
		httpCode: http.StatusUnauthorized,
		message:  message,
	}
}

func GetCustomError(err error) *CustomError {
	if err == nil {
		return nil
	}

	customError := new(CustomError)
	ok := errors.As(err, customError)

	if !ok {
		var echoError *echo.HTTPError
		newOk := errors.As(err, &echoError)
		if newOk {
			return &CustomError{
				httpCode: echoError.Code,
				message:  ErrorMessage(echoError.Error()),
			}
		}
		log.Println(err)
		return &CustomError{
			message:  "an internal error has occurred",
			httpCode: http.StatusInternalServerError,
		}
	}

	return customError
}
