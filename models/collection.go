package models

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          string    `json:"id" gorm:"column:id;unique;not null"`
	Name        string    `json:"name" gorm:"column:name;primary_key"`
	StoreID     string    `json:"-" gorm:"column:store_id;primary_key"`
	Description string    `json:"description" gorm:"column:description;not null"`
	Image       string    `json:"image" gorm:"column:image;not null"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   time.Time `json:"-" gorm:"column:updated_at"`
}

func (c *Collection) TableName() string {
	return "collections"
}

func (c *Collection) ForeignKeys() []string {
	s := Store{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
	}
}
