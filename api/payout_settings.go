package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"time"
)

func createOrUpdatePayoutSettings(ctx echo.Context) error {
	req, err := validators.ValidateCreateShippingMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ShippingMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	m := &models.ShippingMethod{
		ID:                      utils.NewUUID(),
		Name:                    req.Name,
		IsPublished:             req.IsPublished,
		ApproximateDeliveryTime: req.ApproximateDeliveryTime,
		DeliveryCharge:          req.DeliveryCharge,
		WeightUnit:              req.WeightUnit,
		IsFlat:                  req.IsFlat,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}

	au := data.NewMarketplaceRepository()
	if err := au.CreateShippingMethod(db, m); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.ShippingMethodAlreadyExists
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

func getPayoutSettings(ctx echo.Context) error {
	req, err := validators.ValidateCreateShippingMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ShippingMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	m := &models.ShippingMethod{
		ID:                      utils.NewUUID(),
		Name:                    req.Name,
		IsPublished:             req.IsPublished,
		ApproximateDeliveryTime: req.ApproximateDeliveryTime,
		DeliveryCharge:          req.DeliveryCharge,
		WeightUnit:              req.WeightUnit,
		IsFlat:                  req.IsFlat,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}

	au := data.NewMarketplaceRepository()
	if err := au.CreateShippingMethod(db, m); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.ShippingMethodAlreadyExists
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
