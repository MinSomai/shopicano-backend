package models

import "time"

type DiscountType string

const (
	ProductDiscount  DiscountType = "product_discount"
	ShippingDiscount DiscountType = "shipping_discount"
	TotalDiscount    DiscountType = "total_discount"
)

type Discount struct {
	ID             string       `json:"id" gorm:"column:id;primary_key"`
	Code           string       `json:"code" gorm:"column:code;unique_index"`
	IsActive       bool         `json:"is_active" gorm:"column:is_active;index"`
	DiscountAmount int          `json:"discount_amount" gorm:"column:discount_amount"`
	IsFlatDiscount bool         `json:"is_flat_discount" gorm:"column:is_flat_discount"`
	MaxDiscount    int          `json:"max_discount" gorm:"column:max_discount"`
	MaxUsage       int          `json:"max_usage" gorm:"column:max_usage"`
	DiscountType   DiscountType `json:"discount_type" gorm:"column:discount_type;index"`
	StartAt        time.Time    `json:"start_at" gorm:"column:start_at;index"`
	EndAt          time.Time    `json:"end_at" gorm:"column:end_at;index"`
	CreatedAt      time.Time    `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"column:updated_at;index"`
}

func (d *Discount) TableName() string {
	return "discounts"
}
