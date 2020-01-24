package models

type SaleUsage struct {
	SaleID    string `json:"sale_id" gorm:"column:sale_id"`
	OrderID   string `json:"order_id" gorm:"column:order_id"`
	ProductID string `json:"product_id" gorm:"column:product_id"`
}
