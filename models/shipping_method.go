package models

import "time"

const (
	Ounce WeightUnit = "ounce"
	Gram  WeightUnit = "gram"
)

type WeightUnit string

func (wu WeightUnit) IsValid() bool {
	for _, v := range []WeightUnit{Ounce, Gram} {
		if v == wu {
			return true
		}
	}
	return false
}

type ShippingMethod struct {
	ID                      string     `json:"id" sql:"id" gorm:"primary_key"`
	Name                    string     `json:"name" sql:"name" gorm:"unique;not null"`
	ApproximateDeliveryTime int        `json:"approximate_delivery_time" gorm:"approximate_delivery_time" gorm:"index"`
	DeliveryCharge          int64      `json:"delivery_charge" sql:"delivery_charge" gorm:"index"`
	WeightUnit              WeightUnit `json:"weight_unit" sql:"weight_unit"`
	IsFlat                  bool       `json:"is_flat" gorm:"column:is_flat"`
	IsPublished             bool       `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt               time.Time  `json:"created_at" sql:"created_at" gorm:"not null;index"`
	UpdatedAt               time.Time  `json:"updated_at" sql:"updated_at" gorm:"not null"`
}

func (sm *ShippingMethod) TableName() string {
	return "shipping_methods"
}

func (sm *ShippingMethod) CalculateDeliveryCharge(weight int) int64 {
	// TODO : Currently weight doesn't have any impact on charge
	return sm.DeliveryCharge
}
