package data

import "github.com/shopicano/shopicano-backend/models"

type AddressRepository interface {
	CreateAddress(a *models.Address) error
	UpdateAddress(a *models.Address) error
	ListAddresses(userID, string, from, limit int) ([]models.Address, error)
	DeleteAddress(userID, addressID string) error
}
