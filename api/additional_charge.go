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
	"net/http"
	"strconv"
	"time"
)

func RegisterAdditionalChargeRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
		g.POST("/", createAdditionalCharge)
		g.GET("/", listAdditionalCharges)
		g.GET("/:additional_charge_id/", getAdditionalCharge)
		g.DELETE("/:additional_charge_id/", deleteAdditionalCharge)
		g.PATCH("/:additional_charge_id/", updateAdditionalCharge)
	}(*g)
}

func createAdditionalCharge(ctx echo.Context) error {
	req, err := validators.ValidateCreateAdditionalCharge(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AdditionalChargeDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	m := &models.AdditionalCharge{
		ID:           utils.NewUUID(),
		StoreID:      utils.GetStoreID(ctx),
		Name:         req.Name,
		Amount:       req.Amount,
		AmountMin:    req.AmountMin,
		AmountMax:    req.AmountMax,
		IsFlatAmount: req.IsFlatAmount,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	db := app.DB()

	au := data.NewAdditionalChargeRepository()
	if err := au.Create(db, m); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.AdditionalChargeAlreadyExists
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
	resp.Data = m
	return resp.ServerJSON(ctx)
}

func updateAdditionalCharge(ctx echo.Context) error {
	ID := ctx.Param("additional_charge_id")

	resp := core.Response{}

	req := validators.ReqAdditionalChargeUpdate{}
	if err := ctx.Bind(&req); err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AdditionalChargeDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	au := data.NewAdditionalChargeRepository()
	m, err := au.Get(db, utils.GetStoreID(ctx), ID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Additional charge not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.AdditionalChargeNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.IsFlatAmount != nil {
		m.IsFlatAmount = *req.IsFlatAmount
	}
	if req.Amount != nil {
		m.Amount = *req.Amount
	}
	if req.AmountMin != nil {
		m.AmountMin = *req.AmountMin
	}
	if req.AmountMax != nil {
		m.AmountMax = *req.AmountMax
	}
	m.UpdatedAt = time.Now()

	if err := au.Update(db, m); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.AdditionalChargeAlreadyExists
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
	resp.Data = m
	return resp.ServerJSON(ctx)
}

func deleteAdditionalCharge(ctx echo.Context) error {
	ID := ctx.Param("additional_charge_id")

	resp := core.Response{}

	db := app.DB()

	au := data.NewAdditionalChargeRepository()
	if err := au.Delete(db, utils.GetStoreID(ctx), ID); err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Additional charge not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.AdditionalChargeNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusNoContent
	return resp.ServerJSON(ctx)
}

func getAdditionalCharge(ctx echo.Context) error {
	ID := ctx.Param("additional_charge_id")

	resp := core.Response{}

	db := app.DB()

	au := data.NewAdditionalChargeRepository()
	ac, err := au.Get(db, utils.GetStoreID(ctx), ID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Additional charge not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.AdditionalChargeNotFound
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
	resp.Data = ac
	return resp.ServerJSON(ctx)
}

func listAdditionalCharges(ctx echo.Context) error {
	pageQ := ctx.Request().URL.Query().Get("page")
	limitQ := ctx.Request().URL.Query().Get("limit")

	page, err := strconv.ParseInt(pageQ, 10, 64)
	if err != nil {
		page = 1
	}
	limit, err := strconv.ParseInt(limitQ, 10, 64)
	if err != nil {
		limit = 10
	}

	resp := core.Response{}

	db := app.DB()

	from := (page - 1) * limit
	au := data.NewAdditionalChargeRepository()
	d, err := au.List(db, utils.GetStoreID(ctx), int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = d
	return resp.ServerJSON(ctx)
}
