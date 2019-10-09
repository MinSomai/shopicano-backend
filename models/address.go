package models

import (
	"fmt"
	"time"
)

type Address struct {
	ID        string    `json:"id" sql:"id" gorm:"primary_key"`
	UserID    string    `json:"-" sql:"user_id" gorm:"index"`
	Name      string    `json:"name" sql:"name"`
	House     string    `json:"house" sql:"house"`
	Road      string    `json:"road" sql:"road"`
	City      string    `json:"city" sql:"city"`
	Country   string    `json:"country" sql:"country"`
	Postcode  string    `json:"postcode" sql:"postcode"`
	Email     string    `json:"email,omitempty" sql:"email"`
	Phone     string    `json:"phone,omitempty" sql:"phone"`
	CreatedAt time.Time `json:"created_at" sql:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at" sql:"updated_at"`
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
