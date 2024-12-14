package httpimplement

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpdriver "github.com/moazedy/todo/internal/port/driver/http"
	"github.com/moazedy/todo/internal/port/driver/service"
	"github.com/moazedy/todo/internal/port/driver/service/dto"
)

type file struct {
	fileService service.File
}

func NewFile(fileService service.File) httpdriver.File {
	return file{
		fileService: fileService,
	}
}

func (f file) Upload(ctx echo.Context) error {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	req := dto.UploadFileRequest{FileHeader: formFile}

	resp, err := f.fileService.Upload(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (f file) Download(ctx echo.Context) error {
	fileID := ctx.Param("file_name")
	req := dto.DownloadFileRequest{FileName: fileID}

	resp, err := f.fileService.Download(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	ctx.Response().Header().Set("Content-Disposition", "attachment; fileId="+fileID)
	ctx.Response().Header().Set("Content-Type", http.DetectContentType(resp.File))
	return ctx.Blob(http.StatusOK, http.DetectContentType(resp.File), resp.File)
}
