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

func (os *OrderRepositoryImpl) List(db *gorm.DB, userID string, offset, limit int) ([]models.OrderDetailsViewExternal, error) {
	order := models.OrderDetailsViewExternal{}
	var orders []models.OrderDetailsViewExternal

	if err := db.Table(order.TableName()).Offset(offset).Limit(limit).Find(&orders, "user_id = ?", userID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	for i, v := range orders {
		oiv := models.OrderedItemViewExternal{}

		var items []models.OrderedItemViewExternal
		if err := db.Table(oiv.TableName()).Find(&items, "order_id = ?", v.ID).Error; err != nil {
			log.Log().Errorln(err)
			return nil, err
		}

		if len(items) == 0 {
			items = []models.OrderedItemViewExternal{}
		}

		orders[i].Items = items
	}

	if len(orders) == 0 {
		orders = []models.OrderDetailsViewExternal{}
	}
	return orders, nil
}

func (os *OrderRepositoryImpl) ListAsStoreStuff(db *gorm.DB, storeID string, offset, limit int) ([]models.OrderDetailsViewExternal, error) {
	order := models.OrderDetailsViewExternal{}
	var orders []models.OrderDetailsViewExternal

	if err := db.Model(&order).Offset(offset).Limit(limit).Find(&orders, "store_id = ?", storeID).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	for i, v := range orders {
		oiv := models.OrderedItemViewExternal{}

		var items []models.OrderedItemViewExternal
		if err := db.Model(oiv).Find(&items, "order_id = ?", v.ID).Error; err != nil {
			log.Log().Errorln(err)
			return nil, err
		}

		if len(items) == 0 {
			items = []models.OrderedItemViewExternal{}
		}

		orders[i].Items = items
	}

	if len(orders) == 0 {
		orders = []models.OrderDetailsViewExternal{}
	}
	return orders, nil
}

func (os *OrderRepositoryImpl) Search(db *gorm.DB, query, userID string, offset, limit int) ([]models.OrderDetailsView, error) {
	order := models.OrderDetailsView{}
	var orders []models.OrderDetailsView

	if err := db.Model(&order).Offset(offset).Limit(limit).Find(&orders, "user_id = ? AND hash = ?", userID, query).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	for i, v := range orders {
		oiv := models.OrderedItemView{}

		var items []models.OrderedItemView
		if err := db.Model(oiv).Find(&items, "order_id = ?", v.ID).Error; err != nil {
			log.Log().Errorln(err)
			return nil, err
		}

		if len(items) == 0 {
			items = []models.OrderedItemView{}
		}

		orders[i].Items = items
	}

	if len(orders) == 0 {
		orders = []models.OrderDetailsView{}
	}
	return orders, nil
}

func (os *OrderRepositoryImpl) SearchAsStoreStuff(db *gorm.DB, query, storeID string, offset, limit int) ([]models.OrderDetailsView, error) {
	order := models.OrderDetailsView{}
	var orders []models.OrderDetailsView

	if err := db.Model(&order).Offset(offset).Limit(limit).Find(&orders, "store_id = ? AND hash = ?", storeID, query).Error; err != nil {
		log.Log().Errorln(err)
		return nil, err
	}

	for i, v := range orders {
		oiv := models.OrderedItemView{}

		var items []models.OrderedItemView
		if err := db.Model(oiv).Find(&items, "order_id = ?", v.ID).Error; err != nil {
			log.Log().Errorln(err)
			return nil, err
		}

		if len(items) == 0 {
			items = []models.OrderedItemView{}
		}

		orders[i].Items = items
	}

	if len(orders) == 0 {
		orders = []models.OrderDetailsView{}
	}
	return orders, nil
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
