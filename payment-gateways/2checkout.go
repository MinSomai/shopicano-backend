package payment_gateways

import (
	"fmt"
	"github.com/shopicano/shopicano-backend/models"
)

const (
	TwoCheckoutPaymentGatewayName = "2co"
)

type twoCheckoutPaymentGateway struct {
	SuccessCallback string
	FailureCallback string
	PublicKey       string
	PrivateKey      string
	MerchantCode    string
	SecretKey       string
}

func NewTwoCheckoutPaymentGateway(cfg map[string]interface{}) (*twoCheckoutPaymentGateway, error) {
	publicKey := cfg["public_key"].(string)
	privateKey := cfg["private_key"].(string)
	merchantCode := cfg["merchant_code"].(string)
	secretKey := cfg["secret_key"].(string)

	return &twoCheckoutPaymentGateway{
		SuccessCallback: cfg["success_callback"].(string),
		FailureCallback: cfg["failure_callback"].(string),
		PublicKey:       publicKey,
		PrivateKey:      privateKey,
		MerchantCode:    merchantCode,
		SecretKey:       secretKey,
	}, nil
}

func (tco *twoCheckoutPaymentGateway) GetName() string {
	return TwoCheckoutPaymentGatewayName
}

func (tco *twoCheckoutPaymentGateway) Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error) {
	url := "https://sandbox.2checkout.com/checkout/purchase"

	payload := fmt.Sprintf("sid=%s&", tco.MerchantCode)
	payload += fmt.Sprintf("mode=%s&", "2CO")
	payload += fmt.Sprintf("submit=%s&", "Checkout")
	payload += fmt.Sprintf("merchant_order_id=%s&", orderDetails.ID)
	payload += fmt.Sprintf("currency_code=%s&", "USD")
	payload += fmt.Sprintf("street_address=%s&", orderDetails.BillingAddress)
	payload += fmt.Sprintf("city=%s&", orderDetails.BillingCity)
	payload += fmt.Sprintf("state=%s&", orderDetails.BillingCity)
	payload += fmt.Sprintf("zip=%s&", orderDetails.BillingPostcode)
	payload += fmt.Sprintf("country=%s&", orderDetails.BillingCountry)
	payload += fmt.Sprintf("phone=%s&", orderDetails.BillingPhone)
	payload += fmt.Sprintf("email=%s&", orderDetails.BillingEmail)

	payload += fmt.Sprintf("li_0_type=%s&", "product")
	payload += fmt.Sprintf("li_0_name=%s&", fmt.Sprintf("Payment for Order %s", orderDetails.Hash))
	payload += fmt.Sprintf("li_0_price=%s&", fmt.Sprintf("%.2f", float64(orderDetails.GrandTotal)))
	payload += fmt.Sprintf("li_0_quantity=%s&", fmt.Sprintf("%d", 1))
	payload += fmt.Sprintf("li_0_tangible=%s&", "N")

	payload += "purchase_step=payment-method"

	return &PaymentGatewayResponse{
		Result: fmt.Sprintf("%s?%s", url, payload),
	}, nil
}

func (tco *twoCheckoutPaymentGateway) GetConfig() (map[string]interface{}, error) {
	cfg := map[string]interface{}{
		"success_callback_url": tco.SuccessCallback,
		"failure_callback_url": tco.FailureCallback,
		"public_key":           tco.PublicKey,
	}
	return cfg, nil
}

func (tco *twoCheckoutPaymentGateway) ValidateTransaction(orderDetails *models.OrderDetailsView) error {
	return nil
}

func (tco *twoCheckoutPaymentGateway) VoidTransaction(map[string]interface{}) error {
	return nil
}
