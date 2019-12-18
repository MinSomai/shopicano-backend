package models

type OrderedItemViewExternal struct {
	OrderID             string `json:"order_id"`
	ProductID           string `json:"product_id"`
	Name                string `json:"name"`
	Quantity            int    `json:"quantity"`
	Price               int    `json:"price"`
	TotalVat            int    `json:"total_vat"`
	TotalTax            int    `json:"total_tax"`
	SubTotal            int    `json:"sub_total"`
	Description         string `json:"description"`
	SKU                 string `json:"sku"`
	AdditionalImages    string `json:"additional_images"`
	Image               string `json:"image"`
	IsShippable         bool   `json:"is_shippable"`
	IsDigital           bool   `json:"is_digital"`
	DigitalDownloadLink string `json:"digital_download_link"`
}

func (oive *OrderedItemViewExternal) TableName() string {
	o := OrderedItemView{}
	return o.TableName()
}
