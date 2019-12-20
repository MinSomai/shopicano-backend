package models

import (
	"fmt"
	"time"
)

const (
	UserRegistered UserStatus = "registered"
	UserActive     UserStatus = "active"
	UserBanned     UserStatus = "banned"
	UserSuspended  UserStatus = "suspended"

	AdminPerm   Permission = "admin"
	ManagerPerm Permission = "manager"
	UserPerm    Permission = "user"
)

type UserStatus string
type Permission string

type User struct {
	ID                string     `json:"id" sql:"id" gorm:"primary_key"`
	Name              string     `json:"name" sql:"name" gorm:"not null"`
	Email             string     `json:"email" sql:"email" gorm:"unique;not null"`
	ProfilePicture    *string    `json:"profile_picture,omitempty" sql:"profile_picture"`
	Phone             *string    `json:"phone,omitempty" sql:"phone" gorm:"unique"`
	Password          string     `json:"-" sql:"password" gorm:"not null"`
	VerificationToken *string    `json:"-" sql:"verification_token" gorm:"unique"`
	Status            UserStatus `json:"status" sql:"status" gorm:"index;not null"`
	IsEmailVerified   bool       `json:"is_email_verified" json:"is_email_verified"`
	PermissionID      string     `json:"-" sql:"permission_id" gorm:"index;not null"`
	CreatedAt         time.Time  `json:"created_at" sql:"created_at" gorm:"index"`
	UpdatedAt         time.Time  `json:"updated_at" sql:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ForeignKeys() []string {
	up := UserPermission{}

	return []string{
		fmt.Sprintf("permission_id;%s(id);RESTRICT;RESTRICT", up.TableName()),
	}
}

type Session struct {
	ID           string    `json:"-" sql:"id" gorm:"primary_key"`
	UserID       string    `json:"-" sql:"user_id" gorm:"index;not null"`
	AccessToken  string    `json:"access_token" sql:"access_token" gorm:"unique;not null"`
	RefreshToken string    `json:"refresh_token" sql:"refresh_token" gorm:"unique;not null"`
	CreatedAt    time.Time `json:"-" sql:"created_at" gorm:"index;not null"`
	ExpireOn     int64     `json:"expire_on" sql:"expire_on" gorm:"index;not null"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func (s *Session) ForeignKeys() []string {
	u := User{}

	return []string{
		fmt.Sprintf("user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
	}
}

type UserPermission struct {
	ID         string     `json:"id" sql:"id" gorm:"primary_key"`
	Permission Permission `json:"permission" sql:"permission" gorm:"index;not null"`
}

func (up *UserPermission) TableName() string {
	return "user_permissions"
}
