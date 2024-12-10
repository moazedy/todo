package repository

import (
	"context"
)

type File interface {
	Upload(context.Context, []byte, string) error
	Download(context.Context, string) ([]byte, error)
}
