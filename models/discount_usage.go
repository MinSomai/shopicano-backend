package models

type DiscountUsage struct {
	ID         string `json:"id" gorm:"column:id;primary_key"`
	DiscountID string `json:"discount_id" gorm:"column:discount_id;index"`
	UserID     string `json:"user_id" gorm:"column:user_id;index"`
}

func (du *DiscountUsage) TableName() string {
	return "discount_usages"
}
