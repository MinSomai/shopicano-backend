package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type StoreRepository interface {
	GetStoreUserProfile(db *gorm.DB, userID string) (*models.StoreUserProfile, error)
	CreateStore(db *gorm.DB, s *models.Store) error
	FindStoreByID(db *gorm.DB, ID string) (*models.Store, error)
	AddStoreStuff(db *gorm.DB, staff *models.Staff) error
	UpdateStoreStuffPermission(db *gorm.DB, storeID, userID, permissionID string) error
	DeleteStoreStuffPermission(db *gorm.DB, storeID, userID string) error
}
