package tasks

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/data"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	payment_gateways "github.com/shopicano/shopicano-backend/payment-gateways"
	"github.com/shopicano/shopicano-backend/services"
	"github.com/shopicano/shopicano-backend/templates"
	"time"
)

const (
	SendPaymentConfirmationEmailTaskName = "send_payment_confirmation_email"
	SendPaymentRevertedEmailTaskName     = "send_payment_reverted_email"
)

func SendPaymentConfirmationEmailFn(orderID string) error {
	db := app.DB()

	orderDao := data.NewOrderRepository()
	o, err := orderDao.GetDetails(db, orderID)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	userDao := data.NewUserRepository()
	u, err := userDao.Get(db, o.UserID)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := services.SendPaymentConfirmationEmail(u.Email, o); err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}
	return nil
}

func SendPaymentRevertedEmailFn(orderID string) error {
	db := app.DB()

	orderDao := data.NewOrderRepository()
	order, err := orderDao.GetDetails(db, orderID)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	userDao := data.NewUserRepository()
	u, err := userDao.Get(db, order.UserID)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	subject := "Your payment has been reverted"

	params := map[string]interface{}{}
	params["greetings"] = fmt.Sprintf("Hi %s,", order.UserName)
	params["intros"] = subject
	params["orderHash"] = order.Hash
	params["billingAddress"] = fmt.Sprintf("%s\n%s\n%s - %s",
		order.BillingAddress, order.BillingCity, order.BillingCountry, order.BillingPostcode)
	params["isShippable"] = !order.IsAllDigitalProducts

	if !order.IsAllDigitalProducts {
		params["shippingAddress"] = fmt.Sprintf("%s\n%s\n%s - %s",
			*order.ShippingAddress, *order.ShippingCity, *order.ShippingCountry, *order.ShippingPostcode)
	}

	params["shippingCharge"] = fmt.Sprintf("%.2f", float64(order.ShippingCharge)/100)
	params["paymentProcessingFee"] = fmt.Sprintf("%.2f", float64(order.PaymentProcessingFee)/100)
	params["subTotal"] = fmt.Sprintf("%.2f", float64(order.SubTotal)/100)
	params["grandTotal"] = fmt.Sprintf("%.2f", float64(order.GrandTotal)/100)
	params["isCouponApplied"] = false

	switch order.PaymentGateway {
	case payment_gateways.BrainTreePaymentGatewayName:
		params["paymentGateway"] = "Brain Tree"
	case payment_gateways.TwoCheckoutPaymentGatewayName:
		params["paymentGateway"] = "2Checkout"
	case payment_gateways.StripePaymentGatewayName:
		params["paymentGateway"] = "Stripe"
	}

	switch order.Status {
	case models.OrderPending:
		params["status"] = "Pending"
	case models.OrderConfirmed:
		params["status"] = "Confirmed"
	case models.OrderCancelled:
		params["status"] = "Cancelled"
	case models.OrderShipping:
		params["status"] = "Shipping"
	case models.OrderDelivered:
		params["status"] = "Delivered"
	}

	switch order.PaymentStatus {
	case models.PaymentPending:
		params["paymentStatus"] = "Pending"
	case models.PaymentCompleted:
		params["paymentStatus"] = "Completed"
	case models.PaymentFailed:
		params["paymentStatus"] = "Failed"
	case models.PaymentReverted:
		params["paymentStatus"] = "Reverted"
	}

	if order.DiscountedAmount != 0 {
		params["couponCode"] = order.CouponCode
		params["discount"] = fmt.Sprintf("%.2f", float64(order.DiscountedAmount)/100)
		params["isCouponApplied"] = true
	}

	var items []map[string]interface{}

	for _, v := range order.Items {
		items = append(items, map[string]interface{}{
			"name":     v.Name,
			"quantity": v.Quantity,
			"price":    fmt.Sprintf("%.2f", float64(v.Price)/100),
			"subTotal": fmt.Sprintf("%.2f", float64(v.SubTotal)/100),
		})
	}

	params["orderedItems"] = items

	body, err := templates.GenerateInvoiceEmailHTML(params)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}

	if err := services.SendEmail(subject, u.Email, body); err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), time.Second*30)
	}
	return nil
}
