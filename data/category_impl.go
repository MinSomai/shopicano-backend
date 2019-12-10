package data

import (
	"github.com/shopicano/shopicano-backend/app"
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

func (cu *CategoryRepositoryImpl) CreateCategory(c *models.Category) error {
	db := app.DB()
	if err := db.Table(c.TableName()).Create(c).Error; err != nil {
		return err
	}
	return nil
}

func (cu *CategoryRepositoryImpl) ListCategories(from, limit int) ([]models.Category, error) {
	db := app.DB()
	var cols []models.Category
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("is_published = ?", true).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) ListCategoriesWithStore(storeID string, from, limit int) ([]models.Category, error) {
	db := app.DB()
	var cols []models.Category
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ?", storeID).
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) SearchCategories(query string, from, limit int) ([]models.Category, error) {
	db := app.DB()
	var cols []models.Category
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("is_published = ? AND LOWER(name) LIKE ?", true, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) SearchCategoriesWithStore(storeID, query string, from, limit int) ([]models.Category, error) {
	db := app.DB()
	var cols []models.Category
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND LOWER(name) LIKE ?", storeID, "%"+strings.ToLower(query)+"%").
		Offset(from).Limit(limit).
		Order("created_at DESC").Find(&cols).Error; err != nil {
		return nil, err
	}
	return cols, nil
}

func (cu *CategoryRepositoryImpl) DeleteCategory(storeID, categoryID string) error {
	db := app.DB()
	col := models.Category{}
	if err := db.Table(col.TableName()).
		Where("store_id = ? AND id = ?", storeID, categoryID).
		Delete(&col).Error; err != nil {
		return err
	}
	return nil
}
