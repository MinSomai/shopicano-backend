package models

import "time"

type ProductDetails struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	StoreID          string             `json:"store_id"`
	StoreName        string             `json:"store_name"`
	Description      string             `json:"description"`
	IsPublished      bool               `json:"is_published"`
	CategoryID       string             `json:"category_id,omitempty"`
	CategoryName     string             `json:"category_name,omitempty"`
	Image            string             `json:"image,omitempty"`
	IsShippable      bool               `json:"is_shippable"`
	IsDigital        bool               `json:"is_digital"`
	Price            int                `json:"price"`
	ProductCost      int                `json:"product_cost"`
	SKU              string             `json:"sku"`
	Stock            int                `json:"stock"`
	Unit             string             `json:"unit"`
	AdditionalImages string             `json:"additional_images"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Collections      []Collection       `json:"collections,omitempty"`
	Attributes       []ProductAttribute `json:"attributes,omitempty"`
}

type ProductDetailsInternal struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	StoreID             string             `json:"store_id"`
	StoreName           string             `json:"store_name"`
	Description         string             `json:"description"`
	IsPublished         bool               `json:"is_published"`
	CategoryID          string             `json:"category_id,omitempty"`
	CategoryName        string             `json:"category_name,omitempty"`
	Image               string             `json:"image,omitempty"`
	IsShippable         bool               `json:"is_shippable"`
	IsDigital           bool               `json:"is_digital"`
	Price               int                `json:"price"`
	SKU                 string             `json:"sku"`
	Stock               int                `json:"stock"`
	Unit                string             `json:"unit"`
	AdditionalImages    string             `json:"additional_images"`
	DigitalDownloadLink string             `json:"digital_download_link"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	Collections         []Collection       `json:"collections,omitempty"`
	Attributes          []ProductAttribute `json:"attributes,omitempty"`
}
