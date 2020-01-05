package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"github.com/shopicano/shopicano-backend/values"
	"net/http"
)

func RegisterStoreRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaff)
		g.GET("/", getStore)
	}(*g)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreAdmin)
		g.POST("/permission", addStoreStaff)
		g.PATCH("/permission", updateStoreStaffPermission)
		g.DELETE("/permission", deleteStoreStaff)
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

	db := app.DB().Begin()

	su := data.NewStoreRepository()
	if err := su.CreateStore(db, s); err != nil {
		db.Rollback()

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

	st := &models.Staff{
		UserID:       userID,
		StoreID:      s.ID,
		PermissionID: values.AdminGroupID,
		IsCreator:    true,
	}

	if err := su.AddStoreStuff(db, st); err != nil {
		db.Rollback()

		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.UserAlreadyExists
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if err := db.Commit().Error; err != nil {
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

	db := app.DB()

	su := data.NewStoreRepository()
	profile, err := su.GetStoreUserProfile(db, ctx.Get(utils.UserID).(string))
	if err != nil {
		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "Store not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.StoreNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = profile
	return resp.ServerJSON(ctx)
}

func addStoreStaff(ctx echo.Context) error {
	e, p, err := validators.ValidateCreateOrUpdateStoreStaff(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AddStoreStaffDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	uu := data.NewUserRepository()

	u, err := uu.GetByEmail(db, *e)
	if err != nil {
		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "User not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.UserNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	s := &models.Staff{
		UserID:       u.ID,
		StoreID:      utils.GetStoreID(ctx),
		PermissionID: *p,
		IsCreator:    false,
	}

	su := data.NewStoreRepository()
	err = su.AddStoreStuff(db, s)
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Title = "User added to store"
	return resp.ServerJSON(ctx)
}

func updateStoreStaffPermission(ctx echo.Context) error {
	resp := core.Response{}

	db := app.DB()

	su := data.NewStoreRepository()
	profile, err := su.GetStoreUserProfile(db, ctx.Get(utils.UserID).(string))
	if err != nil {
		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "Store not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.StoreNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = profile
	return resp.ServerJSON(ctx)
}

func deleteStoreStaff(ctx echo.Context) error {
	resp := core.Response{}

	db := app.DB()

	su := data.NewStoreRepository()
	profile, err := su.GetStoreUserProfile(db, ctx.Get(utils.UserID).(string))
	if err != nil {
		ok := errors.IsRecordNotFoundError(err)
		if ok {
			resp.Title = "Store not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.StoreNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = profile
	return resp.ServerJSON(ctx)
}
