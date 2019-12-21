package payment_gateways

import (
	"fmt"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

const (
	StripePaymentGatewayName = "stripe"
)

type stripePaymentGateway struct {
	SecretKey       string
	SuccessCallback string
	FailureCallback string
	PublicKey       string
}

func NewStripePaymentGateway(cfg map[string]interface{}) (*stripePaymentGateway, error) {
	return &stripePaymentGateway{
		SecretKey:       cfg["secret_key"].(string),
		SuccessCallback: cfg["success_callback"].(string),
		FailureCallback: cfg["failure_callback"].(string),
		PublicKey:       cfg["public_key"].(string),
	}, nil
}

func (spg *stripePaymentGateway) GetName() string {
	return StripePaymentGatewayName
}

func (spg *stripePaymentGateway) Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error) {
	stripe.Key = spg.SecretKey

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, op := range orderDetails.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(op.Name),
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

	successUrl := fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}&status=success", spg.SuccessCallback)
	failureUrl := fmt.Sprintf("%s?session_id={CHECKOUT_SESSION_ID}&status=failure", spg.FailureCallback)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Params: stripe.Params{
			IdempotencyKey: stripe.String(stripe.NewIdempotencyKey()),
		},
		CustomerEmail:     stripe.String(orderDetails.BillingEmail),
		SuccessURL:        stripe.String(fmt.Sprintf(successUrl, orderDetails.ID)),
		CancelURL:         stripe.String(fmt.Sprintf(failureUrl, orderDetails.ID)),
		ClientReferenceID: stripe.String(orderDetails.ID),
	}

	ss, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return &PaymentGatewayResponse{
		Result: ss.PaymentIntent.ID,
		Nonce:  ss.ID,
	}, nil
}

func (spg *stripePaymentGateway) GetConfig() (map[string]interface{}, error) {
	cfg := map[string]interface{}{
		"success_callback_url": spg.SuccessCallback,
		"failure_callback_url": spg.FailureCallback,
		"public_key":           spg.PublicKey,
	}

	return cfg, nil
}
