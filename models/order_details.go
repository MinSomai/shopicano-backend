package models

import (
	"time"
)

type OrderDetails struct {
	ID                        string
	Hash                      string
	UserID                    string
	Products                  []OrderedProduct
	StoreID                   string
	ShippingAddressID         *Address
	BillingAddressID          Address
	PaymentMethodID           PaymentMethod
	ShippingMethodID          *ShippingMethod
	TotalVat                  int
	TotalTax                  int
	ShippingCharge            int
	PaymentProcessingFee      int
	SubTotal                  int
	PaymentGateway            string
	PaymentGatewayReferenceID string `json:"payment_gateway_reference_id"`
	GrandTotal                int
	IsPaid                    bool
	Status                    OrderStatus
	PaidAt                    *time.Time
	ConfirmedAt               *time.Time
	CompletedAt               *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
