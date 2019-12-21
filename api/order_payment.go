package api

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	payment_gateways "github.com/shopicano/shopicano-backend/payment-gateways"
	"io/ioutil"
	"net/http"
	"time"
)

// payOrder is the IPN callback
func payOrder(ctx echo.Context) error {
	orderID := ctx.Param("order_id")

	resp := core.Response{}

	b, _ := ioutil.ReadAll(ctx.Request().Body)
	log.Log().Infoln(string(b))
	for k, v := range ctx.QueryParams() {
		log.Log().Infoln(k, " Q<==>Q ", v[0])
	}

	for _, v := range ctx.ParamNames() {
		log.Log().Infoln(v, " P<==>P ", ctx.Param(v))
	}

	db := app.DB()

	ou := data.NewOrderRepository()
	m, err := ou.GetDetails(db, orderID)
	if err != nil {
		resp.Title = "Order not found"
		resp.Status = http.StatusNotFound
		resp.Code = errors.OrderNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if m.Status == models.PaymentCompleted {
		resp.Title = "Order already paid"
		resp.Status = http.StatusConflict
		resp.Code = errors.PaymentAlreadyProcessed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	switch m.PaymentGateway {
	case payment_gateways.BrainTreePaymentGatewayName:
		return processPayOrderForBrainTree(ctx, m)
	case payment_gateways.StripePaymentGatewayName:
		return processPayOrderForStripe(ctx, m)
	}
	return serveInvalidPaymentRequest(ctx)
}

type resBrainTreeNonce struct {
	Nonce *string `json:"nonce"`
}

func processPayOrderForBrainTree(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	db := app.DB()
	or := data.NewOrderRepository()

	body := resBrainTreeNonce{}
	if err := ctx.Bind(&body); err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.OrderPaymentDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.Nonce = body.Nonce

	res, err := payment_gateways.GetActivePaymentGateway().Pay(o)
	if err != nil {
		resp.Title = "Failed to process payment"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	now := time.Now().UTC()
	o.TransactionID = &res.Result
	o.Status = models.PaymentCompleted
	o.PaidAt = &now
	o.IsPaid = true

	if err := or.UpdatePaymentInfo(db, o); err != nil {
		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = map[string]interface{}{
		"transaction_id": res.Result,
	}
	return resp.ServerJSON(ctx)
}

func processPayOrderForStripe(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	db := app.DB()
	or := data.NewOrderRepository()

	if ctx.QueryParam("status") == "success" {
		o.Status = models.PaymentCompleted

		now := time.Now().UTC()
		o.PaidAt = &now
		o.IsPaid = true
	} else {
		o.Status = models.PaymentFailed
	}

	if err := or.UpdatePaymentInfo(db, o); err != nil {
		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

// generatePayNonce create payment reference / nonce
func generatePayNonce(ctx echo.Context) error {
	orderID := ctx.Param("order_id")

	resp := core.Response{}

	db := app.DB()

	ou := data.NewOrderRepository()
	m, err := ou.GetDetails(db, orderID)
	if err != nil {
		resp.Title = "Order not found"
		resp.Status = http.StatusNotFound
		resp.Code = errors.OrderNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if m.Status == models.PaymentCompleted {
		resp.Title = "Order already paid"
		resp.Status = http.StatusConflict
		resp.Code = errors.PaymentAlreadyProcessed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	switch m.PaymentGateway {
	case payment_gateways.StripePaymentGatewayName:
		return generateStripePayNonce(ctx, m)
	}

	return serveInvalidPaymentRequest(ctx)
}

func generateStripePayNonce(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	db := app.DB()
	or := data.NewOrderRepository()

	res, err := payment_gateways.GetActivePaymentGateway().Pay(o)
	if err != nil {
		resp.Title = "Failed to process payment"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.TransactionID = &res.Result
	o.Nonce = &res.Nonce

	if err := or.UpdatePaymentInfo(db, o); err != nil {
		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	resp.Status = http.StatusOK
	resp.Data = map[string]interface{}{
		"nonce": res.Nonce,
	}
	return resp.ServerJSON(ctx)
}

func serveInvalidPaymentRequest(ctx echo.Context) error {
	resp := core.Response{}
	resp.Title = "Invalid payment request"
	resp.Status = http.StatusForbidden
	resp.Code = errors.PaymentProcessingFailed
	return resp.ServerJSON(ctx)
}
