package models

import "fmt"

type CollectionOfProduct struct {
	ProductID    string `json:"product_id" gorm:"column:product_id;primary_key"`
	CollectionID string `json:"collection_id" gorm:"column:collection_id;primary_key"`
}

func (cop *CollectionOfProduct) TableName() string {
	return "collection_of_products"
}

func (cop *CollectionOfProduct) ForeignKeys() []string {
	p := Product{}
	c := Collection{}

	return []string{
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
		fmt.Sprintf("collection_id;%s(id);RESTRICT;RESTRICT", c.TableName()),
	}
}
