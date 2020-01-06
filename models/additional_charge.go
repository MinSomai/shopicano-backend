package models

import "time"

type AdditionalCharge struct {
	ID           string    `json:"id" gorm:"column:id;unique_index"`
	StoreID      string    `json:"store_id" gorm:"column:store_id;primary_key"`
	Name         string    `json:"name" gorm:"column:name;primary_key"`
	Amount       int       `json:"amount" gomr:"column:amount"`
	IsFlatAmount bool      `json:"is_flat_amount" gorm:"column:is_flat_amount"`
	AmountMax    int       `json:"amount_max" gorm:"column:amount_max"`
	AmountMin    int       `json:"amount_min" gorm:"column:amount_min"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (ac *AdditionalCharge) TableName() string {
	return "additional_charges"
}

func (ac *AdditionalCharge) CalculateAdditionalCharge(value int) int {
	if ac.Amount == 0 || value == 0 {
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
