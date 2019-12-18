package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type CollectionRepository interface {
	Create(db *gorm.DB, c *models.Collection) error
	List(db *gorm.DB, from, limit int) ([]models.Collection, error)
	Search(db *gorm.DB, query string, from, limit int) ([]models.Collection, error)
	ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.Collection, error)
	SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.Collection, error)
	Delete(db *gorm.DB, storeID, collectionID string) error
	Get(db *gorm.DB, storeID, collectionID string) (*models.Collection, error)
	Update(db *gorm.DB, c *models.Collection) error
}
