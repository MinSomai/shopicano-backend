package core

import (
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go"
	"github.com/shopicano/shopicano-backend/errors"
	"io"
)

type Response struct {
	Code   errors.ErrorCode `json:"code,omitempty"`
	Status int              `json:"-"`
	Title  string           `json:"title,omitempty"`
	Data   interface{}      `json:"data,omitempty"`
	Errors error            `json:"errors,omitempty"`
}

func (r *Response) ServerJSON(ctx echo.Context) error {
	if err := ctx.JSON(r.Status, r); err != nil {
		return err
	}
	return nil
}

func (r *Response) ServeStreamFromMinio(ctx echo.Context, object *minio.Object) error {
	s, _ := object.Stat()
	ctx.Response().Header().Set("Content-Disposition", s.Metadata.Get("Content-Disposition"))
	ctx.Response().Header().Set("Content-Type", s.ContentType)

	if _, err := io.Copy(ctx.Response().Writer, object); err != nil {
		return err
	}
	return nil
}
