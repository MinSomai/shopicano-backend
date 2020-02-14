package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/utils"
	"net/http"
)

var AuthUser = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		token, err := utils.ParseBearerToken(ctx)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		db := app.DB()

		uc := data.NewUserRepository()
		userID, userPermission, err := uc.GetPermission(db, token)
		if err != nil {
			log.Log().Errorln(err)

			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.UserID, userID)
		ctx.Set(utils.UserPermission, *userPermission)
		return next(ctx)
	}
}

var AuthUserWithQueryToken = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		token := ctx.QueryParam("token")

		db := app.DB()

		uc := data.NewUserRepository()
		userID, userPermission, err := uc.GetPermission(db, token)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.UserID, userID)
		ctx.Set(utils.UserPermission, userPermission)
		return next(ctx)
	}
}
