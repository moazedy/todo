package repoimplement

import (
	"context"

	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/pkg/infra/storage"
)

type file struct {
	storageAgent storage.StorageAgent
}

func NewFile(sa storage.StorageAgent) repository.File {
	return file{
		storageAgent: sa,
	}
}

func (f file) Upload(ctx context.Context, data []byte, filename string) error {
	return f.storageAgent.UploadFile(ctx, data, filename)
}

func (f file) Download(ctx context.Context, filename string) ([]byte, error) {
	return f.storageAgent.DownloadFile(ctx, filename)
}
