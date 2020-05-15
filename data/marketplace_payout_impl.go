package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

func (au *MarketplaceRepositoryImpl) CreatePayoutMethod(db *gorm.DB, m *models.PayoutMethod) error {
	if err := db.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (au *MarketplaceRepositoryImpl) UpdatePayoutMethod(db *gorm.DB, m *models.PayoutMethod) error {
	if err := db.Table(m.TableName()).
		Select("name, inputs, is_published, updated_at").
		Updates(map[string]interface{}{
			"name":         m.Name,
			"inputs":       m.Inputs,
			"is_published": m.IsPublished,
			"updated_at":   m.UpdatedAt,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (au *MarketplaceRepositoryImpl) ListPayoutMethods(db *gorm.DB, from, limit int) ([]models.PayoutMethod, error) {
	pom := models.PayoutMethod{}
	var data []models.PayoutMethod
	if err := db.Table(pom.TableName()).
		Order("created_at DESC").
		Offset(from).
		Limit(limit).
		Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *MarketplaceRepositoryImpl) ListPayoutMethodForUser(db *gorm.DB, from, limit int) ([]models.PayoutMethod, error) {
	pom := models.PayoutMethod{}
	var data []models.PayoutMethod
	if err := db.Table(pom.TableName()).
		Order("created_at DESC").
		Offset(from).
		Limit(limit).
		Find(&data, "is_published = ?", true).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (au *MarketplaceRepositoryImpl) DeletePayoutMethod(db *gorm.DB, ID string) error {
	pom := models.PayoutMethod{}
	if err := db.Table(pom.TableName()).
		Delete(&pom, "id = ?", ID).Error; err != nil {
		return err
	}
	return nil
}

func (au *MarketplaceRepositoryImpl) GetPayoutMethod(db *gorm.DB, ID string) (*models.PayoutMethod, error) {
	pom := models.PayoutMethod{}
	if err := db.Table(pom.TableName()).
		Find(&pom, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &pom, nil
}

func (au *MarketplaceRepositoryImpl) GetPayoutMethodForUser(db *gorm.DB, ID string) (*models.PayoutMethod, error) {
	pom := models.PayoutMethod{}
	if err := db.Table(pom.TableName()).
		Find(&pom, "id = ? AND is_published = ?", ID, true).Error; err != nil {
		return nil, err
	}
	return &pom, nil
}
