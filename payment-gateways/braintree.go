package payment_gateways

import (
	"context"
	"github.com/braintree-go/braintree-go"
	"github.com/shopicano/shopicano-backend/models"
)

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

func (bt *brainTreePaymentGateway) Pay(orderDetails *models.OrderDetailsInternal) (*PaymentGatewayResponse, error) {
	//var items []*braintree.TransactionLineItemRequest
	//
	//for _, op := range orderDetails.Items {
	//	unitPrice, err := utils.IntToDecimal(op.Price, 100)
	//	if err != nil {
	//		return nil, err
	//	}
	//	TotalPrice, err := utils.IntToDecimal(op.Price*op.Quantity, 100)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	log.Log().Infoln(unitPrice)
	//
	//	items = append(items, &braintree.TransactionLineItemRequest{
	//		//Name:        op.Name,
	//		UnitAmount:  unitPrice,
	//		Description: op.ProductID,
	//		Quantity:    braintree.NewDecimal(int64(op.Quantity), 0),
	//		TotalAmount: TotalPrice,
	//		Kind:        braintree.TransactionLineItemKindDebit,
	//	})
	//}
	//
	//TotalPrice, err := utils.IntToDecimal(orderDetails.PaymentProcessingFee, 100)
	//if err != nil {
	//	return nil, err
	//}
	//items = append(items, &braintree.TransactionLineItemRequest{
	//	Name:        "Payment Processing Fee",
	//	UnitAmount:  TotalPrice,
	//	Quantity:    braintree.NewDecimal(int64(1), 0),
	//	TotalAmount: TotalPrice,
	//	Kind:        braintree.TransactionLineItemKindDebit,
	//})
	//
	//TotalAmount, err := utils.IntToDecimal(orderDetails.GrandTotal+orderDetails.PaymentProcessingFee, 100)
	//if err != nil {
	//	return nil, err
	//}
	//resp, err := bt.client.Transaction().Create(context.Background(), &braintree.TransactionRequest{
	//	PaymentMethodNonce: orderDetails.Nonce,
	//	Amount:             TotalAmount,
	//	LineItems:          items,
	//	BillingAddress: &braintree.Address{
	//		StreetAddress: fmt.Sprintf("%s, %s",
	//			orderDetails.BillingAddress.House,
	//			orderDetails.BillingAddress.Road),
	//		Region:      orderDetails.BillingAddress.City,
	//		PostalCode:  orderDetails.BillingAddress.Postcode,
	//		CountryName: orderDetails.BillingAddress.Country,
	//	},
	//	Options: &braintree.TransactionOptions{
	//		SubmitForSettlement: true,
	//	},
	//	Type: "sale",
	//})
	//
	//if err != nil {
	//	log.Log().Errorln(err)
	//	return nil, err
	//}
	//
	//return &PaymentGatewayResponse{
	//	Nonce:                      resp.Id,
	//	BrainTreeTransactionStatus: resp.Status,
	//}, nil
	return &PaymentGatewayResponse{}, nil
}

func (bt *brainTreePaymentGateway) GetClientToken() (string, error) {
	return bt.client.ClientToken().Generate(context.Background())
}
