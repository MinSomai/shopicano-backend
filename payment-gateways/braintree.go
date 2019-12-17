package payment_gateways

import (
	"context"
	"fmt"
	"github.com/braintree-go/braintree-go"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"github.com/shopicano/shopicano-backend/utils"
)

const (
	Sale BrainTreeTransactionType = "sale"
)

type BrainTreeTransactionType string

type brainTreePaymentGateway struct {
	SuccessCallback string
	FailureCallback string
	Token           string
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
	}, nil
}

func (bt *brainTreePaymentGateway) GetName() string {
	return "brainTree"
}

func (bt *brainTreePaymentGateway) Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error) {
	var items []*braintree.TransactionLineItemRequest

	for _, op := range orderDetails.Items {
		unitPrice, err := utils.IntToDecimal(op.Price, 1)
		if err != nil {
			return nil, err
		}
		TotalPrice, err := utils.IntToDecimal(op.Price*op.Quantity, 1)
		if err != nil {
			return nil, err
		}

		log.Log().Infoln(unitPrice)

		items = append(items, &braintree.TransactionLineItemRequest{
			Name:        op.Name,
			UnitAmount:  unitPrice,
			Description: op.Description,
			ProductCode: op.SKU,
			Quantity:    braintree.NewDecimal(int64(op.Quantity), 0),
			TotalAmount: TotalPrice,
			Kind:        braintree.TransactionLineItemKindDebit,
		})
	}

	paymentProcessingFee, err := utils.IntToDecimal(orderDetails.PaymentProcessingFee+3, 1)
	if err != nil {
		return nil, err
	}
	items = append(items, &braintree.TransactionLineItemRequest{
		Name:        "Payment Processing Fee",
		UnitAmount:  paymentProcessingFee,
		Quantity:    braintree.NewDecimal(int64(1), 0),
		TotalAmount: paymentProcessingFee,
		Kind:        braintree.TransactionLineItemKindDebit,
	})

	log.Log().Infoln(orderDetails.GrandTotal + orderDetails.PaymentProcessingFee)

	TotalAmount, err := utils.IntToDecimal(orderDetails.GrandTotal+orderDetails.PaymentProcessingFee, 1)
	if err != nil {
		return nil, err
	}

	resp, err := bt.client.Transaction().Create(context.Background(), &braintree.TransactionRequest{
		PaymentMethodNonce: *orderDetails.Nonce,
		Amount:             TotalAmount,
		LineItems:          items,
		BillingAddress: &braintree.Address{
			StreetAddress: fmt.Sprintf("%s, %s",
				orderDetails.BillingHouse,
				orderDetails.BillingRoad),
			Region:      orderDetails.BillingCity,
			PostalCode:  orderDetails.BillingPostcode,
			CountryName: orderDetails.BillingCountry,
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

func (bt *brainTreePaymentGateway) GetClientToken() (string, error) {
	return bt.client.ClientToken().Generate(context.Background())
}
