package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type OrderDetailsView struct {
	ID                      string            `json:"id,omitempty"`
	Hash                    string            `json:"hash,omitempty"`
	ShippingCharge          int64             `json:"shipping_charge,omitempty"`
	PaymentProcessingFee    int64             `json:"payment_processing_fee,omitempty"`
	SubTotal                int64             `json:"sub_total,omitempty"`
	PaymentGateway          string            `json:"payment_gateway,omitempty"`
	Nonce                   *string           `json:"nonce,omitempty"`          // Private
	TransactionID           *string           `json:"transaction_id,omitempty"` //Private
	GrandTotal              int64             `json:"grand_total,omitempty"`
	DiscountedAmount        int64             `json:"discounted_amount,omitempty"`
	CouponCode              string            `json:"coupon_code,omitempty"`
	Status                  OrderStatus       `json:"status,omitempty"`
	PaymentStatus           PaymentStatus     `json:"payment_status,omitempty"`
	IsAllDigitalProducts    bool              `json:"is_all_digital_products,omitempty"`
	CreatedAt               *time.Time        `json:"created_at,omitempty"`
	UpdatedAt               *time.Time        `json:"updated_at,omitempty"`
	ShippingID              *string           `json:"shipping_id,omitempty"`
	ShippingName            *string           `json:"shipping_name,omitempty"`
	ShippingAddress         *string           `json:"shipping_address,omitempty"`
	ShippingCity            *string           `json:"shipping_city,omitempty"`
	ShippingCountry         *string           `json:"shipping_country,omitempty"`
	ShippingPostcode        *string           `json:"shipping_postcode,omitempty"`
	ShippingEmail           *string           `json:"shipping_email,omitempty"`
	ShippingPhone           *string           `json:"shipping_phone,omitempty"`
	BillingID               string            `json:"billing_id,omitempty"`
	BillingName             string            `json:"billing_name,omitempty"`
	BillingAddress          string            `json:"billing_address,omitempty"`
	BillingCity             string            `json:"billing_city,omitempty"`
	BillingCountry          string            `json:"billing_country,omitempty"`
	BillingPostcode         string            `json:"billing_postcode,omitempty"`
	BillingEmail            string            `json:"billing_email,omitempty"`
	BillingPhone            string            `json:"billing_phone,omitempty"`
	StoreID                 string            `json:"store_id,omitempty"`
	StoreName               string            `json:"store_name,omitempty"`
	StoreAddress            string            `json:"store_address,omitempty"`
	StoreCity               string            `json:"store_city,omitempty"`
	StoreCountry            string            `json:"store_country,omitempty"`
	StorePostcode           string            `json:"store_postcode,omitempty"`
	StoreEmail              string            `json:"store_email,omitempty"`
	StorePhone              string            `json:"store_phone,omitempty"`
	StoreStatus             string            `json:"store_status,omitempty"`
	PaymentMethodID         string            `json:"payment_method_id,omitempty"`
	PaymentMethodName       string            `json:"payment_method_name,omitempty"`
	PaymentMethodIsOffline  bool              `json:"payment_method_is_offline"`
	ShippingMethodID        string            `json:"shipping_method_id"`
	ShippingMethodName      string            `json:"shipping_method_name"`
	ApproximateDeliveryTime int               `json:"approximate_delivery_time"`
	Items                   []OrderedItemView `json:"items"`
	UserID                  string            `json:"user_id"`
	UserName                string            `json:"user_name"`
	UserEmail               string            `json:"user_email"`
	UserPhone               *string           `json:"user_phone,omitempty"`
	UserPicture             *string           `json:"user_picture,omitempty"`
	ReviewRating            int               `json:"review_rating,omitempty"`
	ReviewDescription       string            `json:"review_description,omitempty"`
	SellerEarnings          int64             `json:"seller_earnings"`
	PlatformEarnings        int64             `json:"platform_earnings"`
	ActualEarnings          int64             `json:"actual_earnings"`
}

func (odv *OrderDetailsView) TableName() string {
	return "order_details_views"
}

func (odv *OrderDetailsView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT o.id AS id, o.hash AS hash, o.user_id AS user_id, o.is_all_digital_products AS is_all_digital_products,"+
		" u.name AS user_name, u.email AS user_email, u.phone AS user_phone, u.profile_picture AS user_picture,"+
		" o.shipping_charge AS shipping_charge, o.payment_processing_fee AS payment_processing_fee, o.sub_total AS sub_total,"+
		" o.payment_gateway AS payment_gateway, o.nonce AS nonce, o.transaction_id AS transaction_id, o.grand_total AS grand_total,"+
		" o.status AS status, o.payment_status AS payment_status, o.created_at AS created_at, o.discounted_amount AS discounted_amount,"+
		" o.updated_at AS updated_at, sa.id AS shipping_id, sa.name AS shipping_name, sa.address AS shipping_address,"+
		" sa.city AS shipping_city, sa.country AS shipping_country, sa.postcode AS shipping_postcode, cop.code AS coupon_code,"+
		" sa.email AS shipping_email, sa.phone AS shipping_phone, ba.id AS billing_id, ba.name AS billing_name, ba.address AS billing_address,"+
		" ba.city AS billing_city, ba.country AS billing_country, ba.postcode AS billing_postcode, ba.email AS billing_email,"+
		" ba.phone AS billing_phone, s.id AS store_id, s.name AS store_name, sta.address AS store_address, sta.city AS store_city,"+
		" sta.country AS store_country, sta.postcode AS store_postcode, sta.email AS store_email, sta.phone AS store_phone, s.status AS store_status,"+
		" sm.id AS shipping_method_id, sm.name AS shipping_method_name, sm.approximate_delivery_time AS approximate_delivery_time,"+
		" pm.id AS payment_method_id, pm.name AS payment_method_name, pm.is_offline_payment AS payment_method_is_offline,"+
		" rv.rating AS review_rating, rv.description AS review_description, o.seller_earnings AS seller_earnings,"+
		" o.platform_earnings AS platform_earnings, o.actual_earnings AS actual_earnings"+
		" FROM orders AS o"+
		" LEFT JOIN addresses_view AS sa ON o.shipping_address_id = sa.id"+
		" LEFT JOIN addresses_view AS ba ON o.billing_address_id = ba.id"+
		" LEFT JOIN stores AS s ON o.store_id = s.id"+
		" LEFT JOIN addresses_view AS sta ON s.address_id = sta.id"+
		" LEFT JOIN users AS u ON o.user_id = u.id"+
		" LEFT JOIN coupon_usages AS cu ON o.id = cu.order_id"+
		" LEFT JOIN coupons AS cop ON cu.coupon_id = cop.id"+
		" LEFT JOIN reviews AS rv ON o.id = rv.order_id"+
		" LEFT JOIN shipping_methods AS sm ON o.shipping_method_id = sm.id"+
		" LEFT JOIN payment_methods AS pm ON o.payment_method_id = pm.id;", odv.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (odv *OrderDetailsView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", odv.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
