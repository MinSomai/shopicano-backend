package models

import "fmt"

type ShippingForLocation struct {
	LocationID       int64  `json:"location_id" gorm:"column:location_id;primary_key"`
	ShippingMethodID string `json:"shipping_method_id" gorm:"column:shipping_method_id;primary_key"`
}

func (m *ShippingForLocation) TableName() string {
	return "shipping_for_locations"
}

func (m *ShippingForLocation) ForeignKeys() []string {
	sm := ShippingMethod{}
	l := Location{}

	return []string{
		fmt.Sprintf("shipping_method_id;%s(id);RESTRICT;RESTRICT", sm.TableName()),
		fmt.Sprintf("location_id;%s(id);RESTRICT;RESTRICT", l.TableName()),
	}
}
