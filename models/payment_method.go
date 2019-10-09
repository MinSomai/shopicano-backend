package models

import (
	"github.com/shopicano/shopicano-backend/log"
	"time"
)

type PaymentMethod struct {
	ID               string    `json:"id" sql:"id" gorm:"primary_key"`
	Name             string    `json:"name" sql:"name" gorm:"unique;not null"`
	ProcessingFee    int       `json:"processing_fee" gorm:"processing_fee"`
	MinProcessingFee int       `json:"min_processing_fee" gorm:"min_processing_fee"`
	MaxProcessingFee int       `json:"max_processing_fee" sql:"max_processing_fee"`
	IsPublished      bool      `json:"is_published" sql:"is_published" gorm:"index"`
	IsOfflinePayment bool      `json:"is_offline_payment" sql:"is_offline_payment" gorm:"is_offline_payment"`
	CreatedAt        time.Time `json:"created_at" sql:"created_at" gorm:"not null;index"`
	UpdatedAt        time.Time `json:"updated_at" sql:"updated_at" gorm:"not null"`
}

func (pm *PaymentMethod) TableName() string {
	return "payment_methods"
}

func (pm *PaymentMethod) CalculateProcessingFee(bill int) int {
	if pm.IsOfflinePayment {
		return 0
	}

	log.Log().Info("Bill : ", bill)
	log.Log().Info("Pf : ", pm.ProcessingFee)

	fee := ((bill * pm.ProcessingFee) / 100) / 100

	log.Log().Info("Fee : ", fee)

	if fee > pm.MaxProcessingFee && pm.MaxProcessingFee != 0 {
		return pm.MaxProcessingFee
	} else if fee < pm.MinProcessingFee && pm.MinProcessingFee != 0 {
		return pm.MinProcessingFee
	}
	return fee
}
