package models

import "fmt"

type AdditionalChargeOfProduct struct {
	ProductID          string `json:"product_id" gorm:"column:product_id;primary_key"`
	AdditionalChargeID string `json:"additional_charge_id" gorm:"column:additional_charge_id;primary_key"`
}

func (acp *AdditionalChargeOfProduct) TableName() string {
	return "additional_charges_of_products"
}

func (acp *AdditionalChargeOfProduct) ForeignKeys() []string {
	p := Product{}
	ac := AdditionalCharge{}

	return []string{
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
		fmt.Sprintf("additional_charge_id;%s(id);RESTRICT;RESTRICT", ac.TableName()),
	}
}
