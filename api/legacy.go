package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
)

func RegisterLegacyRoutes(g *echo.Group) {
	g.POST("/login/", login)
	g.GET("/logout/", logout)
	g.GET("/refresh-token/", refreshToken)
}

func login(ctx echo.Context) error {
	e, p, err := validators.ValidateLogin(ctx)

	resp := core.Response{}

	if err != nil {
		log.Log().Errorln(err)
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.UserLoginDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	uc := data.NewUserRepository()
	s, err := uc.Login(db, e, p)

	if err != nil {
		log.Log().Errorln(err)

		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Invalid login credentials"
			resp.Status = http.StatusUnauthorized
			resp.Code = errors.LoginCredentialsInvalid
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	_, permission, err := uc.GetPermissionByUserID(db, s.UserID)
	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Data = map[string]interface{}{
		"access_token":  s.AccessToken,
		"refresh_token": s.RefreshToken,
		"expire_on":     s.ExpireOn,
		"permission":    permission,
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

func logout(ctx echo.Context) error {
	token, err := utils.ParseBearerToken(ctx)

	resp := core.Response{}

	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Invalid data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.BearerTokenNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	uc := data.NewUserRepository()
	if err := uc.Logout(db, token); err != nil {
		log.Log().Errorln(err)

		resp.Title = "Failed to logout"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Title = "Logout successful"
	resp.Status = http.StatusNoContent
	return resp.ServerJSON(ctx)
}

func refreshToken(ctx echo.Context) error {
	token, err := utils.ParseBearerToken(ctx)

	resp := core.Response{}

	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Invalid data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.BearerTokenNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	uc := data.NewUserRepository()
	s, err := uc.RefreshToken(db, token)
	if err != nil {
		log.Log().Errorln(err)

		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "User already register"
			resp.Status = http.StatusBadRequest
			resp.Code = errors.BearerTokenNotFound
			resp.Errors = errors.NewError("Invalid refresh token")
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Failed to generate refresh token"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Title = "Refresh token generation successful"
	resp.Status = http.StatusOK
	resp.Data = s
	return resp.ServerJSON(ctx)
}
