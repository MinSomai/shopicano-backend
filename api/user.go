package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/repositories"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"github.com/shopicano/shopicano-backend/values"
	"net/http"
)

func RegisterUserRoutes(g *echo.Group) {
	g.POST("/login/", login)
	g.GET("/logout/", logout)
	g.GET("/refresh-token/", refreshToken)

	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.PUT("/", update)
		g.GET("/", get)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsSignUpEnabled)
		g.POST("/sign-up/", register)
	}(*g)

	//g.PATCH("/:id/status/", r.updateStatus)
	//g.PATCH("/:id/permission/", r.updatePermission)
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

	uc := repositories.NewUserRepository()
	s, err := uc.Login(e, p)

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

	_, permission, err := uc.GetPermissionByUserID(s.UserID)
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

func register(ctx echo.Context) error {
	u, err := validators.ValidateRegister(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.UserSignUpDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	uc := repositories.NewUserRepository()

	u.Password, _ = utils.GeneratePassword(u.Password)
	u.PermissionID = values.UserGroupID

	if err := uc.Register(u); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = "User already register"
			resp.Status = http.StatusConflict
			resp.Code = errors.UserAlreadyExists
			resp.Errors = errors.NewError(msg)
			return resp.ServerJSON(ctx)
		}

		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Title = "User registration successful"
	resp.Data = u
	resp.Status = http.StatusCreated
	return resp.ServerJSON(ctx)
}

func update(ctx echo.Context) error {
	return nil
}

func get(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)

	resp := core.Response{}

	uc := repositories.NewUserRepository()
	u, err := uc.Get(userID)

	if err != nil {
		log.Log().Errorln(err)

		resp.Title = "Failed to get user profile"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Data = map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"status":     u.Status,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
		"permission": ctx.Get(utils.UserPermission),
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

	uc := repositories.NewUserRepository()
	if err := uc.Logout(token); err != nil {
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

	uc := repositories.NewUserRepository()
	s, err := uc.RefreshToken(token)
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
