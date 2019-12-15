package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type CategoryRepository interface {
	Create(db *gorm.DB, c *models.Category) error
	List(db *gorm.DB, from, limit int) ([]models.ResCategorySearch, error)
	Search(db *gorm.DB, query string, from, limit int) ([]models.ResCategorySearch, error)
	ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.ResCategorySearchInternal, error)
	SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.ResCategorySearchInternal, error)
	Delete(db *gorm.DB, storeID, categoryID string) error
	Get(db *gorm.DB, storeID, categoryID string) (*models.Category, error)
	Update(db *gorm.DB, c *models.Category) error
}
