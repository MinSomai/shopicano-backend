package payment_gateways

import (
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/models"
)

type PaymentGateway interface {
	GetName() string
	Pay(orderDetails *models.OrderDetails) (*PaymentGatewayResponse, error)
}

type PaymentGatewayResponse struct {
	ReferenceID string
}

var activePaymentGateway PaymentGateway

func SetActivePaymentGateway(cfg config.PaymentGatewayCfg) error {
	if cfg.Name == "stripe" {
		stripe, err := NewStripePaymentGateway(cfg.Configs)
		if err != nil {
			return err
		}
		activePaymentGateway = stripe
	}
	return nil
}

func GetActivePaymentGateway() PaymentGateway {
	return activePaymentGateway
}
