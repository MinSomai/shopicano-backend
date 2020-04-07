package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type MarketplaceRepository interface {
	CreateShippingMethod(db *gorm.DB, sm *models.ShippingMethod) error
	UpdateShippingMethod(db *gorm.DB, sm *models.ShippingMethod) error
	ListShippingMethods(db *gorm.DB, from, limit int) ([]models.ShippingMethod, error)
	ListActiveShippingMethods(db *gorm.DB, from, limit int) ([]models.ShippingMethod, error)
	ListShippingMethodsByLocation(db *gorm.DB, locationID int64) ([]models.ShippingMethod, error)
	ListShippingMethodsByLocationForUser(db *gorm.DB, locationID int64) ([]models.ShippingMethod, error)
	DeleteShippingMethod(db *gorm.DB, ID string) error
	GetShippingMethod(db *gorm.DB, ID string) (*models.ShippingMethod, error)
	GetShippingMethodForUser(db *gorm.DB, ID string) (*models.ShippingMethod, error)

	CreatePaymentMethod(db *gorm.DB, pm *models.PaymentMethod) error
	UpdatePaymentMethod(db *gorm.DB, pm *models.PaymentMethod) error
	ListPaymentMethods(db *gorm.DB, from, limit int) ([]models.PaymentMethod, error)
	ListActivePaymentMethods(db *gorm.DB, from, limit int) ([]models.PaymentMethod, error)
	ListPaymentMethodsByLocation(db *gorm.DB, locationID int64) ([]models.PaymentMethod, error)
	ListPaymentMethodsByLocationForUser(db *gorm.DB, locationID int64) ([]models.PaymentMethod, error)
	DeletePaymentMethod(db *gorm.DB, ID string) error
	GetPaymentMethod(db *gorm.DB, ID string) (*models.PaymentMethod, error)
	GetPaymentMethodForUser(db *gorm.DB, ID string) (*models.PaymentMethod, error)

	GetSettings(db *gorm.DB) (*models.Settings, error)
	GetSettingsDetails(db *gorm.DB) (*models.SettingsDetails, error)
	UpdateSettings(db *gorm.DB, s *models.Settings) error
}
