package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type LocationRepository interface {
	List(db *gorm.DB, query string, args []interface{}) ([]models.Location, error)
}
