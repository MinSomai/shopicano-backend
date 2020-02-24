package models

import "time"

type SettingsDetails struct {
	ID                           string         `json:"id"`
	Name                         string         `json:"name"`
	Website                      string         `json:"website"`
	Status                       PlatformStatus `json:"status"`
	Address                      string         `json:"address"`
	City                         string         `json:"city"`
	Country                      string         `json:"country"`
	Postcode                     string         `json:"postcode"`
	Email                        string         `json:"email"`
	Phone                        string         `json:"phone"`
	IsSignUpEnabled              bool           `json:"is_sign_up_enabled"`
	IsStoreCreationEnabled       bool           `json:"is_store_creation_enabled"`
	EnabledAutoStoreConfirmation bool           `json:"enabled_auto_store_confirmation"`
	DefaultCommissionRate        int64          `json:"default_commission_rate"`
	TagLine                      string         `json:"tag_line"`
	CreatedAt                    time.Time      `json:"created_at"`
	UpdatedAt                    time.Time      `json:"updated_at"`
}
