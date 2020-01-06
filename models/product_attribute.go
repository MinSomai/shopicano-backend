package models

import "fmt"

type ProductAttribute struct {
	ProductID string `json:"-" gorm:"column:product_id;primary_key"`
	Key       string `json:"key" gorm:"column:key;primary_key"`
	Value     string `json:"value" gorm:"column:value"`
}

func (pa *ProductAttribute) TableName() string {
	return "product_attributes"
}

func (pa *ProductAttribute) ForeignKeys() []string {
	p := Product{}

	return []string{
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
	}
}
