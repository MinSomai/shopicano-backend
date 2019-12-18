package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
)

type OrderRepositoryImpl struct {
}

var orderRepository OrderRepository

func NewOrderRepository() OrderRepository {
	if orderRepository == nil {
		orderRepository = &OrderRepositoryImpl{}
	}
	return orderRepository
}

func (os *OrderRepositoryImpl) Create(db *gorm.DB, o *models.Order) error {
	if err := db.Table(o.TableName()).Create(o).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) UpdatePaymentInfo(db *gorm.DB, o *models.OrderDetailsView) error {
	order := models.Order{}
	if err := db.Table(order.TableName()).
		Select("nonce, transaction_id, is_paid, status, paid_at").
		Updates(map[string]interface{}{
			"nonce":          o.Nonce,
			"transaction_id": o.TransactionID,
			"is_paid":        o.IsPaid,
			"status":         o.Status,
			"paid_at":        o.PaidAt,
		}).Where("id = ?", o.ID).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) AddOrderedItem(db *gorm.DB, oi *models.OrderedItem) error {
	if err := db.Table(oi.TableName()).Create(oi).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) GetDetails(db *gorm.DB, orderID string) (*models.OrderDetailsView, error) {
	order := models.OrderDetailsView{}
	if err := db.Model(&order).First(&order, "id = ?", orderID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	oiv := models.OrderedItemView{}

	var items []models.OrderedItemView
	if err := db.Model(oiv).Find(&items, "order_id = ?", orderID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	order.Items = items
	return &order, nil
}

func (os *OrderRepositoryImpl) GetDetailsExternal(db *gorm.DB, userID, orderID string) (*models.OrderDetailsViewExternal, error) {
	order := models.OrderDetailsViewExternal{}
	if err := db.Model(&order).First(&order, "id = ?", orderID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	oiv := models.OrderedItemView{}

	var items []models.OrderedItemViewExternal
	if err := db.Model(oiv).Find(&items, "order_id = ?", orderID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	order.Items = items
	return &order, nil
}
