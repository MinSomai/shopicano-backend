package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/repository"
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

		uc := repository.NewUserRepository()
		userID, userPermission, err := uc.GetPermission(token)
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
