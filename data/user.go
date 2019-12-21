package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type UserRepository interface {
	Register(db *gorm.DB, u *models.User) error
	Login(db *gorm.DB, email, password string) (*models.Session, error)
	Logout(db *gorm.DB, token string) error
	RefreshToken(db *gorm.DB, token string) (*models.Session, error)
	Update(db *gorm.DB, u *models.User) error
	GetPermission(db *gorm.DB, token string) (string, *models.Permission, error)
	GetPermissionByUserID(db *gorm.DB, userID string) (string, *models.Permission, error)
	Get(db *gorm.DB, userID string) (*models.User, error)
	IsSignUpEnabled(db *gorm.DB) (bool, error)
	IsStoreCreationEnabled(db *gorm.DB) (bool, error)
}
