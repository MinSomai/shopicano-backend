package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type AdditionalChargeRepository interface {
	Create(db *gorm.DB, ac *models.AdditionalCharge) error
	Update(db *gorm.DB, ac *models.AdditionalCharge) error
	List(db *gorm.DB, storeID string, from, limit int) ([]models.AdditionalCharge, error)
	Delete(db *gorm.DB, storeID, ID string) error
	Get(db *gorm.DB, storeID, ID string) (*models.AdditionalCharge, error)
}
