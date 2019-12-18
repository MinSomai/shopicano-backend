package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"github.com/shopicano/shopicano-backend/values"
	"net/http"
)

func RegisterUserRoutes(g *echo.Group) {
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

	db := app.DB()

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

	db := app.DB()

	uc := data.NewUserRepository()
	u, err := uc.Get(db, userID)

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
