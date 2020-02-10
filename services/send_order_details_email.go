package services

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/go-gomail/gomail"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	payment_gateways "github.com/shopicano/shopicano-backend/payment-gateways"
	"github.com/shopicano/shopicano-backend/templates"
)

func SendOrderDetailsEmail(name, email string, order *models.OrderDetailsView) error {
	params := map[string]interface{}{}
	params["greetings"] = fmt.Sprintf("Hi %s,", order.UserName)
	params["intros"] = "Thank you for purchasing from our platform."
	params["orderHash"] = order.Hash
	params["billingAddress"] = fmt.Sprintf("%s\n%s\n%s - %s",
		order.BillingAddress, order.BillingCity, order.BillingCountry, order.BillingPostcode)
	params["isShippable"] = !order.IsAllDigitalProducts

	if !order.IsAllDigitalProducts {
		params["shippingAddress"] = fmt.Sprintf("%s\n%s\n%s - %s",
			*order.ShippingAddress, *order.ShippingCity, *order.ShippingCountry, *order.ShippingPostcode)
	}

	params["shippingCharge"] = order.ShippingCharge
	params["paymentProcessingFee"] = order.PaymentProcessingFee
	params["subTotal"] = order.SubTotal
	params["grandTotal"] = order.GrandTotal
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
		params["discount"] = order.DiscountedAmount
		params["isCouponApplied"] = true
	}

	var items []map[string]interface{}

	for _, v := range order.Items {
		items = append(items, map[string]interface{}{
			"name":     v.Name,
			"quantity": v.Quantity,
			"price":    v.Price,
			"subTotal": v.SubTotal,
		})
	}

	params["orderedItems"] = items

	body, err := templates.GenerateInvoiceEmailHTML(params)
	if err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), TaskRetryDelay)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s", config.EmailService().FromEmailAddress))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Order has been placed")
	m.SetBody("text/html", body)

	if err := EmailDialer().DialAndSend(m); err != nil {
		log.Log().Errorln(err)
		return tasks.NewErrRetryTaskLater(err.Error(), TaskRetryDelay)
	}
	return nil
}
