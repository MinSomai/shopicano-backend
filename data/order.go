package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
	"time"
)

type OrderRepository interface {
	Create(db *gorm.DB, o *models.Order) error
	AddOrderedItem(db *gorm.DB, item *models.OrderedItem) error
	GetDetailsAsUser(db *gorm.DB, userID, orderID string) (*models.OrderDetailsViewExternal, error)
	GetDetailsAsStoreStuff(db *gorm.DB, storeID, orderID string) (*models.OrderDetailsView, error)
	GetDetails(db *gorm.DB, orderID string) (*models.OrderDetailsView, error)
	UpdatePaymentInfo(db *gorm.DB, o *models.OrderDetailsView) error
	List(db *gorm.DB, userID string, offset, limit int) ([]models.OrderDetailsViewExternal, error)
	ListAsStoreStuff(db *gorm.DB, storeID string, offset, limit int) ([]models.OrderDetailsViewExternal, error)
	Search(db *gorm.DB, query, userID string, offset, limit int) ([]models.OrderDetailsView, error)
	SearchAsStoreStuff(db *gorm.DB, query, storeID string, offset, limit int) ([]models.OrderDetailsView, error)
	CountByTimeAsStoreStuff(db *gorm.DB, storeID string, from, end time.Time) (int, error)
}
