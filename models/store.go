package models

import (
	"fmt"
	"time"
)

const (
	StoreRegistered StoreStatus = "registered"
	StoreOpen       StoreStatus = "open"
	StoreActive     StoreStatus = "active"
	StoreClosed     StoreStatus = "closed"
	StoreSuspended  StoreStatus = "suspended"
	StoreBanned     StoreStatus = "banned"
)

type StoreStatus string

type Store struct {
	ID                       string      `json:"id" sql:"id" gorm:"primary_key"`
	Name                     string      `json:"name" sql:"name" gorm:"unique;not null"`
	Address                  string      `json:"address" sql:"address" gorm:"not null"`
	City                     string      `json:"city" sql:"city" gorm:"not null"`
	Country                  string      `json:"country" sql:"country" gorm:"not null"`
	Postcode                 string      `json:"postcode" sql:"postcode" gorm:"not null"`
	Email                    string      `json:"email" sql:"email" gorm:"unique;not null"`
	Phone                    string      `json:"phone" sql:"phone" gorm:"unique;not null"`
	Status                   StoreStatus `json:"status" sql:"status" gorm:"index"`
	IsProductCreationEnabled bool        `json:"-" sql:"is_product_creation_enabled" gorm:"not null"`
	IsOrderCreationEnabled   bool        `json:"is_order_creation_enabled" sql:"is_order_creation_enabled" gorm:"not null"`
	Description              string      `json:"description" sql:"description" gorm:"not null"`
	CreatedAt                time.Time   `json:"created_at" sql:"created_at" gorm:"index;not null"`
	UpdatedAt                time.Time   `json:"updated_at" sql:"updated_at" gorm:"not null"`
}

func (s *Store) TableName() string {
	return "stores"
}

func (s *Store) ForeignKeys() []string {
	return []string{}
}

type Staff struct {
	UserID       string `json:"user_id" sql:"user_id" gorm:"primary_key"`
	StoreID      string `json:"store_id" sql:"store_id" gorm:"primary_key"`
	PermissionID string `json:"permission_id" sql:"permission_id" gorm:"primary_key"`
}

func (sf *Staff) TableName() string {
	return "staffs"
}

func (sf *Staff) ForeignKeys() []string {
	u := User{}
	s := Store{}
	sp := StorePermission{}

	return []string{
		fmt.Sprintf("user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
		fmt.Sprintf("permission_id;%s(id);RESTRICT;RESTRICT", sp.TableName()),
	}
}

type StorePermission struct {
	ID         string     `json:"id" sql:"id" gorm:"primary_key"`
	Permission Permission `json:"permission" sql:"permission" gorm:"index;not null"`
}

func (up *StorePermission) TableName() string {
	return "store_permissions"
}
