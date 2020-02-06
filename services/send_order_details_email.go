package services

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/go-gomail/gomail"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/templates"
)

func SendOrderDetailsEmail(name, email string, order *models.OrderDetailsView) error {
	params := map[string]interface{}{}
	params["greetings"] = fmt.Sprintf("Hi %s,", order.UserName)
	params["intros"] = "Thank you for purchasing from our platform."
	params["orderHash"] = order.Hash
	params["billingAddress"] = fmt.Sprintf("%s<br>%s<br>%s - %s",
		order.BillingAddress, order.BillingCity, order.BillingCountry, order.BillingPostcode)
	params["isShippable"] = !order.IsAllDigitalProducts

	if !order.IsAllDigitalProducts {
		params["shippingAddress"] = fmt.Sprintf("%s<br>%s<br>%s - %s",
			*order.ShippingAddress, *order.ShippingCity, *order.ShippingCountry, *order.ShippingPostcode)
	}

	params["shippingCharge"] = order.ShippingCharge
	params["paymentProcessingFee"] = order.PaymentProcessingFee
	params["paymentGateway"] = order.PaymentGateway
	params["subTotal"] = order.SubTotal
	params["grandTotal"] = order.GrandTotal

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
