package models

type DiscountFor struct {
	DiscountID string `json:"discount_id" gorm:"column:discount_id;primary_key"`
	UserID     string `json:"user_id" gorm:"column:user_id;primary_key"`
}

func (df *DiscountFor) TableName() string {
	return "discount_for"
}
