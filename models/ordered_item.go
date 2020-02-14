package models

import "fmt"

type OrderedItem struct {
	OrderID     string `json:"order_id" gorm:"column:order_id"`
	ProductID   string `json:"product_id" gorm:"column:product_id"`
	Quantity    int    `json:"quantity" gorm:"column:quantity"`
	Price       int    `json:"price" gorm:"column:price"`
	ProductCost int    `json:"product_cost" gorm:"column:product_cost"`
	SubTotal    int    `json:"sub_total" gorm:"column:sub_total"`
}

func (op *OrderedItem) TableName() string {
	return "ordered_items"
}

func (op *OrderedItem) ForeignKeys() []string {
	o := Order{}
	p := Product{}

	return []string{
		fmt.Sprintf("order_id;%s(id);RESTRICT;RESTRICT", o.TableName()),
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
	}
}
