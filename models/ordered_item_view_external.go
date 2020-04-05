package models

type OrderedItemViewExternal struct {
	ID               string                 `json:"id"`
	OrderID          string                 `json:"order_id"`
	ProductID        string                 `json:"product_id"`
	Name             string                 `json:"name"`
	Quantity         int                    `json:"quantity"`
	Price            int64                  `json:"price"`
	SubTotal         int64                  `json:"sub_total"`
	Description      string                 `json:"description"`
	SKU              string                 `json:"sku"`
	AdditionalImages []string               `json:"additional_images"`
	Image            string                 `json:"image"`
	IsShippable      bool                   `json:"is_shippable"`
	IsDigital        bool                   `json:"is_digital"`
	Attributes       []OrderItemAttributeKV `json:"attributes"`
}

func (oive *OrderedItemViewExternal) TableName() string {
	o := OrderedItemView{}
	return o.TableName()
}
