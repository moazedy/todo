package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageAgent struct {
	client     *s3.Client
	bucketName string
}

func NewStorageAgent(client *s3.Client, bucketName string) StorageAgent {
	return StorageAgent{
		client:     client,
		bucketName: bucketName,
	}
}

func (s StorageAgent) UploadFile(ctx context.Context, file []byte, filename string) error {
	_, err := s.client.PutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(filename),
			Body:   bytes.NewReader(file),
		},
	)
	return err
}

func (s StorageAgent) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(filename),
		})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	base64Content := base64.StdEncoding.EncodeToString(data)

	return []byte(base64Content), nil
}
