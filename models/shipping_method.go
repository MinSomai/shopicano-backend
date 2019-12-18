package models

import "time"

type ShippingMethod struct {
	ID                      string    `json:"id" sql:"id" gorm:"primary_key"`
	Name                    string    `json:"name" sql:"name" gorm:"unique;not null"`
	ApproximateDeliveryTime int       `json:"approximate_delivery_time" gorm:"approximate_delivery_time" gorm:"index"`
	DeliveryCharge          int       `json:"delivery_charge" sql:"delivery_charge" gorm:"index"`
	WeightUnit              string    `json:"weight_unit" sql:"weight_unit"`
	IsPublished             bool      `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt               time.Time `json:"created_at" sql:"created_at" gorm:"not null;index"`
	UpdatedAt               time.Time `json:"updated_at" sql:"updated_at" gorm:"not null"`
}

func (sm *ShippingMethod) TableName() string {
	return "shipping_methods"
}

func (sm *ShippingMethod) CalculateDeliveryCharge(weight int) int {
	// TODO : Currently weight doesn't have any impact on charge
	return sm.DeliveryCharge
}
