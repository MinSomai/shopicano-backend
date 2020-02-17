package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/core"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	payment_gateways "github.com/shopicano/shopicano-backend/payment-gateways"
	"github.com/shopicano/shopicano-backend/queue"
	"github.com/shopicano/shopicano-backend/utils"
	"io/ioutil"
	"net/http"
	"time"
)

// payOrder is the IPN callback
func payOrder(ctx echo.Context) error {
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

	if m.PaymentStatus == models.PaymentCompleted {
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

type reqBrainTreeNonce struct {
	Nonce *string `json:"nonce"`
}

func processPayOrderForBrainTree(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	db := app.DB().Begin()
	or := data.NewOrderRepository()

	body := reqBrainTreeNonce{}
	if err := ctx.Bind(&body); err != nil {
		resp.Title = "Invalid data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.OrderPaymentDataInvalid
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.Nonce = body.Nonce

	pg, err := payment_gateways.GetPaymentGatewayByName(o.PaymentGateway)
	if err != nil {
		resp.Title = "Invalid payment gateway"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	res, err := pg.Pay(o)
	if err != nil {
		resp.Title = "Failed to process payment"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	o.TransactionID = &res.Result
	o.PaymentStatus = models.PaymentCompleted

	if err := or.UpdatePaymentInfo(db, o); err != nil {
		db.Rollback()

		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	ol := models.OrderLog{
		ID:        utils.NewUUID(),
		OrderID:   o.ID,
		Action:    string(o.PaymentStatus),
		Details:   "Payment has been updated using BrainTree",
		CreatedAt: time.Now(),
	}
	if err := or.CreateLog(db, &ol); err != nil {
		db.Rollback()

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if o.PaymentStatus == models.PaymentCompleted {
		if err := queue.SendPaymentConfirmationEmail(o.ID); err != nil {
			db.Rollback()

			resp.Title = "Failed to enqueue task"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.FailedToEnqueueTask
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	if err := db.Commit().Error; err != nil {
		resp.Title = "Database query failed"
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

	db := app.DB().Begin()
	or := data.NewOrderRepository()

	pg, err := payment_gateways.GetPaymentGatewayByName(o.PaymentGateway)
	if err != nil {
		db.Rollback()

		resp.Title = "Invalid payment gateway"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentGatewayFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if err := pg.ValidateTransaction(o); err != nil {
		log.Log().Errorln(err)
		o.PaymentStatus = models.PaymentFailed
	} else {
		o.PaymentStatus = models.PaymentCompleted
	}

	if err := or.UpdatePaymentInfo(db, o); err != nil {
		db.Rollback()

		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	ol := models.OrderLog{
		ID:        utils.NewUUID(),
		OrderID:   o.ID,
		Action:    string(o.PaymentStatus),
		Details:   "Payment has been updated using Stripe",
		CreatedAt: time.Now(),
	}
	if err := or.CreateLog(db, &ol); err != nil {
		db.Rollback()

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if o.PaymentStatus == models.PaymentCompleted {
		if err := queue.SendPaymentConfirmationEmail(o.ID); err != nil {
			db.Rollback()

			resp.Title = "Failed to enqueue task"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.FailedToEnqueueTask
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	if err := db.Commit().Error; err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	paymentCompletedCallback := fmt.Sprintf("%s/#/ordertrack/%s", config.App().FrontStoreUrl, o.ID)
	return ctx.Redirect(http.StatusPermanentRedirect, paymentCompletedCallback)
}

func processPayOrderFor2Checkout(ctx echo.Context) error {
	orderID := ctx.QueryParam("merchant_order_id")

	b, _ := ioutil.ReadAll(ctx.Request().Body)
	log.Log().Infoln(string(b))

	resp := core.Response{}

	db := app.DB().Begin()

	ou := data.NewOrderRepository()
	m, err := ou.GetDetails(db, orderID)
	if err != nil {
		db.Rollback()

		resp.Title = "Order not found"
		resp.Status = http.StatusNotFound
		resp.Code = errors.OrderNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if m.PaymentStatus == models.PaymentCompleted {
		db.Rollback()

		resp.Title = "Order already paid"
		resp.Status = http.StatusConflict
		resp.Code = errors.PaymentAlreadyProcessed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if m.PaymentGateway != payment_gateways.TwoCheckoutPaymentGatewayName {
		db.Rollback()
		return serveInvalidPaymentRequest(ctx)
	}

	if ctx.QueryParam("credit_card_processed") == "Y" {
		m.PaymentStatus = models.PaymentCompleted
		trx := ctx.QueryParam("invoice_id")
		m.TransactionID = &trx
	} else {
		m.PaymentStatus = models.PaymentFailed
	}

	or := data.NewOrderRepository()

	if err := or.UpdatePaymentInfo(db, m); err != nil {
		db.Rollback()

		resp.Title = "Failed to update payment info"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	ol := models.OrderLog{
		ID:        utils.NewUUID(),
		OrderID:   m.ID,
		Action:    string(m.PaymentStatus),
		Details:   "Payment has been updated using 2Checkout",
		CreatedAt: time.Now(),
	}
	if err := or.CreateLog(db, &ol); err != nil {
		db.Rollback()

		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	if m.PaymentStatus == models.PaymentCompleted {
		if err := queue.SendPaymentConfirmationEmail(m.ID); err != nil {
			db.Rollback()

			resp.Title = "Failed to enqueue task"
			resp.Status = http.StatusInternalServerError
			resp.Code = errors.FailedToEnqueueTask
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
	}

	if err := db.Commit().Error; err != nil {
		resp.Title = "Database query failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	paymentCompletedCallback := fmt.Sprintf("%s/#/ordertrack/%s", config.App().FrontStoreUrl, orderID)
	return ctx.Redirect(http.StatusPermanentRedirect, paymentCompletedCallback)
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

	if m.PaymentStatus == models.PaymentCompleted {
		resp.Title = "Order already paid"
		resp.Status = http.StatusConflict
		resp.Code = errors.PaymentAlreadyProcessed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	switch m.PaymentGateway {
	case payment_gateways.StripePaymentGatewayName:
		return generateStripePayNonce(ctx, m)
	case payment_gateways.TwoCheckoutPaymentGatewayName:
		return generate2CheckoutPayUrl(ctx, m)
	}
	return serveInvalidPaymentRequest(ctx)
}

func generateStripePayNonce(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	db := app.DB()
	or := data.NewOrderRepository()

	pg, err := payment_gateways.GetPaymentGatewayByName(o.PaymentGateway)
	if err != nil {
		resp.Title = "Invalid payment gateway"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	res, err := pg.Pay(o)
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

func generate2CheckoutPayUrl(ctx echo.Context, o *models.OrderDetailsView) error {
	resp := core.Response{}

	pg, err := payment_gateways.GetPaymentGatewayByName(o.PaymentGateway)
	if err != nil {
		resp.Title = "Invalid payment gateway"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	res, err := pg.Pay(o)
	if err != nil {
		resp.Title = "Failed to process payment"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.PaymentProcessingFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}

	url := res.Result
	log.Log().Infoln(url)
	return ctx.Redirect(http.StatusPermanentRedirect, url)
}

func serveInvalidPaymentRequest(ctx echo.Context) error {
	resp := core.Response{}
	resp.Title = "Invalid payment request"
	resp.Status = http.StatusForbidden
	resp.Code = errors.PaymentProcessingFailed
	return resp.ServerJSON(ctx)
}
