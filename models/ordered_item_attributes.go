package models

import "fmt"

type OrderedItemAttribute struct {
	OrderedItemID  string `json:"ordered_item_id" gorm:"column:ordered_item_id;primary_key"`
	AttributeKey   string `json:"attribute_key" gorm:"column:attribute_key;primary_key"`
	AttributeValue string `json:"attribute_value" gorm:"column:attribute_value"`
}

func (oia *OrderedItemAttribute) TableName() string {
	return "ordered_item_attributes"
}

func (oia *OrderedItemAttribute) ForeignKeys() []string {
	o := OrderedItem{}

	return []string{
		fmt.Sprintf("ordered_item_id;%s(id);RESTRICT;RESTRICT", o.TableName()),
	}
}
