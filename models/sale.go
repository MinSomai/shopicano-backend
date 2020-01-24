package models

import "time"

type Sale struct {
	ID             string    `json:"id" gorm:"column:id;primary_key"`
	CollectionID   *string   `json:"collection_id" gorm:"column:collection_id;index"`
	CategoryID     *string   `json:"category_id" gorm:"column:category_id;index"`
	IsActive       bool      `json:"is_active" gorm:"column:is_active;index"`
	DiscountAmount int       `json:"discount_amount" gorm:"column:discount_amount"`
	IsFlatDiscount bool      `json:"is_flat_discount" gorm:"column:is_flat_discount"`
	MaxDiscount    int       `json:"max_discount" gorm:"column:max_discount"`
	Message        string    `json:"message" gorm:"column:message"`
	StartAt        time.Time `json:"start_at" gorm:"column:start_at;index"`
	EndAt          time.Time `json:"end_at" gorm:"column:end_at;index"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;index"`
}

func (s *Sale) TableName() string {
	return "sales"
}
