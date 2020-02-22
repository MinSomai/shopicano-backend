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
	"strconv"
)

func RegisterStoreRoutes(publicEndpoints, platformEndpoints *echo.Group) {
	storesPublicPath := publicEndpoints.Group("/stores")
	storesPlatformPath := platformEndpoints.Group("/stores")

	func(g echo.Group) {
		g.GET("/:store_id/", getStore)
	}(*storesPublicPath)

	func(g echo.Group) {
		g.Use(middlewares.JWTAuth())
		g.Use(middlewares.HasStore())
		g.Use(middlewares.IsStoreManager())
		g.GET("/", getStoreForOwner)
	}(*storesPublicPath)

	func(g echo.Group) {
		g.Use(middlewares.JWTAuth())
		g.Use(middlewares.HasStore())
		g.Use(middlewares.IsStoreAdmin())
		g.POST("/:store_id/staffs/", addStoreStaff)
		g.PATCH("/:store_id/staffs/", updateStoreStaffPermission)
		g.DELETE("/:store_id/staffs/:user_id/", deleteStoreStaff)
		g.GET("/:store_id/staffs/", listStaffs)
	}(*storesPublicPath)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreCreationEnabled)
		g.Use(middlewares.JWTAuth())
		g.POST("/", createStore)
	}(*storesPublicPath)

	func(g echo.Group) {
		g.Use(middlewares.IsPlatformManager)
		g.GET("/", listStores)
		g.PATCH("/:store_id/", updateStoreAsPlatformOwner)
	}(*storesPlatformPath)
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

	storeID := ctx.Param("store_id")

	db := app.DB()

	su := data.NewStoreRepository()
	store, err := su.FindStoreByID(db, storeID)
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
	resp.Data = store
	return resp.ServerJSON(ctx)
}

func getStoreForOwner(ctx echo.Context) error {
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
	e, p, err := validators.ValidateCreateStoreStaff(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AddStoreStaffDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB().Begin()

	uu := data.NewUserRepository()
	su := data.NewStoreRepository()

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

	exists, err := su.IsAlreadyStaff(db, s.UserID)
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if exists {
		resp.Title = "User already staff"
		resp.Status = http.StatusForbidden
		resp.Code = errors.UserAlreadyStaff
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	err = su.AddStoreStuff(db, s)
	if err != nil {
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

	resp.Status = http.StatusOK
	resp.Title = "Staff added to store"
	return resp.ServerJSON(ctx)
}

func updateStoreStaffPermission(ctx echo.Context) error {
	uID, pID, err := validators.ValidateUpdateStoreStaff(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AddStoreStaffDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB().Begin()

	uu := data.NewUserRepository()
	su := data.NewStoreRepository()

	u, err := uu.Get(db, *uID)
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
		PermissionID: *pID,
		IsCreator:    false,
	}

	exists, err := su.IsAlreadyStaff(db, s.UserID)
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if !exists {
		resp.Title = "Staff doesn't exists"
		resp.Status = http.StatusNotFound
		resp.Code = errors.StaffDoesNotExists
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	err = su.UpdateStoreStuffPermission(db, s)
	if err != nil {
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

	resp.Status = http.StatusOK
	resp.Title = "Staff permission updated"
	return resp.ServerJSON(ctx)
}

func deleteStoreStaff(ctx echo.Context) error {
	uID := ctx.Param("user_id")

	resp := core.Response{}

	db := app.DB().Begin()
	su := data.NewStoreRepository()

	err := su.DeleteStoreStuffPermission(db, utils.GetStoreID(ctx), uID)
	if err != nil {
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

	resp.Status = http.StatusNoContent
	resp.Title = "Staff removed"
	return resp.ServerJSON(ctx)
}

func listStaffs(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")
	query := ctx.Request().URL.Query().Get("query")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	from := (page - 1) * limit

	resp := core.Response{}

	db := app.DB()
	su := data.NewStoreRepository()

	var r interface{}

	if query == "" {
		r, err = su.ListStaffs(db, utils.GetStoreID(ctx), int(from), int(limit))
	} else {
		r, err = su.SearchStaffs(db, utils.GetStoreID(ctx), query, int(from), int(limit))
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = r
	return resp.ServerJSON(ctx)
}

func listStores(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")
	query := ctx.Request().URL.Query().Get("query")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	from := (page - 1) * limit

	resp := core.Response{}

	db := app.DB()
	su := data.NewStoreRepository()

	var r interface{}

	if query == "" {
		r, err = su.List(db, int(from), int(limit))
	} else {
		r, err = su.Search(db, query, int(from), int(limit))
	}

	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = r
	return resp.ServerJSON(ctx)
}

func updateStoreAsPlatformOwner(ctx echo.Context) error {
	storeID := ctx.Param("store_id")

	status, commissionRate, err := validators.ValidateUpdateStoreStatus(ctx)

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
	store, err := su.FindStoreByID(db, storeID)
	if err != nil {
		db.Rollback()

		if errors.IsRecordNotFoundError(err) {
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

	if status != nil {
		store.Status = *status
	}
	if commissionRate != nil {
		store.CommissionRate = *commissionRate
	}

	if err := su.UpdateStoreStatus(db, store); err != nil {
		resp.Title = "Failed to update store"
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

	resp.Status = http.StatusOK
	resp.Data = store
	return resp.ServerJSON(ctx)
}
