package payment_gateways

import (
	"context"
	"fmt"
	"github.com/braintree-go/braintree-go"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
)

const (
	BrainTreePaymentGatewayName = "brain_tree"
)

const (
	Sale BrainTreeTransactionType = "sale"
)

type BrainTreeTransactionType string

type brainTreePaymentGateway struct {
	SuccessCallback string
	FailureCallback string
	Token           string
	PublicKey       string
	PrivateKey      string
	client          *braintree.Braintree
}

func NewBrainTreePaymentGateway(cfg map[string]interface{}) (*brainTreePaymentGateway, error) {
	publicKey := cfg["public_key"].(string)
	privateKey := cfg["private_key"].(string)
	merchantID := cfg["merchant_id"].(string)

	c := braintree.New(braintree.Sandbox, merchantID, publicKey, privateKey)

	return &brainTreePaymentGateway{
		client:          c,
		SuccessCallback: cfg["success_callback"].(string),
		FailureCallback: cfg["failure_callback"].(string),
		Token:           cfg["token"].(string),
		PublicKey:       publicKey,
		PrivateKey:      privateKey,
	}, nil
}

func (bt *brainTreePaymentGateway) GetName() string {
	return BrainTreePaymentGatewayName
}

func (bt *brainTreePaymentGateway) Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error) {
	var items []*braintree.TransactionLineItemRequest

	items = append(items, &braintree.TransactionLineItemRequest{
		Name:        fmt.Sprintf("Payment for Order #%s", orderDetails.Hash),
		UnitAmount:  braintree.NewDecimal(int64(orderDetails.GrandTotal), 0),
		Quantity:    braintree.NewDecimal(int64(1), 0),
		TotalAmount: braintree.NewDecimal(int64(orderDetails.GrandTotal), 0),
		Kind:        braintree.TransactionLineItemKindDebit,
	})

	resp, err := bt.client.Transaction().Create(context.Background(), &braintree.TransactionRequest{
		PaymentMethodNonce: *orderDetails.Nonce,
		Amount:             braintree.NewDecimal(int64(orderDetails.GrandTotal), 0),
		LineItems:          items,
		BillingAddress: &braintree.Address{
			StreetAddress: fmt.Sprintf("%s", orderDetails.BillingAddress),
			Region:        orderDetails.BillingCity,
			PostalCode:    orderDetails.BillingPostcode,
			CountryName:   orderDetails.BillingCountry,
		},
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
		},
		Type: string(Sale),
	})

	if err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	return &PaymentGatewayResponse{
		Result:                     resp.Id,
		BrainTreeTransactionStatus: resp.Status,
	}, nil
}

func (bt *brainTreePaymentGateway) GetConfig() (map[string]interface{}, error) {
	token, err := bt.client.ClientToken().Generate(context.Background())
	if err != nil {
		return nil, err
	}

	cfg := map[string]interface{}{
		"client_token":         token,
		"success_callback_url": bt.SuccessCallback,
		"failure_callback_url": bt.FailureCallback,
		"public_key":           bt.PublicKey,
	}

	return cfg, nil
}
