package payment_gateways

import (
	"fmt"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

type stripePaymentGateway struct {
	SecretKey       string
	SuccessCallback string
	FailureCallback string
}

func NewStripePaymentGateway(cfg map[string]interface{}) (*stripePaymentGateway, error) {
	return &stripePaymentGateway{
		SecretKey:       cfg["secret_key"].(string),
		SuccessCallback: cfg["success_callback"].(string),
		FailureCallback: cfg["failure_callback"].(string),
	}, nil
}

func (spg *stripePaymentGateway) GetName() string {
	return "stripe"
}

func (spg *stripePaymentGateway) Pay(orderDetails *models.OrderDetailsInternal) (*PaymentGatewayResponse, error) {
	stripe.Key = spg.SecretKey

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, op := range orderDetails.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			//Name:        stripe.String(op.Name),
			Amount:      stripe.Int64(int64(op.Price)),
			Currency:    stripe.String("usd"),
			Description: stripe.String(op.ProductID),
			Quantity:    stripe.Int64(int64(op.Quantity)),
		})
	}

	log.Log().Infoln(orderDetails.PaymentProcessingFee)

	lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
		Name:     stripe.String("Payment Processing Fee"),
		Amount:   stripe.Int64(int64(orderDetails.PaymentProcessingFee)),
		Currency: stripe.String("usd"),
		Quantity: stripe.Int64(int64(1)),
	})

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Params: stripe.Params{
			IdempotencyKey: stripe.String(stripe.NewIdempotencyKey()),
		},
		CustomerEmail: stripe.String(orderDetails.BillingAddress.Email),
		SuccessURL:    stripe.String(fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}&order_id=%s", spg.SuccessCallback, orderDetails.ID)),
		CancelURL:     stripe.String(fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}&order_id=%s", spg.FailureCallback, orderDetails.ID)),
	}

	ss, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return &PaymentGatewayResponse{
		Nonce: ss.ID,
	}, nil
}

func (spg *stripePaymentGateway) GetClientToken() (string, error) {
	return "", nil
}
