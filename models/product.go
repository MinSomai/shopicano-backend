package models

import (
	"fmt"
	"time"
)

type Product struct {
	ID                  string    `json:"id" gorm:"column:id;unique"`
	Name                string    `json:"name" gorm:"column:name;primary_key"`
	Description         string    `json:"description" gorm:"column:description"`
	IsPublished         bool      `json:"is_published" gorm:"column:is_published;index"`
	StoreID             string    `json:"store_id" gorm:"column:store_id;primary_key"`
	CategoryID          *string   `json:"category_id,omitempty" gorm:"column:category_id;index"`
	SKU                 string    `json:"sku" gorm:"column:sku;unique"`
	Stock               int       `json:"stock" gorm:"column:stock;index"`
	MaxQuantityCount    int       `json:"max_quantity_count" gorm:"column:max_quantity_count;not null;default:10"`
	Unit                string    `json:"unit" gorm:"column:unit"`
	Price               int       `json:"price" gorm:"column:price;index"`
	ProductCost         int       `json:"product_cost" gorm:"column:product_cost;index"`
	AdditionalImages    string    `json:"additional_images" gorm:"column:additional_images"`
	Image               string    `json:"image,omitempty" gorm:"column:image"`
	IsShippable         bool      `json:"is_shippable" gorm:"column:is_shippable;index"`
	IsDigital           bool      `json:"is_digital" gorm:"column:is_digital;index"`
	DigitalDownloadLink string    `json:"-" gorm:"column:digital_download_link"`
	DownloadCounter     int       `json:"download_counter" gorm:"column:download_counter;default:0"`
	CreatedAt           time.Time `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"column:updated_at;index"`
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
