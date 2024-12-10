package dto

import (
	"context"
	"mime/multipart"

	"github.com/moazedy/todo/pkg/cerror"
	"github.com/moazedy/todo/pkg/validator"
)

type UploadFileRequest struct {
	FileHeader *multipart.FileHeader `json:"fileHeader"`
}

func (ufr UploadFileRequest) Validate(ctx context.Context) error {
	if ufr.FileHeader == nil {
		return cerror.NewBadRequestError("file header can not be empty")
	}

	if ufr.FileHeader.Filename == "" {
		return cerror.NewBadRequestError("file name can not be empty")
	}

	if ufr.FileHeader.Size == 0 {
		return cerror.NewBadRequestError("file has no size")
	}

	return nil
}

type UploadFileResponse struct {
	FileID string `json:"fileId"`
}

type DownloadFileRequest struct {
	FilID string `json:"fileId" validate:"required,uuid"`
}

func (dfr DownloadFileRequest) Validate(ctx context.Context) error {
	if err := validator.Validate(ctx, dfr); err != nil {
		return cerror.NewBadRequestError(cerror.ErrorMessage(err.Error()))
	}

	return nil
}

type DownloadFileResponse struct {
	File []byte `json:"file"`
}
