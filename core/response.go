package core

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"strconv"
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
	//ctx.Response().Header().Set("cache-control", fmt.Sprintf("max-age=%d", 60*60*24))
	ctx.Response().Header().Set("X-Platform", "Shopicano")
	ctx.Response().Header().Set("X-Platform-Developer", "Coders Garage")
	ctx.Response().Header().Set("X-Platform-Connect", "www.shopicano.com")

	heightQ := ctx.QueryParam("height")
	widthQ := ctx.QueryParam("width")
	qualityQ := ctx.QueryParam("quality")

	height, _ := strconv.ParseInt(heightQ, 10, 64)
	width, _ := strconv.ParseInt(widthQ, 10, 64)
	quality, _ := strconv.ParseInt(qualityQ, 10, 64)

	img, err := imaging.Decode(object)
	if err != nil {
		return nil
	}

	log.Log().Infoln("H : ", height, ",W : ", width, ",Q : ", quality)

	if height != 0 && width != 0 {
		img = imaging.Fit(img, int(height), int(width), imaging.Lanczos)
	} else if height != 0 && width == 0 {
		img = imaging.Fit(img, int(height), 0, imaging.Lanczos)
	} else if height == 0 && width != 0 {
		img = imaging.Fit(img, 0, int(width), imaging.Lanczos)
	}

	if quality <= 0 || quality > 100 {
		return imaging.Encode(ctx.Response().Writer, img, imaging.JPEG)
	}
	return imaging.Encode(ctx.Response().Writer, img, imaging.JPEG, imaging.JPEGQuality(int(quality)))
}
