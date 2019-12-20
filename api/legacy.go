package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/queue"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"github.com/shopicano/shopicano-backend/values"
	"net/http"
)

func RegisterLegacyRoutes(g *echo.Group) {
	g.POST("/login/", login)
	g.GET("/logout/", logout)
	g.GET("/refresh-token/", refreshToken)
	g.GET("/email-verification/", emailVerification)

	func(g echo.Group) {
		g.Use(middlewares.IsSignUpEnabled)
		g.POST("/register/", register)
	}(*g)
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

	db := app.DB().Begin()

	uc := data.NewUserRepository()

	u.Password, _ = utils.GeneratePassword(u.Password)
	u.PermissionID = values.UserGroupID

	if err := uc.Register(db, u); err != nil {
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

	if err := queue.SendSignUpVerificationEmail(u.ID); err != nil {
		db.Rollback()

		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.FailedToEnqueueTask
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if err := db.Commit().Error; err != nil {
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

func emailVerification(ctx echo.Context) error {
	resp := core.Response{}

	userID := ctx.QueryParam("uid")
	token := ctx.QueryParam("token")

	db := app.DB().Begin()

	uc := data.NewUserRepository()

	u, err := uc.Get(db, userID)
	if err != nil {
		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "User not registered"
			resp.Status = http.StatusNotFound
			resp.Code = errors.UserNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Email verification failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if u.IsEmailVerified {
		resp.Title = "Email already verified"
		resp.Status = http.StatusOK
		return resp.ServerJSON(ctx)
	}

	if u.VerificationToken == nil || *u.VerificationToken != token {
		resp.Title = "Email verification failed"
		resp.Status = http.StatusForbidden
		resp.Code = errors.VerificationTokenIsInvalid
		return resp.ServerJSON(ctx)
	}

	u.VerificationToken = nil
	u.Status = models.UserActive
	u.IsEmailVerified = true
	if err := uc.Update(db, u); err != nil {
		resp.Title = "Email verification failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		return resp.ServerJSON(ctx)
	}

	if err := db.Commit().Error; err != nil {
		resp.Title = "Email verification failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		return resp.ServerJSON(ctx)
	}

	resp.Title = "Email verification succeed"
	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}
