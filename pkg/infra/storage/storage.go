package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageAgent struct {
	client     *s3.S3
	bucketName string
}

func NewStorageAgent(client *s3.S3, bucketName string) StorageAgent {
	return StorageAgent{
		client:     client,
		bucketName: bucketName,
	}
}

func (s StorageAgent) UploadFile(ctx context.Context, file []byte, filename string) error {
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(file)),
		Bucket: &s.bucketName,
		Key:    &filename,
	},
	)
	return err
}

func (s StorageAgent) DownloadFile(ctx context.Context, filename string) ([]byte, error) {
	obj, err := s.client.GetObject(
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

func CreateAWSS3Client(endpoint, accessKey, secretKey, bucketName string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		fmt.Println(err.Error())
	}

	return s3Client
}
