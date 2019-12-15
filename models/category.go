package models

import (
	"fmt"
	"time"
)

type Category struct {
	ID          string    `json:"id" sql:"id" gorm:"unique;not null"`
	Name        string    `json:"name" sql:"name" gorm:"primary_key"`
	StoreID     string    `json:"-" sql:"store_id" gorm:"primary_key"`
	Description string    `json:"description" sql:"description" gorm:"not null"`
	Image       string    `json:"image" sql:"image" gorm:"not null"`
	IsPublished bool      `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" sql:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at" sql:"updated_at"`
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
