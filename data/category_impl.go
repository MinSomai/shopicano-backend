package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
	"strings"
)

type CategoryRepositoryImpl struct {
}

var categoryRepository CategoryRepository

func NewCategoryRepository() CategoryRepository {
	if categoryRepository == nil {
		categoryRepository = &CategoryRepositoryImpl{}
	}

	return categoryRepository
}

func (cu *CategoryRepositoryImpl) Create(db *gorm.DB, c *models.Category) error {
	if err := db.Table(c.TableName()).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CategoryRepositoryImpl) List(db *gorm.DB, from, limit int) ([]models.ResCategorySearch, error) {
	var cols []models.ResCategorySearch
	col := models.Category{}
	if err := db.Table(fmt.Sprintf("%s AS c", col.TableName())).
		Select("COUNT(p.id) AS count, c.id, c.name, c.description, c.image, c.store_id").
		Joins("LEFT JOIN products AS p ON p.category_id = c.id").
		Group("c.name, c.description, c.image, c.id, c.updated_at, c.store_id").
		Where("c.is_published = ?", true).
		Offset(from).Limit(limit).
		Order("c.updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) ListAsStoreStuff(db *gorm.DB, storeID string, from, limit int) ([]models.ResCategorySearchInternal, error) {
	var cols []models.ResCategorySearchInternal
	col := models.Category{}
	if err := db.Table(fmt.Sprintf("%s AS c", col.TableName())).
		Select("COUNT(p.id) AS count, c.id, c.name, c.description, c.image, c.store_id, c.created_at, c.is_published, c.updated_at").
		Joins("LEFT JOIN products AS p ON p.category_id = c.id").
		Group("c.name, c.description, c.image, c.id, c.updated_at, c.store_id, c.is_published, c.created_at").
		Where("c.store_id = ?", storeID).
		Offset(from).Limit(limit).
		Order("c.updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) Search(db *gorm.DB, query string, from, limit int) ([]models.ResCategorySearch, error) {
	var cols []models.ResCategorySearch
	col := models.Category{}
	if err := db.Table(fmt.Sprintf("%s AS c", col.TableName())).
		Select("COUNT(p.id) AS count, c.id, c.name, c.description, c.image, c.store_id").
		Joins("LEFT JOIN products AS p ON p.category_id = c.id").
		Group("c.name, c.description, c.image, c.id, c.updated_at, c.store_id").
		Where("c.is_published = ? AND LOWER(c.name) LIKE ?", true, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("c.updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) SearchAsStoreStuff(db *gorm.DB, storeID, query string, from, limit int) ([]models.ResCategorySearchInternal, error) {
	var cols []models.ResCategorySearchInternal
	col := models.Category{}
	if err := db.Table(fmt.Sprintf("%s AS c", col.TableName())).
		Select("COUNT(p.id) AS count, c.id, c.name, c.description, c.image, c.store_id, c.created_at, c.is_published, c.updated_at").
		Joins("LEFT JOIN products AS p ON p.category_id = c.id").
		Group("c.name, c.description, c.image, c.id, c.updated_at, c.store_id, c.is_published, c.created_at").
		Where("c.store_id = ? AND LOWER(c.name) LIKE ?", storeID, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("c.updated_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) Delete(db *gorm.DB, storeID, categoryID string) error {
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, categoryID).
		Delete(&col).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CategoryRepositoryImpl) Get(db *gorm.DB, storeID, categoryID string) (*models.Category, error) {
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, categoryID).
		First(&col).Error; err != nil {
		return nil, err
	}
	return &col, nil
}

func (cu *CategoryRepositoryImpl) Update(db *gorm.DB, c *models.Category) error {
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", c.StoreID, c.ID).
		Select("name", "description", "image", "is_published", "updated_at").
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
