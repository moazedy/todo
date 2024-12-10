package service

import (
	"context"

	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

type File interface {
	Upload(context.Context, dto.UploadFileRequest) (dto.UploadFileResponse, error)
	Download(context.Context, dto.DownloadFileRequest) (dto.DownloadFileResponse, error)
}
