package repositories

import "github.com/shopicano/shopicano-backend/models"

type CategoryRepository interface {
	CreateCategory(c *models.Category) error
	ListCategories(from, limit int) ([]models.Category, error)
	SearchCategories(query string, from, limit int) ([]models.Category, error)
	ListCategoriesWithStore(storeID string, from, limit int) ([]models.Category, error)
	SearchCategoriesWithStore(storeID, query string, from, limit int) ([]models.Category, error)
	DeleteCategory(storeID, categoryID string) error
}
