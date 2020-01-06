package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/helpers"
	"github.com/shopicano/shopicano-backend/models"
)

type ProductRepository interface {
	Create(db *gorm.DB, p *models.Product) error
	Update(db *gorm.DB, p *models.Product) error
	List(db *gorm.DB, from, limit int) ([]models.ProductDetails, error)
	Search(db *gorm.DB, query string, from, limit int) ([]models.ProductDetails, error)
	ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.ProductDetailsInternal, error)
	SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.ProductDetailsInternal, error)
	Delete(db *gorm.DB, storeID, productID string) error
	Get(db *gorm.DB, productID string) (*models.Product, error)
	GetAsStoreStuff(db *gorm.DB, storeID, productID string) (*models.Product, error)
	GetDetails(db *gorm.DB, productID string) (*models.ProductDetails, error)
	GetDetailsAsStoreStuff(db *gorm.DB, storeID, productID string) (*models.ProductDetailsInternal, error)
	GetForOrder(db *gorm.DB, storeID, productID string, quantity int) (*models.Product, error)
	StatsAsStoreStuff(db *gorm.DB, storeID string, offset, limit int) ([]helpers.ProductStats, error)
	Stats(db *gorm.DB, offset, limit int) ([]helpers.ProductStats, error)
	AddAttribute(db *gorm.DB, v *models.ProductAttribute) error
	RemoveAttribute(db *gorm.DB, productID, attributeKey string) error
}
