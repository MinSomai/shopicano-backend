package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type AdditionalChargeRepositoryImpl struct {
}

var additionalChargeRepository AdditionalChargeRepository

func NewAdditionalChargeRepository() AdditionalChargeRepository {
	if additionalChargeRepository == nil {
		additionalChargeRepository = &AdditionalChargeRepositoryImpl{}
	}
	return additionalChargeRepository
}

func (au *AdditionalChargeRepositoryImpl) Create(db *gorm.DB, ac *models.AdditionalCharge) error {
	if err := db.Table(ac.TableName()).Create(ac).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdditionalChargeRepositoryImpl) Update(db *gorm.DB, ac *models.AdditionalCharge) error {
	if err := db.Table(ac.TableName()).
		Select("name, amount, is_flat_amount, amount_max, amount_min, updated_at").
		Where("id = ? AND store_id = ?", ac.ID, ac.StoreID).
		Updates(map[string]interface{}{
			"name":           ac.Name,
			"amount":         ac.Amount,
			"is_flat_amount": ac.IsFlatAmount,
			"amount_max":     ac.AmountMax,
			"amount_min":     ac.AmountMin,
			"updated_at":     ac.UpdatedAt,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdditionalChargeRepositoryImpl) List(db *gorm.DB, storeID string, from, limit int) ([]models.AdditionalCharge, error) {
	var data []models.AdditionalCharge
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("updated_at DESC").Find(&data, "store_id = ?", storeID).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdditionalChargeRepositoryImpl) Delete(db *gorm.DB, storeID, ID string) error {
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Where("id = ? AND store_id = ?", ID, storeID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdditionalChargeRepositoryImpl) Get(db *gorm.DB, storeID, ID string) (*models.AdditionalCharge, error) {
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Where("id = ? AND store_id = ?", ID, storeID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}
