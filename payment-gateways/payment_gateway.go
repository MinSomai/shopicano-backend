package payment_gateways

import (
	"github.com/braintree-go/braintree-go"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/models"
)

type PaymentGateway interface {
	GetName() string
	GetClientToken() (string, error)
	Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error)
}

type PaymentGatewayResponse struct {
	Result                     string
	BrainTreeTransactionStatus braintree.TransactionStatus
}

var activePaymentGateway PaymentGateway

func SetActivePaymentGateway(cfg config.PaymentGatewayCfg) error {
	if cfg.Name == "stripe" {
		stripe, err := NewStripePaymentGateway(cfg.Configs)
		if err != nil {
			return err
		}
		activePaymentGateway = stripe
	} else if cfg.Name == "brainTree" {
		bt, err := NewBrainTreePaymentGateway(cfg.Configs)
		if err != nil {
			return err
		}
		activePaymentGateway = bt
	}
	return nil
}

func GetActivePaymentGateway() PaymentGateway {
	return activePaymentGateway
}
