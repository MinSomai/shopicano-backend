package models

import (
	"fmt"
	"time"
)

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
