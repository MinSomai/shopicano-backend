package models

import "time"

type OrderDetailsViewExternal struct {
	ID                     string                    `json:"id,omitempty"`
	Hash                   string                    `json:"hash"`
	ShippingCharge         int64                     `json:"shipping_charge"`
	PaymentProcessingFee   int64                     `json:"payment_processing_fee"`
	SubTotal               int64                     `json:"sub_total"`
	PaymentGateway         string                    `json:"payment_gateway"`
	Nonce                  *string                   `json:"nonce,omitempty"`
	TransactionID          *string                   `json:"transaction_id,omitempty"`
	GrandTotal             int64                     `json:"grand_total"`
	Status                 OrderStatus               `json:"status"`
	PaymentStatus          PaymentStatus             `json:"payment_status"`
	CreatedAt              *time.Time                `json:"created_at"`
	UpdatedAt              *time.Time                `json:"updated_at"`
	ShippingID             *string                   `json:"shipping_id,omitempty"`
	ShippingName           *string                   `json:"shipping_name,omitempty"`
	ShippingAddress        *string                   `json:"shipping_address,omitempty"`
	ShippingCity           *string                   `json:"shipping_city,omitempty"`
	ShippingCountry        *string                   `json:"shipping_country,omitempty"`
	ShippingPostcode       *string                   `json:"shipping_postcode,omitempty"`
	ShippingEmail          *string                   `json:"shipping_email,omitempty"`
	ShippingPhone          *string                   `json:"shipping_phone,omitempty"`
	BillingID              string                    `json:"billing_id"`
	BillingName            string                    `json:"billing_name"`
	BillingAddress         string                    `json:"billing_address"`
	BillingCity            string                    `json:"billing_city"`
	BillingCountry         string                    `json:"billing_country"`
	BillingPostcode        string                    `json:"billing_postcode"`
	BillingEmail           string                    `json:"billing_email"`
	BillingPhone           string                    `json:"billing_phone"`
	StoreID                string                    `json:"store_id"`
	StoreName              string                    `json:"store_name"`
	StoreAddress           string                    `json:"store_address"`
	StoreCity              string                    `json:"store_city"`
	StoreCountry           string                    `json:"store_country"`
	StorePostcode          string                    `json:"store_postcode"`
	StoreEmail             string                    `json:"store_email"`
	StorePhone             string                    `json:"store_phone"`
	StoreStatus            string                    `json:"store_status"`
	PaymentMethodID        string                    `json:"payment_method_id"`
	PaymentMethodName      string                    `json:"payment_method_name"`
	PaymentMethodIsOffline bool                      `json:"payment_method_is_offline"`
	Items                  []OrderedItemViewExternal `json:"items"`
	UserID                 string                    `json:"user_id"`
	UserName               string                    `json:"user_name"`
	UserEmail              string                    `json:"user_email"`
	UserPhone              *string                   `json:"user_phone,omitempty"`
	UserPicture            *string                   `json:"user_picture,omitempty"`
}

func (odi *OrderDetailsViewExternal) TableName() string {
	od := OrderDetailsView{}
	return od.TableName()
}

type OrderedItemDetailsInternal struct {
	OrderID     string `json:"order_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	SubTotal    int    `json:"sub_total"`
}
