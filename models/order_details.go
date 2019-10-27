package models

import (
	"time"
)

type OrderDetails struct {
	ID                   string           `json:"id"`
	Hash                 string           `json:"hash"`
	UserID               string           `json:"user_id"`
	Products             []OrderedProduct `json:"products"`
	StoreID              string           `json:"store_id"`
	ShippingAddress      *Address         `json:"shipping_address,omitempty"`
	BillingAddress       Address          `json:"billing_address"`
	PaymentMethod        PaymentMethod    `json:"payment_method"`
	ShippingMethod       *ShippingMethod  `json:"shipping_method,omitempty"`
	TotalVat             int              `json:"total_vat"`
	TotalTax             int              `json:"total_tax"`
	ShippingCharge       int              `json:"shipping_charge"`
	PaymentProcessingFee int              `json:"payment_processing_fee"`
	SubTotal             int              `json:"sub_total"`
	PaymentGateway       string           `json:"payment_gateway"`
	Nonce                string           `json:"nonce,omitempty"`
	GrandTotal           int              `json:"grand_total"`
	IsPaid               bool             `json:"is_paid"`
	Status               OrderStatus      `json:"status"`
	PaidAt               *time.Time       `json:"paid_at,omitempty"`
	ConfirmedAt          *time.Time       `json:"confirmed_at,omitempty"`
	CompletedAt          *time.Time       `json:"completed_at,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
}
