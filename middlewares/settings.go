package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"net/http"
)

func IsSignUpEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		db := app.DB()

		uc := data.NewUserRepository()
		ok, err := uc.IsSignUpEnabled(db)
		if err != nil {
			log.Log().Errorln(err)
			resp.Status = http.StatusInternalServerError
			resp.Title = "Database query failed"
			return resp.ServerJSON(ctx)
		}
		if !ok {
			resp.Status = http.StatusForbidden
			resp.Title = "Sign up disabled in settings"
			return resp.ServerJSON(ctx)
		}
		return next(ctx)
	}
}

func IsStoreCreationEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		db := app.DB()

		uc := data.NewUserRepository()
		ok, err := uc.IsStoreCreationEnabled(db)
		if err != nil {
			log.Log().Errorln(err)
			resp.Status = http.StatusInternalServerError
			resp.Title = "Database query failed"
			return resp.ServerJSON(ctx)
		}
		if !ok {
			resp.Status = http.StatusForbidden
			resp.Title = "Store creation disabled in settings"
			return resp.ServerJSON(ctx)
		}
		return next(ctx)
	}
}
