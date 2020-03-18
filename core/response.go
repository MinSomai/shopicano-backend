package core

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go"
	"github.com/shopicano/shopicano-backend/errors"
	"strings"
)

type Response struct {
	Code   errors.ErrorCode `json:"code,omitempty"`
	Status int              `json:"-"`
	Title  string           `json:"title,omitempty"`
	Data   interface{}      `json:"data,omitempty"`
	Errors error            `json:"errors,omitempty"`
}

func (r *Response) ServerJSON(ctx echo.Context) error {
	ctx.Response().Header().Set("X-Platform", "Shopicano")
	ctx.Response().Header().Set("X-Platform-Developer", "Coders Garage")
	ctx.Response().Header().Set("X-Platform-Connect", "www.shopicano.com")
	ctx.Response().Header().Set("Content-Type", "application/json")

	if err := ctx.JSON(r.Status, r); err != nil {
		return err
	}
	return nil
}

func (r *Response) ServeStreamFromMinio(ctx echo.Context, object *minio.Object) error {
	s, _ := object.Stat()

	fileName := fmt.Sprintf("%s.%s", s.ETag, s.Key[strings.LastIndex(s.Key, ".")+1:])
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	ctx.Response().Header().Set("Content-Type", s.ContentType)
	ctx.Response().Header().Set("cache-control", "max-age=3600")
	ctx.Response().Header().Set("X-Platform", "Shopicano")
	ctx.Response().Header().Set("X-Platform-Developer", "Coders Garage")
	ctx.Response().Header().Set("X-Platform-Connect", "www.shopicano.com")

	img, err := imaging.Decode(object)
	if err != nil {
		return nil
	}

	img = imaging.Resize(img, 0, 1024, imaging.Lanczos)

	if err := imaging.Encode(ctx.Response().Writer, img, imaging.PNG); err != nil {
		return err
	}
	return nil
}
