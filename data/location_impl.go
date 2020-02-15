package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type LocationRepositoryImpl struct {
}

var locationRepository LocationRepository

func NewLocationRepository() LocationRepository {
	if locationRepository == nil {
		locationRepository = &LocationRepositoryImpl{}
	}
	return locationRepository
}

func (l *LocationRepositoryImpl) List(db *gorm.DB, query string, args []interface{}) ([]models.Location, error) {
	loc := models.Location{}
	var locations []models.Location
	if err := db.Table(loc.TableName()).
		Where(query, args...).
		Order("id ASC").
		Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}
