package data

import (
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/models"
)

type AddressRepositoryImpl struct {
}

var addressRepository *AddressRepositoryImpl

func NewAddressRepository() *AddressRepositoryImpl {
	if addressRepository == nil {
		addressRepository = &AddressRepositoryImpl{}
	}
	return addressRepository
}

func (au *AddressRepositoryImpl) CreateAddress(a *models.Address) error {
	db := app.DB()
	if err := db.Table(a.TableName()).Create(a).Error; err != nil {
		return err
	}
	return nil
}

func (au *AddressRepositoryImpl) UpdateAddress(a *models.Address) error {
	db := app.DB()
	if err := db.Table(a.TableName()).
		Where("id = ? AND user_id = ?", a.ID, a.UserID).Updates(a).Error; err != nil {
		return err
	}
	return nil
}

func (au *AddressRepositoryImpl) ListAddresses(userId string, from, limit int) ([]models.Address, error) {
	db := app.DB()
	var addresses []models.Address
	address := models.Address{}
	if err := db.Table(address.TableName()).
		Offset(from).Limit(limit).
		Where("user_id = ?", userId).
		Order("created_at DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (au *AddressRepositoryImpl) DeleteAddress(userID, addressID string) error {
	db := app.DB()
	address := models.Address{}
	if err := db.Table(address.TableName()).
		Where("user_id = ? AND id = ?", userID, addressID).
		Delete(&address).Error; err != nil {
		return err
	}
	return nil
}
