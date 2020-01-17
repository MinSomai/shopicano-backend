package models

import (
	"fmt"
	"time"
)

type Address struct {
	ID        string    `json:"id" gorm:"column:id;primary_key"`
	UserID    string    `json:"-" gorm:"column:user_id;index;not null"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	Address   string    `json:"address" gorm:"column:address;not null"`
	City      string    `json:"city" gorm:"column:city;not null"`
	Country   string    `json:"country" gorm:"column:country;not null"`
	Postcode  string    `json:"postcode" gorm:"column:postcode;not null"`
	Email     string    `json:"email,omitempty" gorm:"column:email"`
	Phone     string    `json:"phone,omitempty" gorm:"column:phone"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (a *Address) TableName() string {
	return "addresses"
}

func (a *Address) ForeignKeys() []string {
	u := User{}

	return []string{
		fmt.Sprintf("user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
	}
}
