package data

import (
	"github.com/jinzhu/gorm"
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

func (cu *CollectionRepositoryImpl) Create(db *gorm.DB, c *models.Collection) error {
	if err := db.Table(c.TableName()).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CollectionRepositoryImpl) List(db *gorm.DB, from, limit int) ([]models.Collection, error) {
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.Collection, error) {
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ?", storeID).
		Offset(from).Limit(limit).
		Order("updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) Search(db *gorm.DB, query string, from, limit int) ([]models.Collection, error) {
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("is_published = ? AND LOWER(name) LIKE ?", true, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.Collection, error) {
	var cols []models.Collection
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND LOWER(name) LIKE ?", storeID, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CollectionRepositoryImpl) Delete(db *gorm.DB, storeID, collectionID string) error {
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, collectionID).
		Delete(&col).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CollectionRepositoryImpl) Get(db *gorm.DB, storeID, collectionID string) (*models.Collection, error) {
	col := models.Collection{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, collectionID).
		First(&col).Error; err != nil {
		return nil, err
	}
	return &col, nil
}

func (cu *CollectionRepositoryImpl) Update(db *gorm.DB, c *models.Collection) error {
	if err := db.Table(c.TableName()).
		Where("id = ?", c.ID).
		Select("name", "description", "is_published", "image", "updated_at").
		Updates(map[string]interface{}{
			"name":         c.Name,
			"description":  c.Description,
			"is_published": c.IsPublished,
			"image":        c.Image,
			"updated_at":   c.UpdatedAt,
		}).Error; err != nil {
		return err
	}
	return nil
}
