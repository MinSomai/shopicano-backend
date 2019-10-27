package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/repositories"
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

		uc := repositories.NewUserRepository()
		userID, userPermission, err := uc.GetPermission(token)
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

var IsSignUpEnabled = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		uc := repositories.NewUserRepository()
		ok, err := uc.IsSignUpEnabled()
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

		uc := repositories.NewUserRepository()
		ok, err := uc.IsStoreCreationEnabled()
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
