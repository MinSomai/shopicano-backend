package models

import "time"

type PayoutMethod struct {
	ID          string    `json:"id" gorm:"column:id;primary_key"`
	Name        string    `json:"name" gorm:"column:name;unique;not null"`
	Inputs      string    `json:"inputs" gorm:"column:inputs;not null"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (pom *PayoutMethod) TableName() string {
	return "payout_methods"
}
