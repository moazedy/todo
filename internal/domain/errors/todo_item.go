package errors

import "github.com/moazedy/todo/pkg/cerror"

var ErrTodoItemNotFound = cerror.NewNotFoundError("requested todo item not found")

var ErrFileTypeNotAllowed = cerror.NewForbiddenError("file type not allowed")

var ErrFileIsOverSized = cerror.NewForbiddenError("file is over sized")
