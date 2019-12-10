package data

import (
	"github.com/shopicano/shopicano-backend/app"
	"github.com/shopicano/shopicano-backend/models"
	"strings"
)

type CollectionRepositoryImpl struct {
}

var collectionRepository CollectionRepository

func NewCollectionRepository() CollectionRepository {
	if collectionRepository == nil {
		collectionRepository = &CollectionRepositoryImpl{}
	}

	return collectionRepository
}

func (cu *CollectionRepositoryImpl) CreateCollection(c *models.Collection) error {
	db := app.DB()
	if err := db.Table(c.TableName()).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CollectionRepositoryImpl) ListCollections(from, limit int) ([]models.Collection, error) {
	db := app.DB()
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) ListCollectionsWithStore(storeID string, from, limit int) ([]models.Collection, error) {
	db := app.DB()
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ?", storeID).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) SearchCollections(query string, from, limit int) ([]models.Collection, error) {
	db := app.DB()
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("is_published = ? AND LOWER(name) LIKE ?", true, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) SearchCollectionsWithStore(storeID, query string, from, limit int) ([]models.Collection, error) {
	db := app.DB()
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND LOWER(name) LIKE ?", storeID, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) DeleteCollection(storeID, collectionID string) error {
	db := app.DB()
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, collectionID).
		Delete(&col).Error; err != nil {
		return err
	}
	return nil
}
