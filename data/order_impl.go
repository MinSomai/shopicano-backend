package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/log"
	"github.com/shopicano/shopicano-backend/models"
	"time"
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

func (os *OrderRepositoryImpl) CreateLog(db *gorm.DB, ol *models.OrderLog) error {
	if err := db.Table(ol.TableName()).Create(ol).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) UpdatePaymentInfo(db *gorm.DB, o *models.OrderDetailsView) error {
	order := models.Order{}
	if err := db.Table(order.TableName()).
		Where("id = ?", o.ID).
		Select("nonce, transaction_id, payment_status").
		Updates(map[string]interface{}{
			"nonce":          o.Nonce,
			"transaction_id": o.TransactionID,
			"payment_status": o.PaymentStatus,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) UpdateStatus(db *gorm.DB, o *models.Order) error {
	order := models.Order{}
	if err := db.Table(order.TableName()).
		Where("id = ?", o.ID).
		Select("status").
		Updates(map[string]interface{}{
			"status": o.Status,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (os *OrderRepositoryImpl) UpdatePaymentStatus(db *gorm.DB, o *models.Order) error {
	order := models.Order{}
	if err := db.Table(order.TableName()).
		Where("id = ?", o.ID).
		Select("payment_status").
		Updates(map[string]interface{}{
			"payment_status": o.PaymentStatus,
		}).Error; err != nil {
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

func (os *OrderRepositoryImpl) GetOrderedItem(db *gorm.DB, orderID, productID string) (*models.OrderedItem, error) {
	oi := models.OrderedItem{}
	if err := db.Table(oi.TableName()).Find(&oi, "order_id = ? AND product_id = ?", orderID, productID).Error; err != nil {
		return nil, err
	}
	return &oi, nil
}

func (os *OrderRepositoryImpl) List(db *gorm.DB, userID string, offset, limit int) ([]models.OrderDetailsViewExternal, error) {
	order := models.OrderDetailsViewExternal{}
	var orders []models.OrderDetailsViewExternal

	if err := db.Table(order.TableName()).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders, "user_id = ?", userID).Error; err != nil {
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

	if err := db.Model(&order).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders, "store_id = ?", storeID).Error; err != nil {
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

	if err := db.Model(&order).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders, "user_id = ? AND hash = ?", userID, query).Error; err != nil {
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

	if err := db.Model(&order).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders, "store_id = ? AND hash = ?", storeID, query).Error; err != nil {
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

func (os *OrderRepositoryImpl) GetDetailsAsStoreStuff(db *gorm.DB, storeID, orderID string) (*models.OrderDetailsView, error) {
	order := models.OrderDetailsView{}
	if err := db.Model(&order).First(&order, "id = ? AND store_id = ?", orderID, storeID).Error; err != nil {
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

func (os *OrderRepositoryImpl) GetDetailsAsUser(db *gorm.DB, userID, orderID string) (*models.OrderDetailsViewExternal, error) {
	order := models.OrderDetailsViewExternal{}
	if err := db.Model(&order).First(&order, "id = ? AND user_id = ?", orderID, userID).Error; err != nil {
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

func (os *OrderRepositoryImpl) CountByTimeAsStoreStuff(db *gorm.DB, storeID string, from, end time.Time) (int, error) {
	order := models.OrderDetailsViewExternal{}

	var count int

	if err := db.Table(order.TableName()).
		Where("store_id = ? AND (created_at >= ? AND created_at <= ?)", storeID, from, end).
		Count(&count).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return count, nil
}

func (os *OrderRepositoryImpl) CountByTimeByStatus(db *gorm.DB, storeID string, from, end time.Time, status models.OrderStatus) (int, error) {
	order := models.OrderDetailsViewExternal{}

	var count int

	if err := db.Table(order.TableName()).
		Where("store_id = ? AND (created_at >= ? AND created_at <= ?) AND status = ?", storeID, from, end, status).
		Count(&count).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return count, nil
}

func (os *OrderRepositoryImpl) CountAsStoreStuff(db *gorm.DB, storeID string) (int, error) {
	order := models.OrderDetailsViewExternal{}

	var count int

	if err := db.Table(order.TableName()).
		Where("store_id = ?", storeID).
		Count(&count).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return count, nil
}

func (os *OrderRepositoryImpl) Earnings(db *gorm.DB, storeID string) (int, error) {
	order := models.OrderDetailsViewExternal{}

	r := struct {
		Earnings int `json:"earnings"`
	}{}

	if err := db.Table(order.TableName()).
		Select("SUM(order_details_views.grand_total) AS earnings").
		Where("store_id = ? AND payment_status = ?", storeID, models.PaymentCompleted).
		Scan(&r).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return r.Earnings, nil
}

func (os *OrderRepositoryImpl) EarningsByTime(db *gorm.DB, storeID string, from, end time.Time) (int, error) {
	order := models.OrderDetailsViewExternal{}

	r := struct {
		Earnings int `json:"earnings"`
	}{}

	if err := db.Table(order.TableName()).
		Select("SUM(order_details_views.grand_total) AS earnings").
		Where("store_id = ? AND (created_at >= ? AND created_at <= ?)", storeID, from, end).
		Scan(&r).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return r.Earnings, nil
}

func (os *OrderRepositoryImpl) EarningsByTimeByStatus(db *gorm.DB, storeID string, from, end time.Time, status models.PaymentStatus) (int, error) {
	order := models.OrderDetailsViewExternal{}

	r := struct {
		Earnings int `json:"earnings"`
	}{}

	if err := db.Table(order.TableName()).
		Select("SUM(order_details_views.grand_total) AS earnings").
		Where("store_id = ? AND (created_at >= ? AND created_at <= ?) AND payment_status = ?", storeID, from, end, status).
		Scan(&r).Error; err != nil {
		log.Log().Errorln(err)
		return 0, err
	}

	return r.Earnings, nil
}
