package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type LocationRepository interface {
	List(db *gorm.DB, query string, args []interface{}) ([]models.Location, error)
	UpdateByID(db *gorm.DB, locationID int64, toggle int64) error
	UpdateAll(db *gorm.DB, toggle int64) error
}
