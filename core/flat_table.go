package core

import "github.com/jinzhu/gorm"

type FlatTable interface {
	TableName() string
	Populate(db *gorm.DB) error
}
