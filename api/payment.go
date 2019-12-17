package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	gateway "github.com/shopicano/shopicano-backend/payment-gateways"
	"github.com/shopicano/shopicano-backend/validators"
	"net/http"
)

func RegisterPaymentRoutes(g *echo.Group) {
	g.POST("/brain-tree/", createPaymentBrainTree)
	g.POST("/stripe/", createPaymentStripe)
	g.POST("/success/", succeedPayment)
	g.POST("/failure/", failedPayment)
	g.GET("/token/", getPaymentGatewayClientToken)
}

func createPaymentBrainTree(ctx echo.Context) error {
	orderID := ctx.Request().URL.Query().Get("order_id")

	resp := core.Response{}

	req, err := validators.ValidateCreateReqBrainTreePayment(ctx)
	if err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.ProductCreationDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	db := app.DB()

	repo := data.NewOrderRepository()
	od, err := repo.GetDetailsInternal(db, orderID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ProductNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	od.Nonce = req.Nonce

	_, err = gateway.GetActivePaymentGateway().Pay(od)
	if err != nil {
		resp.Title = "Payment gateway failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = map[string]interface{}{
		"redirect": fmt.Sprintf("http://localhost:8080/#/confirmation?order_id=%s", od.ID),
	}

	return resp.ServerJSON(ctx)
}

func createPaymentStripe(ctx echo.Context) error {
	orderID := ctx.Param("order_id")

	resp := core.Response{}

	db := app.DB()

	repo := data.NewOrderRepository()
	od, err := repo.GetDetailsInternal(db, orderID)
	if err != nil {
		if errors.IsRecordNotFoundError(err) {
			resp.Title = "Product not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.ProductNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	res, err := gateway.GetActivePaymentGateway().Pay(od)
	if err != nil {
		resp.Title = "Payment gateway failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = map[string]interface{}{
		"session_id": res.Nonce,
	}

	return resp.ServerJSON(ctx)
}

func succeedPayment(ctx echo.Context) error {
	return nil
}

func failedPayment(ctx echo.Context) error {
	return nil
}

func getPaymentGatewayClientToken(ctx echo.Context) error {
	resp := core.Response{}

	token, err := gateway.GetActivePaymentGateway().GetClientToken()
	if err != nil {
		resp.Title = "Failed to get payment gateway client token"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = map[string]interface{}{
		"client_token": token,
	}
	return resp.ServerJSON(ctx)
}
