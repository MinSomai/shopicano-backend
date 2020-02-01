package models

import (
	"fmt"
	"time"
)

const (
	OrderPending   OrderStatus = "order_pending"
	OrderCancelled OrderStatus = "order_cancelled"
	OrderConfirmed OrderStatus = "order_confirmed"
	OrderShipping  OrderStatus = "order_shipping"
	OrderDelivered OrderStatus = "order_delivered"

	PaymentPending   PaymentStatus = "payment_pending"
	PaymentCompleted PaymentStatus = "payment_completed"
	PaymentFailed    PaymentStatus = "payment_failed"
	PaymentReverted  PaymentStatus = "payment_reverted"
)

type OrderStatus string
type PaymentStatus string

type Order struct {
	ID                   string        `json:"id" gorm:"column:id;primary_key"`
	Hash                 string        `json:"hash" gorm:"column:hash;unique_index;not null"`
	UserID               string        `json:"user_id" gorm:"column:user_id;index;not null"`
	StoreID              string        `json:"store_id" gorm:"column:store_id;index;not null"`
	ShippingAddressID    *string       `json:"shipping_address_id;omitempty" gorm:"column:shipping_address_id"`
	BillingAddressID     string        `json:"billing_address_id" gorm:"column:billing_address_id;not null"`
	PaymentMethodID      string        `json:"payment_method_id" gorm:"column:payment_method_id;not null"`
	ShippingMethodID     *string       `json:"shipping_method_id;omitempty" gorm:"column:shipping_method_id"`
	ShippingCharge       int           `json:"shipping_charge" gomr:"column:shipping_charge"`
	PaymentProcessingFee int           `json:"payment_processing_fee" gorm:"column:payment_processing_fee"`
	SubTotal             int           `json:"sub_total" gorm:"column:sub_total"`
	IsAllDigitalProducts bool          `json:"is_all_digital_products" gorm:"column:is_all_digital_products;index"`
	PaymentGateway       *string       `json:"payment_gateway" gorm:"column:payment_gateway"`
	Nonce                *string       `json:"nonce" gomr:"column:nonce"`
	TransactionID        *string       `json:"transaction_id" gorm:"column:transaction_id;unique_index"`
	GrandTotal           int           `json:"grand_total" gorm:"column:grand_total"`
	DiscountedAmount     int           `json:"discounted_amount" gorm:"column:discounted_amount"`
	Status               OrderStatus   `json:"status" gorm:"column:status"`
	PaymentStatus        PaymentStatus `json:"payment_status" gorm:"column:payment_status"`
	CreatedAt            time.Time     `json:"created_at" gorm:"column:created_at;index;not null"`
	UpdatedAt            time.Time     `json:"updated_at" gorm:"column:updated_at"`
}

func (o *Order) TableName() string {
	return "orders"
}

func (o *Order) ForeignKeys() []string {
	s := Store{}
	u := User{}
	a := Address{}
	sm := ShippingMethod{}
	pm := PaymentMethod{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
		fmt.Sprintf("user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
		fmt.Sprintf("shipping_address_id;%s(id);RESTRICT;RESTRICT", a.TableName()),
		fmt.Sprintf("billing_address_id;%s(id);RESTRICT;RESTRICT", a.TableName()),
		fmt.Sprintf("payment_method_id;%s(id);RESTRICT;RESTRICT", pm.TableName()),
		fmt.Sprintf("shipping_method_id;%s(id);RESTRICT;RESTRICT", sm.TableName()),
	}
}
