package models

import (
	"fmt"
	"time"
)

type Category struct {
	ID          string    `json:"id" gorm:"column:id;unique;not null"`
	Name        string    `json:"name" gorm:"column:name;primary_key"`
	StoreID     string    `json:"-" gorm:"column:store_id;primary_key"`
	Description string    `json:"description" gorm:"column:description;not null"`
	Image       string    `json:"image" gorm:"column:image;not null"`
	IsPublished bool      `json:"is_published" gorm:"column:is_published;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) ForeignKeys() []string {
	s := Store{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
	}
}

type ResCategorySearch struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	StoreID     string `json:"store_id"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Count       int64  `json:"count"`
}

type ResCategorySearchInternal struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StoreID     string    `json:"store_id"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Count       int64     `json:"count"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
