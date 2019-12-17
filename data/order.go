package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type OrderRepository interface {
	Create(db *gorm.DB, o *models.Order) error
	AddOrderedItem(db *gorm.DB, item *models.OrderedItem) error
	GetDetailsInternal(db *gorm.DB, orderID string) (*models.OrderDetailsInternal, error)
	GetDetails(db *gorm.DB, userID, orderID string) (*models.OrderDetails, error)
}
