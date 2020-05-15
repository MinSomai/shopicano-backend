package models

import "time"

type BusinessAccountType struct {
	ID          string    `json:"id" gorm:"column:id;primary_key"`
	Name        string    `json:"name" gorm:"column:name;unique_index;not null"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published;default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (bat *BusinessAccountType) TableName() string {
	return "business_account_types"
}
