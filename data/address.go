package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type AddressRepository interface {
	CreateAddress(db *gorm.DB, a *models.Address) error
	UpdateAddress(db *gorm.DB, a *models.Address) error
	GetAddress(db *gorm.DB, userID, addressID string) (*models.Address, error)
	ListAddresses(db *gorm.DB, userID string, from, limit int) ([]models.Address, error)
	DeleteAddress(db *gorm.DB, userID, addressID string) error
}
