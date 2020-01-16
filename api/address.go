package api

import (
	"github.com/labstack/echo/v4"
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

func RegisterAddressRoutes(g *echo.Group) {
	func(g echo.Group) {
		g.Use(middlewares.AuthUser)
		g.POST("/", createAddress)
		g.GET("/", listAddresses)
		g.DELETE("/:id/", deleteAddress)
		g.PUT("/:id/", updateAddress)
	}(*g)
}

func createAddress(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)

	req, err := validators.ValidateCreateAddress(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.CollectionCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	add := &models.Address{
		ID:        utils.NewUUID(),
		Name:      req.Name,
		Street:    req.Street,
		City:      req.City,
		Postcode:  req.Postcode,
		Phone:     req.Phone,
		Email:     req.Email,
		Country:   req.Country,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	au := data.NewAddressRepository()
	if err := au.CreateAddress(add); err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusCreated
	resp.Data = add
	return resp.ServerJSON(ctx)
}

func updateAddress(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)
	addressID := ctx.Param("id")

	req, err := validators.ValidateCreateAddress(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.CollectionCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	add := &models.Address{
		ID:        addressID,
		Name:      req.Name,
		Street:    req.Street,
		City:      req.City,
		Postcode:  req.Postcode,
		Phone:     req.Phone,
		Email:     req.Email,
		Country:   req.Country,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	au := data.NewAddressRepository()
	if err := au.UpdateAddress(add); err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = add
	return resp.ServerJSON(ctx)
}

func deleteAddress(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)
	addressID := ctx.Param("id")

	resp := core.Response{}

	au := data.NewAddressRepository()
	if err := au.DeleteAddress(userID, addressID); err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Address not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.AddressNotFound
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

func listAddresses(ctx echo.Context) error {
	userID := ctx.Get(utils.UserID).(string)
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

	from := (page - 1) * limit
	au := data.NewAddressRepository()
	addresses, err := au.ListAddresses(userID, int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = addresses
	return resp.ServerJSON(ctx)
}
