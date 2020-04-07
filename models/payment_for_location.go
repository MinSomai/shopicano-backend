package models

import "fmt"

type PaymentForLocation struct {
	LocationID      int64  `json:"location_id" gorm:"column:location_id;primary_key"`
	PaymentMethodID string `json:"payment_method_id" gorm:"column:payment_method_id;primary_key"`
}

func (m *PaymentForLocation) TableName() string {
	return "payment_for_locations"
}

func (m *PaymentForLocation) ForeignKeys() []string {
	pm := PaymentMethod{}
	l := Location{}

	return []string{
		fmt.Sprintf("payment_method_id;%s(id);RESTRICT;RESTRICT", pm.TableName()),
		fmt.Sprintf("location_id;%s(id);RESTRICT;RESTRICT", l.TableName()),
	}
}
