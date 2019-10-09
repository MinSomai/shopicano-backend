package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/middlewares"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/repository"
	"github.com/shopicano/shopicano-backend/utils"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
	"strconv"
	"time"
)

func RegisterAdminRoutes(g *echo.Group) {
	g.GET("/shipping-methods/", listShippingMethods)
	g.GET("/payment-methods/", listPaymentMethods)

	g.Use(middlewares.IsPlatformAdmin)
	g.POST("/shipping-methods/", createShippingMethod)
	g.PUT("/shipping-methods/:id/", updateShippingMethod)
	g.DELETE("/shipping-methods/:id/", deleteShippingMethod)
	g.GET("/shipping-methods/with-admin/", listShippingMethodsWithAdmin)

	g.POST("/payment-methods/", createPaymentMethod)
	g.PUT("/payment-methods/:id/", updatePaymentMethod)
	g.DELETE("/payment-methods/:id/", deletePaymentMethod)
	g.GET("/payment-methods/with-admin/", listPaymentMethodsWithAdmin)

	g.POST("/additional-charges/", createAdditionalCharge)
	g.PUT("/additional-charges/:id/", updateAdditionalCharge)
	g.DELETE("/additional-charges/:id/", deleteAdditionalCharge)
	g.GET("/additional-charges/with-admin/", listAdditionalChargesWithAdmin)

	func(g echo.Group) {
		g.Use(middlewares.IsStoreStaffWithStoreActivation)
		g.GET("/additional-charges/", listAdditionalCharges)
	}(*g)
}

func createShippingMethod(ctx echo.Context) error {
	req, err := validators.ValidateCreateShippingMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ShippingMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	m := &models.ShippingMethod{
		ID:                      utils.NewUUID(),
		Name:                    req.Name,
		IsPublished:             req.IsPublished,
		ApproximateDeliveryTime: req.ApproximateDeliveryTime,
		DeliveryCharge:          req.DeliveryCharge,
		WeightUnit:              req.WeightUnit,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}

	au := repository.NewAdminRepository()
	if err := au.CreateShippingMethod(m); err != nil {
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

func updateShippingMethod(ctx echo.Context) error {
	ID := ctx.Param("id")

	req, err := validators.ValidateCreateShippingMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ShippingMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	au := repository.NewAdminRepository()
	m, err := au.GetShippingMethod(ID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Shipping method not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ShippingMethodNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	m.Name = req.Name
	m.IsPublished = req.IsPublished
	m.ApproximateDeliveryTime = req.ApproximateDeliveryTime
	m.DeliveryCharge = req.DeliveryCharge
	m.WeightUnit = req.WeightUnit
	m.UpdatedAt = time.Now()

	if err := au.UpdateShippingMethod(m); err != nil {
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

func deleteShippingMethod(ctx echo.Context) error {
	ID := ctx.Param("id")

	resp := core.Response{}

	au := repository.NewAdminRepository()
	if err := au.DeleteShippingMethod(ID); err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Shipping method not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ShippingMethodNotFound
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

func listShippingMethods(ctx echo.Context) error {
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
	au := repository.NewAdminRepository()
	data, err := au.ListActiveShippingMethods(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
}

func listShippingMethodsWithAdmin(ctx echo.Context) error {
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
	au := repository.NewAdminRepository()
	data, err := au.ListShippingMethods(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
}

func createPaymentMethod(ctx echo.Context) error {
	req, err := validators.ValidateCreatePaymentMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.PaymentMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	m := &models.PaymentMethod{
		ID:          utils.NewUUID(),
		Name:        req.Name,
		IsPublished: req.IsPublished,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	au := repository.NewAdminRepository()
	if err := au.CreatePaymentMethod(m); err != nil {
		msg, ok := errors.IsDuplicateKeyError(err)
		if ok {
			resp.Title = msg
			resp.Status = http.StatusConflict
			resp.Code = errors.PaymentMethodAlreadyExists
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

func updatePaymentMethod(ctx echo.Context) error {
	ID := ctx.Param("id")

	req, err := validators.ValidateCreatePaymentMethod(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.PaymentMethodCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	au := repository.NewAdminRepository()
	m, err := au.GetPaymentMethod(ID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Payment method not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.PaymentMethodNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	m.Name = req.Name
	m.IsPublished = req.IsPublished
	m.UpdatedAt = time.Now()

	if err := au.UpdatePaymentMethod(m); err != nil {
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

func deletePaymentMethod(ctx echo.Context) error {
	ID := ctx.Param("id")

	resp := core.Response{}

	au := repository.NewAdminRepository()
	if err := au.DeletePaymentMethod(ID); err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Payment method not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.PaymentMethodNotFound
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

func listPaymentMethods(ctx echo.Context) error {
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
	au := repository.NewAdminRepository()
	data, err := au.ListActivePaymentMethods(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
}

func listPaymentMethodsWithAdmin(ctx echo.Context) error {
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
	au := repository.NewAdminRepository()
	data, err := au.ListPaymentMethods(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
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
		ID:          utils.NewUUID(),
		Name:        req.Name,
		Amount:      req.Amount,
		ChargeType:  models.AdditionalChargeType(req.ChargeType),
		AmountType:  models.AdditionalChargeAmountType(req.AmountType),
		AmountMin:   req.AmountMin,
		AmountMax:   req.AmountMax,
		IsPublished: req.IsPublished,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	au := repository.NewAdminRepository()
	if err := au.CreateAdditionalCharge(m); err != nil {
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
	ID := ctx.Param("id")

	req, err := validators.ValidateCreateAdditionalCharge(ctx)

	resp := core.Response{}

	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.AdditionalChargeDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	au := repository.NewAdminRepository()
	m, err := au.GetAdditionalCharge(ID)
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

	m.Name = req.Name
	m.IsPublished = req.IsPublished
	m.Amount = req.Amount
	m.AmountMin = req.AmountMin
	m.AmountType = models.AdditionalChargeAmountType(req.AmountType)
	m.ChargeType = models.AdditionalChargeType(req.ChargeType)
	m.AmountMax = req.AmountMax
	m.UpdatedAt = time.Now()

	if err := au.UpdateAdditionalCharge(m); err != nil {
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
	ID := ctx.Param("id")

	resp := core.Response{}

	au := repository.NewAdminRepository()
	if err := au.DeleteAdditionalCharge(ID); err != nil {
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

	from := (page - 1) * limit
	au := repository.NewAdminRepository()
	data, err := au.ListActiveAdditionalCharges(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
}

func listAdditionalChargesWithAdmin(ctx echo.Context) error {
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
	au := repository.NewAdminRepository()
	data, err := au.ListAdditionalCharges(int(from), int(limit))
	if err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = data
	return resp.ServerJSON(ctx)
}
