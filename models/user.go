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

func (us UserStatus) IsValid() bool {
	for _, s := range []UserStatus{UserActive, UserBanned, UserSuspended} {
		if s == us {
			return true
		}
	}
	return false
}

type User struct {
	ID                            string     `json:"id" gorm:"column:id;primary_key"`
	Name                          string     `json:"name" gorm:"column:name;not null"`
	Email                         string     `json:"email" gorm:"column:email;unique;not null"`
	ProfilePicture                *string    `json:"profile_picture,omitempty" gorm:"column:profile_picture"`
	Phone                         *string    `json:"phone,omitempty" gorm:"column:phone;unique"`
	Password                      string     `json:"-" gorm:"column:password;not null"`
	VerificationToken             *string    `json:"-" gorm:"column:verification_token;unique"`
	ResetPasswordToken            *string    `json:"-" gorm:"column:reset_password_token;index"`
	ResetPasswordTokenGeneratedAt *time.Time `json:"-" gorm:"column:reset_password_token_generated_at"`
	Status                        UserStatus `json:"status" gorm:"column:status;index;not null"`
	IsEmailVerified               bool       `json:"is_email_verified" gorm:"column:is_email_verified"`
	PermissionID                  string     `json:"-" gorm:"column:permission_id;index;not null"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt                     time.Time  `json:"updated_at" gorm:"column:updated_at"`
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
	ID           string    `json:"-" gorm:"column:id;primary_key"`
	UserID       string    `json:"-" gorm:"column:user_id;index;not null"`
	AccessToken  string    `json:"access_token" gorm:"column:access_token;unique;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"column:refresh_token;unique;not null"`
	CreatedAt    time.Time `json:"-" gorm:"column:created_at;index;not null"`
	ExpireOn     int64     `json:"expire_on" gorm:"column:expire_on;index;not null"`
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
	ID         string     `json:"id" gorm:"column:id;primary_key"`
	Permission Permission `json:"permission" gorm:"column:permission;index;not null"`
}

func (up *UserPermission) TableName() string {
	return "user_permissions"
}
