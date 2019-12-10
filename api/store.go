package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
)

func RegisterStoreRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaff)
		g.GET("/", getStore)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreCreationEnabled)
		g.Use(middlewares.AuthUser)
		g.POST("/", createStore)
	}(*g)
}

func createStore(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)

	s, err := validators.ValidateCreateStore(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.StoreCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	su := data.NewStoreRepository()
	if err := su.CreateStore(s, userID); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.StoreAlreadyExists
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusCreated
	resp.Data = s
	return resp.ServerJSON(ctx)
}

func getStore(ctx echo.Context) error {
	resp := core.Response{}

	su := data.NewStoreRepository()
	profile, _ := su.GetStoreUserProfile(ctx.Get(utils.UserID).(string))

	resp.Status = http.StatusOK
	resp.Data = profile
	return resp.ServerJSON(ctx)
}
