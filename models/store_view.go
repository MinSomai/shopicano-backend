package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type StoreView struct {
	ID                       string      `json:"id"`
	Name                     string      `json:"name"`
	Status                   StoreStatus `json:"status"`
	LogoImage                string      `json:"logo_image"`
	CoverImage               string      `json:"cover_image"`
	CommissionRate           int64       `json:"commission_rate"`
	IsProductCreationEnabled bool        `json:"is_product_creation_enabled"`
	IsOrderCreationEnabled   bool        `json:"is_order_creation_enabled"`
	IsAutoConfirmEnabled     bool        `json:"is_auto_confirm_enabled"`
	Description              string      `json:"description"`
	Address                  string      `json:"address"`
	City                     string      `json:"city"`
	Country                  string      `json:"country"`
	Postcode                 string      `json:"postcode"`
	Email                    string      `json:"email"`
	Phone                    string      `json:"phone"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
}

func (sv *StoreView) TableName() string {
	return "stores_view"
}

func (sv *StoreView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT s.id AS id, s.name AS name, s.status AS status, s.logo_image AS logo_image,"+
		" s.cover_image AS cover_image, s.commission_rate AS commission_rate, s.is_product_creation_enabled AS is_product_creation_enabled,"+
		" s.is_order_creation_enabled AS is_order_creation_enabled, s.is_auto_confirm_enabled AS is_auto_confirm_enabled,"+
		" s.description AS description, av.address AS address, av.city AS city, av.country AS country, av.postcode AS postcode,"+
		" av.email AS email, av.phone AS phone, s.created_at AS created_at, s.updated_at AS updated_at"+
		" FROM stores AS s"+
		" LEFT JOIN addresses_view AS av ON s.address_id = av.id", sv.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (sv *StoreView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", sv.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
