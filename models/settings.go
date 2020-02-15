package models

import "time"

type Settings struct {
	ID                     string    `json:"id" gorm:"column:id;primary_key"`
	Name                   string    `json:"name" gorm:"column:name;not null"`
	URL                    string    `json:"url" gorm:"column:url;not null"`
	IsActive               bool      `json:"is_active" gorm:"column:is_active;not null"`
	CompanyName            string    `json:"company_name" gorm:"column:company_name;not null"`
	CompanyAddress         string    `json:"company_address" gorm:"column:company_address;not null"`
	CompanyCity            string    `json:"company_city" gorm:"column:company_city;not null"`
	CompanyCountry         string    `json:"company_country" gorm:"column:company_country;not null"`
	CompanyPostcode        string    `json:"company_postcode" gorm:"column:company_postcode;not null"`
	CompanyEmail           string    `json:"company_email" gorm:"column:company_email;not null"`
	CompanyPhone           string    `json:"company_phone" gorm:"column:company_phone;not null"`
	IsSignUpEnabled        bool      `json:"is_sign_up_enabled" gorm:"column:is_sign_up_enabled;not null"`
	IsStoreCreationEnabled bool      `json:"is_store_creation_enabled" gorm:"column:is_store_creation_enabled;not null"`
	DefaultCommissionRate  int64     `json:"default_commission_rate" gorm:"column:default_commission_rate;not null;default:0"`
	TagLine                string    `json:"tag_line" gorm:"column:tag_line;not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (s *Settings) TableName() string {
	return "settings"
}

func (s *Settings) ForeignKeys() []string {
	return []string{}
}
