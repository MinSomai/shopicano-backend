package repository

import (
	"github.com/shopicano/shopicano-backend/models"
)

type StoreRepository interface {
	GetStoreUserProfile(userID string) (*models.StoreUserProfile, error)
	CreateStore(s *models.Store, userID string) error
	FindStoreByID(ID string) (*models.Store, error)
}
