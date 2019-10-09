package repository

import "github.com/shopicano/shopicano-backend/models"

type AdminRepository interface {
	CreateShippingMethod(sm *models.ShippingMethod) error
	UpdateShippingMethod(sm *models.ShippingMethod) error
	ListShippingMethods(from, limit int) ([]models.ShippingMethod, error)
	ListActiveShippingMethods(from, limit int) ([]models.ShippingMethod, error)
	DeleteShippingMethod(ID string) error
	GetShippingMethod(ID string) (*models.ShippingMethod, error)

	CreatePaymentMethod(pm *models.PaymentMethod) error
	UpdatePaymentMethod(pm *models.PaymentMethod) error
	ListPaymentMethods(from, limit int) ([]models.PaymentMethod, error)
	ListActivePaymentMethods(from, limit int) ([]models.PaymentMethod, error)
	DeletePaymentMethod(ID string) error
	GetPaymentMethod(ID string) (*models.PaymentMethod, error)

	CreateAdditionalCharge(ac *models.AdditionalCharge) error
	UpdateAdditionalCharge(ac *models.AdditionalCharge) error
	ListAdditionalCharges(from, limit int) ([]models.AdditionalCharge, error)
	ListActiveAdditionalCharges(from, limit int) ([]models.AdditionalCharge, error)
	DeleteAdditionalCharge(ID string) error
	GetAdditionalCharge(ID string) (*models.AdditionalCharge, error)
}
