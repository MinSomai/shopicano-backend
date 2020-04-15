package models

import (
	"fmt"
	"time"
)

type PaymentGateway struct {
	ID          string    `json:"id" gorm:"column:id;primary_key"`
	DisplayName string    `json:"display_name" gorm:"column:display_name"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (pg *PaymentGateway) TableName() string {
	return "payment_gateways"
}

type PaymentGatewayConfig struct {
	PaymentGatewayID string `json:"payment_gateway_id" gorm:"column:payment_gateway_id;index"`
	Key              string `json:"key" gorm:"column:key"`
	Value            string `json:"value" gorm:"column:value"`
	ValueType        string `json:"value_type" gorm:"column:value_type"`
}

func (pgc *PaymentGatewayConfig) TableName() string {
	return "payment_gateway_configs"
}

func (pgc *PaymentGatewayConfig) ForeignKeys() []string {
	pg := PaymentGateway{}

	return []string{
		fmt.Sprintf("payment_gateway_id;%s(id);RESTRICT;RESTRICT", pg.TableName()),
	}
}
