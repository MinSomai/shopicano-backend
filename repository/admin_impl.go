package repository

import (
	"github.com/shopicano/shopicano-backend/app"
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

func (au *AdminRepositoryImpl) CreateShippingMethod(sm *models.ShippingMethod) error {
	db := app.DB()
	if err := db.Table(sm.TableName()).Create(sm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) UpdateShippingMethod(sm *models.ShippingMethod) error {
	db := app.DB()
	if err := db.Table(sm.TableName()).
		Where("id = ?", sm.ID).Save(sm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) ListShippingMethods(from, limit int) ([]models.ShippingMethod, error) {
	db := app.DB()
	var data []models.ShippingMethod
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) ListActiveShippingMethods(from, limit int) ([]models.ShippingMethod, error) {
	db := app.DB()
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

func (au *AdminRepositoryImpl) DeleteShippingMethod(ID string) error {
	db := app.DB()
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) GetShippingMethod(ID string) (*models.ShippingMethod, error) {
	db := app.DB()
	m := models.ShippingMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}

func (au *AdminRepositoryImpl) CreatePaymentMethod(pm *models.PaymentMethod) error {
	db := app.DB()
	if err := db.Table(pm.TableName()).Create(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) UpdatePaymentMethod(pm *models.PaymentMethod) error {
	db := app.DB()
	if err := db.Table(pm.TableName()).
		Where("id = ?", pm.ID).Save(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) ListPaymentMethods(from, limit int) ([]models.PaymentMethod, error) {
	db := app.DB()
	var data []models.PaymentMethod
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) ListActivePaymentMethods(from, limit int) ([]models.PaymentMethod, error) {
	db := app.DB()
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

func (au *AdminRepositoryImpl) DeletePaymentMethod(ID string) error {
	db := app.DB()
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) GetPaymentMethod(ID string) (*models.PaymentMethod, error) {
	db := app.DB()
	m := models.PaymentMethod{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}

func (au *AdminRepositoryImpl) CreateAdditionalCharge(pm *models.AdditionalCharge) error {
	db := app.DB()
	if err := db.Table(pm.TableName()).Create(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) UpdateAdditionalCharge(pm *models.AdditionalCharge) error {
	db := app.DB()
	if err := db.Table(pm.TableName()).
		Where("id = ?", pm.ID).Save(pm).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) ListAdditionalCharges(from, limit int) ([]models.AdditionalCharge, error) {
	db := app.DB()
	var data []models.AdditionalCharge
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) ListActiveAdditionalCharges(from, limit int) ([]models.AdditionalCharge, error) {
	db := app.DB()
	var data []models.AdditionalCharge
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *AdminRepositoryImpl) DeleteAdditionalCharge(ID string) error {
	db := app.DB()
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

func (au *AdminRepositoryImpl) GetAdditionalCharge(ID string) (*models.AdditionalCharge, error) {
	db := app.DB()
	m := models.AdditionalCharge{}
	if err := db.Table(m.TableName()).
		Where("id = ?", ID).
		First(&m).Error; err != nil {
		return &m, err
	}
	return &m, nil
}
