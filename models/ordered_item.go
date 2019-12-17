package models

import "fmt"

type OrderedItem struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price" sql:"price"`
	TotalVat  int    `json:"total_vat" sql:"total_vat"`
	TotalTax  int    `json:"total_tax" sql:"total_tax"`
	SubTotal  int    `json:"sub_total" sql:"sub_total"`
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
