package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type StoreUserProfile struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Address            string      `json:"address"`
	City               string      `json:"city"`
	Country            string      `json:"country"`
	PostCode           string      `json:"postcode"`
	Email              string      `json:"email"`
	Phone              string      `json:"phone"`
	Status             StoreStatus `json:"status"`
	Description        string      `json:"description"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	UserID             string      `json:"user_id"`
	UserName           string      `json:"user_name"`
	UserEmail          string      `json:"user_email"`
	UserProfilePicture string      `json:"user_profile_picture"`
	UserPhone          string      `json:"user_phone"`
	UserStatus         string      `json:"user_status"`
	UserPermission     Permission  `json:"user_permission"`
	StorePermission    Permission  `json:"store_permission"`
}

func (sup *StoreUserProfile) TableName() string {
	return "store_user_profiles"
}

func (sup *StoreUserProfile) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT s.id, s.name, s.address, s.city, s.country, s.postcode,"+
		" s.email, s.phone, s.status, s.description, s.created_at, up.permission AS user_permission,"+
		" s.updated_at, u.id AS user_id, u.name AS user_name, u.email AS user_email, u.profile_picture AS user_profile_picture,"+
		" u.phone AS user_phone, u.status AS user_status, sp.permission AS store_permission FROM stores AS s"+
		" LEFT JOIN staffs AS st ON s.id = st.store_id JOIN users AS u ON st.user_id = u.id"+
		" LEFT JOIN user_permissions AS up ON u.permission_id = up.id"+
		" LEFT JOIN store_permissions AS sp ON sp.id = st.permission_id;", sup.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (sup *StoreUserProfile) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", sup.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
