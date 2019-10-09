package models

import "time"

type Settings struct {
	ID                     string    `json:"id" sql:"id" json:"primary_key"`
	Name                   string    `json:"name" sql:"name" gorm:"not null"`
	URL                    string    `json:"url" sql:"url" gorm:"not null"`
	IsActive               bool      `json:"is_active" sql:"is_active" gorm:"not null"`
	CompanyName            string    `json:"company_name" sql:"company_name" gorm:"not null"`
	CompanyAddress         string    `json:"company_address" sql:"company_address" gorm:"not null"`
	CompanyCity            string    `json:"company_city" sql:"company_city" gorm:"not null"`
	CompanyCountry         string    `json:"company_country" sql:"company_country" gorm:"not null"`
	CompanyPostcode        string    `json:"company_postcode" sql:"company_postcode" gorm:"not null"`
	CompanyEmail           string    `json:"company_email" sql:"company_email" gorm:"not null"`
	CompanyPhone           string    `json:"company_phone" sql:"company_phone" gorm:"not null"`
	IsSignUpEnabled        bool      `json:"is_sign_up_enabled" sql:"is_sign_up_enabled" gorm:"not null"`
	IsStoreCreationEnabled bool      `json:"is_store_creation_enabled" sql:"is_store_creation_enabled" gorm:"not null"`
	TagLine                string    `json:"tag_line" sql:"tag_line" gorm:"not null"`
	CreatedAt              time.Time `json:"created_at" sql:"created_at" gorm:"not null"`
	UpdatedAt              time.Time `json:"updated_at" sql:"updated_at" gorm:"not null"`
}

func (s *Settings) TableName() string {
	return "settings"
}

func (s *Settings) ForeignKeys() []string {
	return []string{}
}
