package middlewares

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"net/http"
)

var IsStoreStaffWithStoreActivation = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		token, err := utils.ParseBearerToken(ctx)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Access token missing"
			return resp.ServerJSON(ctx)
		}

		db := app.DB()

		uc := data.NewUserRepository()
		userID, _, err := uc.GetPermission(db, token)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		su := data.NewStoreRepository()
		store, err := su.GetStoreUserProfile(userID)

		if err != nil {
			log.Log().Errorln(err)
			if err == gorm.ErrRecordNotFound {
				resp.Status = http.StatusNotFound
				resp.Title = "Store not found"
				return resp.ServerJSON(ctx)
			}
			resp.Status = http.StatusInternalServerError
			resp.Title = "Database query failed"
			return resp.ServerJSON(ctx)
		}

		if store.Status == models.StoreRegistered || store.Status == models.StoreBanned || store.Status == models.StoreSuspended {
			resp.Status = http.StatusForbidden
			resp.Title = "Store isn't active"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.StoreID, store.ID)
		ctx.Set(utils.UserID, store.UserID)
		ctx.Set(utils.StorePermission, store.StorePermission)
		return next(ctx)
	}
}

var MightBeStoreStaffWithStoreActivation = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		token, err := utils.ParseBearerToken(ctx)
		if err != nil {
			return next(ctx)
		}

		db := app.DB()

		uc := data.NewUserRepository()
		userID, u, err := uc.GetPermission(db, token)
		if err != nil {
			return next(ctx)
		}

		su := data.NewStoreRepository()
		store, err := su.GetStoreUserProfile(userID)

		if err != nil {
			log.Log().Errorln(err)
			if err == gorm.ErrRecordNotFound {
				ctx.Set(utils.UserID, userID)
				ctx.Set(utils.UserPermission, *u)
				return next(ctx)
			}

			resp.Status = http.StatusInternalServerError
			resp.Title = "Database query failed"
			return resp.ServerJSON(ctx)
		}

		if store.Status != models.StoreActive {
			resp.Status = http.StatusForbidden
			resp.Title = "Store isn't active"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.StoreID, store.ID)
		ctx.Set(utils.UserID, store.UserID)
		ctx.Set(utils.StorePermission, store.StorePermission)
		return next(ctx)
	}
}

var IsStoreStaff = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		resp := core.Response{}

		token, err := utils.ParseBearerToken(ctx)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Access token missing"
			return resp.ServerJSON(ctx)
		}

		db := app.DB()

		uc := data.NewUserRepository()
		userID, _, err := uc.GetPermission(db, token)
		if err != nil {
			resp.Status = http.StatusUnauthorized
			resp.Title = "Unauthorized request"
			return resp.ServerJSON(ctx)
		}

		su := data.NewStoreRepository()
		store, err := su.GetStoreUserProfile(userID)

		if err != nil {
			log.Log().Errorln(err)
			if err == gorm.ErrRecordNotFound {
				resp.Status = http.StatusNotFound
				resp.Title = "Store not found"
				return resp.ServerJSON(ctx)
			}
			resp.Status = http.StatusInternalServerError
			resp.Title = "Database query failed"
			return resp.ServerJSON(ctx)
		}

		ctx.Set(utils.StoreID, store.ID)
		ctx.Set(utils.UserID, store.UserID)
		ctx.Set(utils.StorePermission, store.StorePermission)
		return next(ctx)
	}
}
