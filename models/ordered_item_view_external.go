package models

type OrderedItemViewExternal struct {
	OrderID             string `json:"order_id"`
	ProductID           string `json:"product_id"`
	Name                string `json:"name"`
	Quantity            int64  `json:"quantity"`
	Price               int64  `json:"price"`
	TotalVat            int64  `json:"total_vat"`
	TotalTax            int64  `json:"total_tax"`
	SubTotal            int64  `json:"sub_total"`
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
