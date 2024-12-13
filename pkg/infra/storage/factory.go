package storage

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageAgent interface {
	UploadFile(ctx context.Context, file []byte, filename string) error
	DownloadFile(ctx context.Context, filename string) ([]byte, error)
}

type StorageAgentFactory interface {
	NewStorageAgent() StorageAgent
}

type storageAgentFactory struct {
	isMock     bool
	client     *s3.S3
	bucketName string
}

func (saf storageAgentFactory) NewStorageAgent() StorageAgent {
	if saf.isMock {
		return newMockStorageAgent()
	} else {
		return newStorageAgent(saf.client, saf.bucketName)
	}
}

func NewStorageAgentFactory(isMock bool, client *s3.S3, bucketName string) StorageAgentFactory {
	return storageAgentFactory{
		isMock:     isMock,
		client:     client,
		bucketName: bucketName,
	}
}
