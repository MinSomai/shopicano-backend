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

/**
ps.id AS id, ps.country_id AS country_id, l.name AS country_name, ps.business_name AS business_name,
ps.payout_minimum_threshold AS payout_minimum_threshold, ps.payout_method_details AS payout_method_details,
ps.vat_number AS vat_number, ps.payout_method_id AS payout_method_id, pom.name AS payout_method_name,
pom.inputs AS payout_method_inputs, ps.updated_at AS updated_at, ps.created_at AS created_at, ps.store_id AS store_id,
ps.business_address_id AS business_address_id, a.address AS business_address_address, a.city AS business_address_city,
a.state AS business_address_state, a.postcode AS business_address_post_code, ps.account_type_id AS account_type_id,
bat.name AS business_account_type_name
*/

type PayoutSettingsDetails struct {
	ID                      string    `json:"id"`
	CountryID               int64     `json:"country_id"`
	CountryName             string    `json:"country_name"`
	BusinessName            string    `json:"business_name"`
	PayoutMinimumThreshold  int64     `json:"payout_minimum_threshold"`
	PayoutMethodDetails     string    `json:"payout_method_details"`
	VatNumber               string    `json:"vat_number"`
	PayoutMethodID          string    `json:"payout_method_id"`
	PayoutMethodName        string    `json:"payout_method_name"`
	PayoutMethodInputs      string    `json:"payout_method_inputs"`
	UpdatedAt               time.Time `json:"updated_at"`
	CreatedAt               time.Time `json:"created_at"`
	StoreID                 string    `json:"store_id"`
	BusinessAddressID       string    `json:"business_address_id"`
	BusinessAddressAddress  string    `json:"business_address_address"`
	BusinessAddressCity     string    `json:"business_address_city"`
	BusinessAddressState    string    `json:"business_address_state"`
	BusinessAddressPostcode string    `json:"business_address_postcode"`
	BusinessAccountTypeID   string    `json:"business_account_type_id"`
	BusinessAccountTypeName string    `json:"business_account_type_name"`
}

func (psd *PayoutSettingsDetails) TableName() string {
	ps := PayoutSettings{}
	return ps.TableName()
}
