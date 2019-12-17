package models

import "time"

type OrderDetails struct {
	ID                     string      `json:"id"`
	Hash                   string      `json:"hash"`
	TotalVat               int64       `json:"total_vat"`
	TotalTax               int64       `json:"total_tax"`
	ShippingCharge         int64       `json:"shipping_charge"`
	PaymentProcessingFee   int64       `json:"payment_processing_fee"`
	SubTotal               int64       `json:"sub_total"`
	PaymentGateway         string      `json:"payment_gateway"`
	Nonce                  string      `json:"nonce"`
	TransactionID          string      `json:"transaction_id"`
	GrandTotal             int64       `json:"grand_total"`
	IsPaid                 bool        `json:"is_paid"`
	Status                 OrderStatus `json:"status"`
	PaidAt                 *time.Time  `json:"paid_at"`
	ConfirmedAt            *time.Time  `json:"confirmed_at"`
	CompletedAt            *time.Time  `json:"completed_at"`
	CreatedAt              *time.Time  `json:"created_at"`
	UpdatedAt              *time.Time  `json:"updated_at"`
	ShippingID             string      `json:"shipping_id"`
	ShippingName           string      `json:"shipping_name"`
	ShippingHouse          string      `json:"shipping_house"`
	ShippingRoad           string      `json:"shipping_road"`
	ShippingCity           string      `json:"shipping_city"`
	ShippingCountry        string      `json:"shipping_country"`
	ShippingPostcode       string      `json:"shipping_postcode"`
	ShippingEmail          string      `json:"shipping_email"`
	ShippingPhone          string      `json:"shipping_phone"`
	BillingID              string      `json:"billing_id"`
	BillingName            string      `json:"billing_name"`
	BillingHouse           string      `json:"billing_house"`
	BillingRoad            string      `json:"billing_road"`
	BillingCity            string      `json:"billing_city"`
	BillingCountry         string      `json:"billing_country"`
	BillingPostcode        string      `json:"billing_postcode"`
	BillingEmail           string      `json:"billing_email"`
	BillingPhone           string      `json:"billing_phone"`
	StoreID                string      `json:"store_id"`
	StoreName              string      `json:"store_name"`
	StoreAddress           string      `json:"store_address"`
	StoreCity              string      `json:"store_city"`
	StoreCountry           string      `json:"store_country"`
	StorePostcode          string      `json:"store_postcode"`
	StoreEmail             string      `json:"store_email"`
	StorePhone             string      `json:"store_phone"`
	StoreStatus            string      `json:"store_status"`
	PaymentMethodID        string      `json:"payment_method_id"`
	PaymentMethodName      string      `json:"payment_method_name"`
	PaymentMethodIsOffline bool        `json:"payment_method_is_offline"`
}

func (od *OrderDetails) TableName() string {
	return "order_details_views"
}

type OrderedItemDetails struct {
	OrderID     string `json:"order_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	Vat         int    `json:"vat"`
	Tax         int    `json:"tax"`
	Total       int    `json:"total"`
}
