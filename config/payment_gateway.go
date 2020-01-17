package config

import (
	"github.com/spf13/viper"
)

type PaymentGatewayCfg struct {
	Name    string
	Configs map[string]interface{}
}

var paymentGateway PaymentGatewayCfg

func PaymentGateway() PaymentGatewayCfg {
	return paymentGateway
}

func LoadPaymentGateway() {
	mu.Lock()
	defer mu.Unlock()

	paymentGateway = PaymentGatewayCfg{
		Name:    viper.GetString("payment_gateway.name"),
		Configs: viper.GetStringMap("payment_gateway"),
	}
}
