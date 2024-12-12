package srvimplement

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/moazedy/todo/internal/domain/errors"
	"github.com/moazedy/todo/internal/domain/value"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

type file struct {
	fileRepo    repository.File
	fileMaxSize int64
}

func NewFile(fileRepo repository.File, fileMaxSize int64) service.File {
	return file{
		fileRepo:    fileRepo,
		fileMaxSize: fileMaxSize,
	}
}

func (f file) Upload(ctx context.Context, req dto.UploadFileRequest) (out dto.UploadFileResponse, err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	if req.FileHeader.Size > f.fileMaxSize {
		err = errors.ErrFileIsOverSized
		return
	}

	file, err := req.FileHeader.Open()
	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}

	uId := uuid.New()

	fileType := strings.Split(http.DetectContentType(data), "/")[1]

	_, typeAllowed := value.AllowedFileTypes[fileType]
	if !typeAllowed {
		err = errors.ErrFileTypeNotAllowed
		return
	}

	fileName := fmt.Sprintf("%s.%s", uId.String(), fileType)

	err = f.fileRepo.Upload(ctx, data, fileName)
	if err != nil {
		return
	}

	out.FileName = fileName
	return
}

func (f file) Download(ctx context.Context, req dto.DownloadFileRequest) (out dto.DownloadFileResponse, err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	out.File, err = f.fileRepo.Download(ctx, req.FileName)

	return
}
