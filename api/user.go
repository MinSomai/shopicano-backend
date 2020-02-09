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
	"net/http"
	"time"
)

func RegisterUserRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.PUT("/", update)
		g.GET("/", get)
	}(*g)

	//g.PATCH("/:id/status/", r.updateStatus)
	//g.PATCH("/:id/permission/", r.updatePermission)
}

func update(ctx echo.Context) error {
	req, err := validators.ValidateUserUpdate(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.StoreCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	userID := ctx.Get(utils.UserID).(string)

	db := app.DB().Begin()

	uc := data.NewUserRepository()
	u, err := uc.Get(db, userID)
	if err != nil {
		db.Rollback()
		log.Log().Errorln(err)

		resp.Title = "Failed to get user profile"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.ProfilePicture != nil {
		u.ProfilePicture = req.ProfilePicture
	}
	if req.Phone != nil {
		u.Phone = req.Phone
	}

	if req.Password != nil {
		pass, err := utils.GeneratePassword(*req.Password)
		if err != nil {
			db.Rollback()
			log.Log().Errorln(err)

			resp.Title = "Failed to generate password"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.PasswordEncryptionFailed
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		u.Password = pass
	}

	u.UpdatedAt = time.Now().UTC()

	err = uc.Update(db, u)
	if err != nil {
		db.Rollback()
		log.Log().Errorln(err)

		resp.Title = "Failed to update profile"
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

	resp.Data = map[string]interface{}{
		"id":              u.ID,
		"name":            u.Name,
		"email":           u.Email,
		"status":          u.Status,
		"phone":           u.Phone,
		"profile_picture": u.ProfilePicture,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
		"permission":      ctx.Get(utils.UserPermission),
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
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
		"id":              u.ID,
		"name":            u.Name,
		"email":           u.Email,
		"status":          u.Status,
		"phone":           u.Phone,
		"profile_picture": u.ProfilePicture,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
		"permission":      ctx.Get(utils.UserPermission),
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}
