package models

import (
	"fmt"
	"time"
)

type Review struct {
	ID          string    `json:"id" gorm:"column:id;primary_key"`
	OrderID     string    `json:"-" gorm:"column:order_id;index;unique"`
	Rating      int       `json:"rating" gorm:"column:rating;index"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
}

func (r *Review) TableName() string {
	return "reviews"
}

func (r *Review) ForeignKeys() []string {
	o := Order{}

	return []string{
		fmt.Sprintf("order_id;%s(id);RESTRICT;RESTRICT", o.TableName()),
	}
}
