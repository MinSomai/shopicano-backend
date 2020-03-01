package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type OrderedItemView struct {
	ID               string                 `json:"id"`
	OrderID          string                 `json:"order_id"`
	ProductID        string                 `json:"product_id"`
	Name             string                 `json:"name"`
	Quantity         int                    `json:"quantity"`
	Price            int64                  `json:"price"`
	ProductCost      int64                  `json:"product_cost"`
	SubTotal         int64                  `json:"sub_total"`
	Description      string                 `json:"description"`
	SKU              string                 `json:"sku"`
	AdditionalImages string                 `json:"additional_images"`
	Image            string                 `json:"image"`
	IsShippable      bool                   `json:"is_shippable"`
	IsDigital        bool                   `json:"is_digital"`
	Attributes       []OrderItemAttributeKV `json:"attributes"`
}

func (oiv *OrderedItemView) TableName() string {
	return "ordered_item_views"
}

func (oiv *OrderedItemView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT oi.id AS id, oi.order_id AS order_id, oi.product_id AS product_id, p.name AS name,"+
		" oi.quantity AS quantity, oi.price AS price, oi.product_cost AS product_cost, oi.sub_total AS sub_total,"+
		" p.description AS description, p.sku AS sku, p.additional_images AS additional_images, p.image AS image,"+
		" p.is_shippable AS is_shippable, p.is_digital AS is_digital, p.digital_download_link AS digital_download_link"+
		" FROM ordered_items AS oi"+
		" LEFT JOIN products AS p ON oi.product_id = p.id;", oiv.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (oiv *OrderedItemView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", oiv.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
