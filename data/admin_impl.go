package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type AdminRepositoryImpl struct {
}

var adminRepository AdminRepository

func NewAdminRepository() AdminRepository {
	if adminRepository == nil {
		adminRepository = &AdminRepositoryImpl{}
	}
	return adminRepository
}

func (au *AdminRepositoryImpl) CreateShippingMethod(db *gorm.DB, sm *models.ShippingMethod) error {
	if err := db.Table(sm.TableName()).Create(sm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) UpdateShippingMethod(db *gorm.DB, sm *models.ShippingMethod) error {
	if err := db.Table(sm.TableName()).
		Where("id = ?", sm.ID).Save(sm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) ListShippingMethods(db *gorm.DB, from, limit int) ([]models.ShippingMethod, error) {
	var data []models.ShippingMethod
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) ListActiveShippingMethods(db *gorm.DB, from, limit int) ([]models.ShippingMethod, error) {
	var data []models.ShippingMethod
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) DeleteShippingMethod(db *gorm.DB, ID string) error {
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) GetShippingMethod(db *gorm.DB, ID string) (*models.ShippingMethod, error) {
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}

func (au *AdminRepositoryImpl) CreatePaymentMethod(db *gorm.DB, pm *models.PaymentMethod) error {
	if err := db.Table(pm.TableName()).Create(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) UpdatePaymentMethod(db *gorm.DB, pm *models.PaymentMethod) error {
	if err := db.Table(pm.TableName()).
		Where("id = ?", pm.ID).Save(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) ListPaymentMethods(db *gorm.DB, from, limit int) ([]models.PaymentMethod, error) {
	var data []models.PaymentMethod
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) ListActivePaymentMethods(db *gorm.DB, from, limit int) ([]models.PaymentMethod, error) {
	var data []models.PaymentMethod
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) DeletePaymentMethod(db *gorm.DB, ID string) error {
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) GetPaymentMethod(db *gorm.DB, ID string) (*models.PaymentMethod, error) {
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}

func (au *AdminRepositoryImpl) GetSettings(db *gorm.DB) (*models.SettingsDetails, error) {
	settings := models.Settings{}
	settingsDetails := models.SettingsDetails{}
	a := models.Address{}
	if err := db.Table(fmt.Sprintf("%s AS s", settings.TableName())).
		Where("s.id = ?", "1").
		Select("s.id AS id, s.name AS name, s.website AS website, s.is_active AS is_active, a.address AS address, a.city AS city, a.country AS country, a.postcode AS postcode, a.email AS email, a.phone AS phone, s.is_sign_up_enabled AS is_sign_up_enabled, s.is_store_creation_enabled AS is_store_creation_enabled, s.default_commission_rate AS default_commission_rate, s.tag_line AS tag_line, s.created_at AS created_at, s.updated_at AS updated_at").
		Joins(fmt.Sprintf("LEFT JOIN %s AS a ON s.company_address_id = a.id", a.TableName())).
		Find(&settingsDetails).Error; err != nil {
		return nil, err
	}
	return &settingsDetails, nil
}
