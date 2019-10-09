package models

import (
	"fmt"
	"time"
)

const (
	Pending   OrderStatus = "pending"
	Confirmed OrderStatus = "confirmed"
	Shipping  OrderStatus = "shipping"
	Cancelled OrderStatus = "cancelled"
	Delivered OrderStatus = "delivered"
)

type OrderStatus string

type Order struct {
	ID                        string      `json:"id" sql:"id" gorm:"primary_key"`
	Hash                      string      `json:"hash" json:"hash" gorm:"unique;not null"`
	UserID                    string      `json:"user_id" sql:"user_id" gorm:"index;not null"`
	StoreID                   string      `json:"store_id" sql:"store_id" gorm:"index;not null"`
	ShippingAddressID         *string     `json:"shipping_address_id;omitempty" sql:"shipping_address_id"`
	BillingAddressID          string      `json:"billing_address_id" sql:"billing_address_id;not null"`
	PaymentMethodID           string      `json:"payment_method_id" sql:"payment_method_id;not null"`
	ShippingMethodID          *string     `json:"shipping_method_id;omitempty" json:"shipping_method_id"`
	TotalVat                  int         `json:"total_vat" sql:"total_vat"`
	TotalTax                  int         `json:"total_tax" sql:"total_tax"`
	ShippingCharge            int         `json:"shipping_charge" sql:"shipping_charge"`
	PaymentProcessingFee      int         `json:"payment_processing_fee" sql:"payment_processing_fee"`
	SubTotal                  int         `json:"sub_total" sql:"sub_total"`
	PaymentGateway            string      `json:"payment_gateway" sql:"payment_gateway"`
	PaymentGatewayReferenceID string      `json:"payment_gateway_reference_id" sql:"payment_gateway_reference_id"`
	GrandTotal                int         `json:"grand_total" sql:"grand_total"`
	IsPaid                    bool        `json:"is_paid" sql:"is_paid"`
	Status                    OrderStatus `json:"status" sql:"status"`
	PaidAt                    *time.Time  `json:"paid_at;omitempty" sql:"paid_at" gorm:"index"`
	ConfirmedAt               *time.Time  `json:"confirmed_at;omitempty" sql:"confirmed_at" gorm:"index"`
	CompletedAt               *time.Time  `json:"completed_at;omitempty" sql:"completed_at" gorm:"index"`
	CreatedAt                 time.Time   `json:"created_at" sql:"created_at" gorm:"index;not null"`
	UpdatedAt                 time.Time   `json:"updated_at" sql:"updated_at"`
}

func (o *Order) TableName() string {
	return "orders"
}

func (o *Order) ForeignKeys() []string {
	s := Store{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
	}
}

type OrderedProduct struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price" sql:"price"`
	TotalVat  int    `json:"total_vat" sql:"total_vat"`
	TotalTax  int    `json:"total_tax" sql:"total_tax"`
	SubTotal  int    `json:"sub_total" sql:"sub_total"`
	//CurrencyID string `json:"currency_id" sql:"currency_id"`
}

func (op *OrderedProduct) TableName() string {
	return "ordered_products"
}

func (op *OrderedProduct) ForeignKeys() []string {
	o := Order{}
	p := Product{}

	return []string{
		fmt.Sprintf("order_id;%s(id);RESTRICT;RESTRICT", o.TableName()),
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
	}
}
