package models

import (
	"fmt"
	"time"
)

type Settings struct {
	ID                     string    `json:"id" gorm:"column:id;primary_key"`
	Name                   string    `json:"name" gorm:"column:name;not null"`
	Website                string    `json:"website" gorm:"column:website;not null"`
	IsActive               bool      `json:"is_active" gorm:"column:is_active;not null"`
	CompanyAddressID       string    `json:"company_address_id" gorm:"column:company_address_id;not null"`
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
	a := Address{}
	return []string{fmt.Sprintf("company_address_id;%s(id);RESTRICT;RESTRICT", a.TableName())}
}
