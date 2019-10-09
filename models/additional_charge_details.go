package models

type AdditionalChargeDetails struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	ChargeType AdditionalChargeType       `json:"charge_type"`
	Amount     int                        `json:"amount"`
	AmountType AdditionalChargeAmountType `json:"amount_type"`
	AmountMax  int                        `json:"amount_max"`
	AmountMin  int                        `json:"amount_min"`
}

func (acd *AdditionalChargeDetails) CalculateAdditionalCharge(value int) int {
	if acd.Amount == 0 {
		return 0
	}

	switch acd.AmountType {
	case Percent:
		charge := (value * acd.Amount) / 100
		if charge > acd.AmountMax && acd.AmountMax != 0 {
			return acd.AmountMax
		} else if charge < acd.AmountMin && acd.AmountMin != 0 {
			return acd.AmountMin
		}
		return charge
	case Fixed:
		charge := acd.Amount
		if charge > acd.AmountMax && acd.AmountMax != 0 {
			return acd.AmountMax
		} else if charge < acd.AmountMin && acd.AmountMin != 0 {
			return acd.AmountMin
		}
		return charge
	}
	return 0
}
