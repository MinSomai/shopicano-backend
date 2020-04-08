package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type AddressView struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Name      string    `json:"name,omitempty"`
	Address   string    `json:"address,omitempty"`
	City      string    `json:"city,omitempty"`
	Country   string    `json:"country,omitempty"`
	State     string    `json:"state,omitempty"`
	Postcode  string    `json:"postcode,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (av *AddressView) TableName() string {
	return "addresses_view"
}

func (av *AddressView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT a.id AS id, a.user_id AS user_id, a.name AS name, a.address AS address, a.city AS city,"+
		" loc.name AS country, a.state AS state, a.postcode AS postcode, a.email AS email, a.phone AS phone, a.created_at AS created_at,"+
		" a.updated_at AS updated_at"+
		" FROM addresses AS a"+
		" LEFT JOIN locations AS loc ON a.country_id = loc.id", av.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (av *AddressView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", av.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
