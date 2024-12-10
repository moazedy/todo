package srvimplement

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/moazedy/todo/internal/port/driven/db/repository"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

type file struct {
	fileRepo repository.File
}

func NewFile(fileRepo repository.File) service.File {
	return file{
		fileRepo: fileRepo,
	}
}

func (f file) Upload(ctx context.Context, req dto.UploadFileRequest) (out dto.UploadFileResponse, err error) {
	if err = req.Validate(ctx); err != nil {
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

	// TODO : check allowed file types
	fileType := strings.Split(http.DetectContentType(data), "/")[1]
	fileName := fmt.Sprintf("%s.%s", uId.String(), fileType)

	err = f.fileRepo.Upload(ctx, data, fileName)
	if err != nil {
		return
	}

	out.FileID = uId.String()
	return
}

func (f file) Download(ctx context.Context, req dto.DownloadFileRequest) (out dto.DownloadFileResponse, err error) {
	if err = req.Validate(ctx); err != nil {
		return
	}

	// TODO : file is being saved by .type extention, but in request the type is bypassed.
	out.File, err = f.fileRepo.Download(ctx, req.FilID)

	return
}
