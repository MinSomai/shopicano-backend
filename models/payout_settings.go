package models

import (
	"fmt"
	"time"
)

type PayoutSettings struct {
	ID                     string    `json:"id" gorm:"column:id;primary_key"`
	StoreID                string    `json:"store_id" gorm:"column:store_id;unique_index;not null"`
	CountryID              int64     `json:"country_id" gorm:"column:country_id;not null"`
	AccountTypeID          string    `json:"account_type_id" gorm:"column:account_type_id;not null"`
	BusinessName           string    `json:"business_name" gorm:"column:business_name;not null"`
	BusinessAddressID      string    `json:"business_address_id" gorm:"column:business_address_id;not null"`
	VatNumber              string    `json:"vat_number" gorm:"column:vat_number"`
	PayoutMethodID         string    `json:"payout_method_id" gorm:"column:payout_method_id;not null"`
	PayoutMethodDetails    string    `json:"payout_method_details" gorm:"column:payout_method_details;not null"`
	PayoutMinimumThreshold int64     `json:"payout_minimum_threshold" gorm:"column:payout_minimum_threshold;not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"column:created_at;index;not null"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (pos *PayoutSettings) TableName() string {
	return "payout_settings"
}

func (pos *PayoutSettings) ForeignKeys() []string {
	s := Store{}
	l := Location{}
	bat := BusinessAccountType{}
	pom := PayoutMethod{}
	a := Address{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
		fmt.Sprintf("country_id;%s(id);RESTRICT;RESTRICT", l.TableName()),
		fmt.Sprintf("account_type_id;%s(id);RESTRICT;RESTRICT", bat.TableName()),
		fmt.Sprintf("payout_method_id;%s(id);RESTRICT;RESTRICT", pom.TableName()),
		fmt.Sprintf("business_address_id;%s(id);RESTRICT;RESTRICT", a.TableName()),
	}
}
