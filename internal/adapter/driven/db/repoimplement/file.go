package repoimplement

import (
	"context"

	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/pkg/infra/storage"
)

type file struct {
	storageAgentFactory storage.StorageAgentFactory
}

func NewFile(sa storage.StorageAgentFactory) repository.File {
	return file{
		storageAgentFactory: sa,
	}
}

func (f file) Upload(ctx context.Context, data []byte, filename string) error {
	storageAgent := f.storageAgentFactory.NewStorageAgent()
	return storageAgent.UploadFile(ctx, data, filename)
}

func (f file) Download(ctx context.Context, filename string) ([]byte, error) {
	storageAgent := f.storageAgentFactory.NewStorageAgent()
	return storageAgent.DownloadFile(ctx, filename)
}
