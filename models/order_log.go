package models

import (
	"fmt"
	"time"
)

type OrderLog struct {
	ID        string    `json:"id" gorm:"column:id;primary_key"`
	OrderID   string    `json:"order_id" gorm:"column:order_id;index"`
	Action    string    `json:"action" gorm:"column:action"`
	Details   string    `json:"details" gorm:"column:details"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (ol *OrderLog) TableName() string {
	return "order_logs"
}

func (ol *OrderLog) ForeignKeys() []string {
	o := Order{}

	return []string{
		fmt.Sprintf("order_id;%s(id);RESTRICT;RESTRICT", o.TableName()),
	}
}
