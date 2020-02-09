package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"net/http"
)

var IsPlatformAdmin = func(next echo.HandlerFunc) echo.HandlerFunc {
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
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		if *userPermission != models.AdminPerm {
			resp.Status = http.StatusForbidden
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.UserID, userID)
		ctx.Set(utils.UserPermission, userPermission)
		return next(ctx)
	}
}

var IsPlatformManager = func(next echo.HandlerFunc) echo.HandlerFunc {
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
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		if !(*userPermission == models.AdminPerm || *userPermission == models.ManagerPerm) {
			resp.Status = http.StatusForbidden
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.UserID, userID)
		ctx.Set(utils.UserPermission, userPermission)
		return next(ctx)
	}
}

var IsSignUpEnabled = func(next echo.HandlerFunc) echo.HandlerFunc {
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

var IsStoreCreationEnabled = func(next echo.HandlerFunc) echo.HandlerFunc {
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
