package errors

import "github.com/moazedy/todo/pkg/cerror"

var ErrTodoItemNotFound = cerror.NewNotFoundError("requested todo item not found")
