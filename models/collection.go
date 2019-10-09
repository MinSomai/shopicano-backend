package models

import (
	"fmt"
	"time"
)

type Collection struct {
	ID          string    `json:"id" sql:"id" gorm:"unique;not null"`
	Name        string    `json:"name" sql:"name" gorm:"primary_key"`
	StoreID     string    `json:"-" sql:"store_id" gorm:"primary_key"`
	Description string    `json:"description" sql:"description" gorm:"not null"`
	Image       string    `json:"image" sql:"image" gorm:"not null"`
	IsPublished bool      `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"-" sql:"updated_at"`
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
