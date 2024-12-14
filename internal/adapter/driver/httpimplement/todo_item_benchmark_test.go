package httpimplement

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
	"github.com/moazedy/todo/pkg/infra/config"
)

func init() {
	cfg := config.Config{
		Postgres: config.PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "password",
			Name:     "todo_benchmark",
			Driver:   "postgres",
		},

		Storage: config.AwsS3Config{
			Endpoint:    "http://localhost:9000",
			Bucket:      "todo-benchmark",
			AccessKey:   "mys3accesskey",
			SecretKey:   "mys3secretkey",
			MaxFileSize: 500000,
		},

		Queue: config.SQS{
			IsMock: true,
		},
	}

	e := echo.New()
	register(e, cfg)
}

func BenchmarkCreateTodoItem(b *testing.B) {
	req := dto.CreateTodoItemRequest{
		Description: "some decs",
		DueDate:     time.Now(),
		FileName:    uuid.NewString() + "pdf",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := items.todoService.Create(context.Background(), req)
		if err != nil {
			b.Fatalf("Service failed: %v", err)
		}
	}
}

func createTempFile(fileName string, size int) (*os.File, error) {
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}

	data := make([]byte, size)
	_, err = rand.Read(data)
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Write(data)
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func BenchmarkUploadToS3(b *testing.B) {
	fileName := "large_test.txt"
	const fileSize = 21 * 1024 * 10 // 210 KB
	tempFile, err := createTempFile(fileName, fileSize)
	if err != nil {
		b.Fatalf("failed to create temp file : %v \n", err)
	}

	defer os.Remove(tempFile.Name()) // Cleanup after test

	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		b.Fatalf("failed to create form file : %v \n", err)
	}

	_, err = io.Copy(part, tempFile)
	if err != nil {
		b.Fatalf("failed to copy file : %v \n", err)
	}

	writer.Close()

	req := multipart.NewReader(&buff, writer.Boundary())
	form, err := req.ReadForm(32 << 20) // Limit the memory size for form parsing
	if err != nil {
		b.Fatalf("failed to read form: %v \n", err)
	}

	fileHeader := form.File["file"][0]

	uploadRequest := dto.UploadFileRequest{
		FileHeader: fileHeader,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := items.fileService.Upload(context.Background(), uploadRequest)
		if err != nil {
			b.Fatalf("Service failed: %v", err)
		}
	}
}

func BenchmarkSendMessageToSQS(b *testing.B) {
	req := dto.CreateTodoItemRequest{
		Description: "some decs",
		DueDate:     time.Now(),
		FileName:    uuid.NewString() + "pdf",
	}

	messageByte, err := json.Marshal(req)
	if err != nil {
		b.Fatalf("failed to marshal message : %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := items.queueRepo.SendMessage(context.Background(), string(messageByte))
		if err != nil {
			b.Fatalf("Service failed: %v", err)
		}
	}
}
