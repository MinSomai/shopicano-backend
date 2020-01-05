package models

import "time"

type AdditionalCharge struct {
	ID           string    `json:"id" sql:"id" gorm:"primary_key"`
	Name         string    `json:"name" sql:"name" gorm:"unique"`
	Amount       int       `json:"amount" sql:"amount"`
	IsFlatAmount bool      `json:"is_flat_amount" json:"is_flat_amount"`
	AmountMax    int       `json:"amount_max" sql:"amount_max"`
	AmountMin    int       `json:"amount_min" sql:"amount_min"`
	IsPublished  bool      `json:"is_published" sql:"is_published" gorm:"index"`
	CreatedAt    time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" sql:"updated_at"`
}

func (ac *AdditionalCharge) TableName() string {
	return "additional_charges"
}

func (ac *AdditionalCharge) CalculateAdditionalCharge(value int) int {
	if ac.Amount == 0 {
		return 0
	}

	if ac.IsFlatAmount {
		charge := (value * ac.Amount) / 100
		if charge > ac.AmountMax && ac.AmountMax != 0 {
			return ac.AmountMax
		} else if charge < ac.AmountMin && ac.AmountMin != 0 {
			return ac.AmountMin
		}
		return charge
	}

	charge := ac.Amount
	if charge > ac.AmountMax && ac.AmountMax != 0 {
		return ac.AmountMax
	} else if charge < ac.AmountMin && ac.AmountMin != 0 {
		return ac.AmountMin
	}
	return charge
}
