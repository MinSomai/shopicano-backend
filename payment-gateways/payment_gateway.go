package payment_gateways

import (
	"errors"
	"github.com/braintree-go/braintree-go"
	"github.com/shopicano/shopicano-backend/config"
	"github.com/shopicano/shopicano-backend/models"
)

type PaymentGateway interface {
	GetName() string
	GetConfig() (map[string]interface{}, error)
	Pay(orderDetails *models.OrderDetailsView) (*PaymentGatewayResponse, error)
}

type PaymentGatewayResponse struct {
	Result                     string
	Nonce                      string
	BrainTreeTransactionStatus braintree.TransactionStatus
}

var activePaymentGateway PaymentGateway

func SetActivePaymentGateway(cfg config.PaymentGatewayCfg) error {
	if cfg.Name == StripePaymentGatewayName {
		stripe, err := NewStripePaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return err
		}
		activePaymentGateway = stripe
	} else if cfg.Name == BrainTreePaymentGatewayName {
		bt, err := NewBrainTreePaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return err
		}
		activePaymentGateway = bt
	} else if cfg.Name == TwoCheckoutPaymentGatewayName {
		tco, err := NewTwoCheckoutPaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return err
		}
		activePaymentGateway = tco
	}
	return nil
}

func GetActivePaymentGateway() PaymentGateway {
	return activePaymentGateway
}

func GetPaymentGatewayByName(name string) (PaymentGateway, error) {
	cfg := config.PaymentGateway()
	if name == StripePaymentGatewayName {
		stripe, err := NewStripePaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return stripe, nil
	} else if name == BrainTreePaymentGatewayName {
		bt, err := NewBrainTreePaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return bt, nil
	} else if name == TwoCheckoutPaymentGatewayName {
		tco, err := NewTwoCheckoutPaymentGateway(cfg.Configs[cfg.Name].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return tco, nil
	}
	return nil, errors.New("payment gateway not found")
}
