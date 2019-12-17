package core

import "github.com/jinzhu/gorm"

type View interface {
	TableName() string
	CreateView(tx *gorm.DB) error
	DropView(tx *gorm.DB) error
}
