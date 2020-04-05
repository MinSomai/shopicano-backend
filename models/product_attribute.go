package models

import "fmt"

type ProductKV struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Image string `json:"image"`
}

type OrderItemAttributeKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductAttribute struct {
	ID        string `json:"id" gorm:"column:id;unique_index;not null"`
	ProductID string `json:"-" gorm:"column:product_id;primary_key"`
	Key       string `json:"key" gorm:"column:key;primary_key"`
	Value     string `json:"value" gorm:"column:value;primary_key"`
	Image     string `json:"image" gorm:"column:image"`
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
