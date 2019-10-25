package repositories

import "github.com/shopicano/shopicano-backend/models"

type UserRepository interface {
	Register(u *models.User) error
	Login(email, password string) (*models.Session, error)
	Logout(token string) error
	RefreshToken(token string) (*models.Session, error)
	Update(userID string, ud *models.User) (*models.User, error)
	GetPermission(token string) (string, *models.Permission, error)
	GetPermissionByUserID(userID string) (string, *models.Permission, error)
	Get(userID string) (*models.User, error)
	IsSignUpEnabled() (bool, error)
	IsStoreCreationEnabled() (bool, error)
}
