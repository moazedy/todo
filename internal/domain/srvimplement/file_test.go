package srvimplement

import (
	"bytes"
	"context"
	"crypto/rand"
	"mime/multipart"
	"testing"

	"github.com/moazedy/todo/internal/adapter/driven/db/repoimplement"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
	"github.com/moazedy/todo/pkg/infra/storage"
	"github.com/stretchr/testify/assert"
)

var (
	mockStorageAgentFactory = storage.NewStorageAgentFactory(true, nil, "mybocket")
	mockFileRepo            = repoimplement.NewFile(mockStorageAgentFactory)
	fileService             = NewFile(mockFileRepo, 20)
)

func TestUpload_Success(t *testing.T) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fileName := "test.txt"
	fileContent := []byte("This is a test file.")

	part, err := writer.CreateFormFile("file", fileName)
	assert.NoError(t, err)

	_, err = part.Write(fileContent)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := multipart.NewReader(&b, writer.Boundary())
	form, err := req.ReadForm(1024)
	assert.NoError(t, err)

	fileHeader := form.File["file"][0]

	uploadRequest := dto.UploadFileRequest{
		FileHeader: fileHeader,
	}

	resp, err := fileService.Upload(context.Background(), uploadRequest)
	assert.Nil(t, err)
	assert.NotEqual(t, resp.FileName, "")
}

func TestUpload_FailureSize(t *testing.T) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fileName := "large_test.txt"

	const fileSize = 21 * 1024
	fileContent := make([]byte, fileSize)
	rand.Read(fileContent) // Fill with random data

	part, err := writer.CreateFormFile("file", fileName)
	assert.NoError(t, err)

	_, err = part.Write(fileContent)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := multipart.NewReader(&b, writer.Boundary())
	form, err := req.ReadForm(32 << 20) // Limit the memory size for form parsing
	assert.NoError(t, err)

	fileHeader := form.File["file"][0]

	uploadRequest := dto.UploadFileRequest{
		FileHeader: fileHeader,
	}

	_, err = fileService.Upload(context.Background(), uploadRequest)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "file is over sized")
}

/*
func TestUpload_FailureType(t *testing.T) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	fileName := "test"
	fileContent := []byte{2, 6, 7, 8}

	part, err := writer.CreateFormFile("file", fileName)
	assert.NoError(t, err)

	_, err = part.Write(fileContent)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := multipart.NewReader(&b, writer.Boundary())
	form, err := req.ReadForm(1024)
	assert.NoError(t, err)

	fileHeader := form.File["file"][0]

	uploadRequest := dto.UploadFileRequest{
		FileHeader: fileHeader,
	}

	_, err = fileService.Upload(context.Background(), uploadRequest)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "file type not allowed")
}
*/

func TestDownload_Success(t *testing.T) {
	req := dto.DownloadFileRequest{
		FileName: "testfilename",
	}
	resp, err := fileService.Download(context.Background(), req)
	assert.Nil(t, err)

	assert.Equal(t, "this is the test file", string(resp.File))
}

func TestDownload_FailureFileName(t *testing.T) {
	req := dto.DownloadFileRequest{
		// FileName: "testfilename",
	}
	_, err := fileService.Download(context.Background(), req)
	assert.NotNil(t, err)

	assert.ErrorContains(t, err, "Error:Field validation for 'FileName' failed on the 'required' tag")
}
