package models

import "time"

const (
	Vat AdditionalChargeType = "vat"
	Tax AdditionalChargeType = "tax"

	Percent AdditionalChargeAmountType = "percent"
	Fixed   AdditionalChargeAmountType = "fixed"
)

type AdditionalChargeType string
type AdditionalChargeAmountType string

type AdditionalCharge struct {
	ID          string                     `json:"id" sql:"id" gorm:"primary_key"`
	Name        string                     `json:"name" sql:"name" gorm:"unique"`
	ChargeType  AdditionalChargeType       `json:"charge_type" sql:"charge_type" gorm:"index;not null"`
	Amount      int                        `json:"amount" sql:"amount"`
	AmountType  AdditionalChargeAmountType `json:"amount_type" sql:"amount_type" gorm:"index;not null"`
	AmountMax   int                        `json:"amount_max" sql:"amount_max"`
	AmountMin   int                        `json:"amount_min" sql:"amount_min"`
	IsPublished bool                       `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt   time.Time                  `json:"created_at" sql:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at" sql:"updated_at"`
}

func (ac *AdditionalCharge) TableName() string {
	return "additional_charges"
}

func (ac *AdditionalCharge) CalculateAdditionalCharge(value int) int {
	if ac.Amount == 0 {
		return 0
	}

	switch ac.AmountType {
	case Percent:
		charge := (value * ac.Amount) / 100
		if charge > ac.AmountMax && ac.AmountMax != 0 {
			return ac.AmountMax
		} else if charge < ac.AmountMin && ac.AmountMin != 0 {
			return ac.AmountMin
		}
		return charge
	case Fixed:
		charge := ac.Amount
		if charge > ac.AmountMax && ac.AmountMax != 0 {
			return ac.AmountMax
		} else if charge < ac.AmountMin && ac.AmountMin != 0 {
			return ac.AmountMin
		}
		return charge
	}
	return 0
}
