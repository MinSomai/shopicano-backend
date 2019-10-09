package models

import (
	"fmt"
	"time"
)

type Product struct {
	ID                  string    `json:"id" sql:"id" gorm:"unique"`
	Name                string    `json:"name" sql:"name" gorm:"primary_key"`
	Description         string    `json:"description" sql:"description"`
	IsPublished         bool      `json:"is_published" sql:"is_published" gorm:"index"`
	StoreID             string    `json:"store_id" sql:"store_id" gorm:"primary_key"`
	CategoryID          *string   `json:"category_id,omitempty" sql:"category_id" gorm:"index"`
	SKU                 string    `json:"sku" sql:"sku" gorm:"unique"`
	Stock               int       `json:"stock" sql:"stock" gorm:"index"`
	Unit                string    `json:"unit" sql:"unit"`
	Price               int       `json:"price" sql:"price" gorm:"index"`
	AdditionalImages    string    `json:"additional_images" sql:"additional_images"`
	Image               string    `json:"image,omitempty" sql:"image"`
	IsShippable         bool      `json:"is_shippable" sql:"is_shippable" gorm:"index"`
	IsDigital           bool      `json:"is_digital" sql:"is_digital" gorm:"index"`
	DigitalDownloadLink string    `json:"-" sql:"digital_download_link"`
	CreatedAt           time.Time `json:"created_at" sql:"created_at" gorm:"index"`
	UpdatedAt           time.Time `json:"updated_at" sql:"updated_at" gorm:"index"`
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) ForeignKeys() []string {
	s := Store{}
	c := Category{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
		fmt.Sprintf("category_id;%s(id);RESTRICT;RESTRICT", c.TableName()),
	}
}

type ProductOfCollection struct {
	CollectionID string `json:"collection_id" sql:"collection_id" gorm:"primary_key"`
	ProductID    string `json:"product_id" sql:"product_id" gorm:"primary_key"`
	StoreID      string `json:"store_id" sql:"store_id" gorm:"primary_key"`
}

func (cop *ProductOfCollection) TableName() string {
	return "product_of_collections"
}

func (cop *ProductOfCollection) ForeignKeys() []string {
	col := Collection{}
	p := Product{}
	s := Store{}

	return []string{
		fmt.Sprintf("product_id;%s(id);RESTRICT;RESTRICT", p.TableName()),
		fmt.Sprintf("collection_id;%s(id);RESTRICT;RESTRICT", col.TableName()),
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
	}
}
