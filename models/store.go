package models

import (
	"fmt"
	"time"
)

const (
	StoreRegistered StoreStatus = "registered"
	StoreActive     StoreStatus = "active"
	StoreSuspended  StoreStatus = "suspended"
	StoreBanned     StoreStatus = "banned"
)

type StoreStatus string

func (s *StoreStatus) IsValid() bool {
	for _, status := range []StoreStatus{StoreRegistered, StoreActive, StoreSuspended, StoreBanned} {
		if status == *s {
			return true
		}
	}
	return false
}

type Store struct {
	ID                       string      `json:"id" gorm:"column:id;primary_key"`
	Name                     string      `json:"name" gorm:"column:name;unique;not null"`
	Address                  string      `json:"address" gorm:"column:address;not null"`
	City                     string      `json:"city" gorm:"column:city;not null"`
	Country                  string      `json:"country" gorm:"column:country;not null"`
	Postcode                 string      `json:"postcode" gorm:"column:postcode;not null"`
	Email                    string      `json:"email" gorm:"column:email;unique;not null"`
	Phone                    string      `json:"phone" gorm:"column:phone;unique;not null"`
	Status                   StoreStatus `json:"status" gorm:"column:status;index"`
	CommissionRate           int64       `json:"commission_rate" gorm:"column:commission_rate;not null;default:0"`
	IsProductCreationEnabled bool        `json:"is_product_creation_enabled" gorm:"column:is_product_creation_enabled;not null;index"`
	IsOrderCreationEnabled   bool        `json:"is_order_creation_enabled" gorm:"column:is_order_creation_enabled;not null;index"`
	IsAutoConfirmEnabled     bool        `json:"is_auto_confirm_enabled" json:"column:is_auto_confirm_enabled;not null;index"`
	Description              string      `json:"description" gorm:"column:description;not null"`
	CreatedAt                time.Time   `json:"created_at" gorm:"column:created_at;index;not null"`
	UpdatedAt                time.Time   `json:"updated_at" gorm:"column:updated_at;not null"`
}

func (s *Store) TableName() string {
	return "stores"
}

func (s *Store) ForeignKeys() []string {
	return []string{}
}

type Staff struct {
	UserID       string `json:"user_id" gorm:"column:user_id;primary_key"`
	StoreID      string `json:"store_id" gorm:"column:store_id;primary_key"`
	PermissionID string `json:"permission_id" gorm:"column:permission_id;primary_key"`
	IsCreator    bool   `json:"is_creator" gorm:"column:is_creator;index"`
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
	ID         string     `json:"id" gorm:"column:id;primary_key"`
	Permission Permission `json:"permission" gorm:"column:permission;index;not null"`
}

func (up *StorePermission) TableName() string {
	return "store_permissions"
}
