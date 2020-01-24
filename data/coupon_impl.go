package data

import (
	"github.com/jinzhu/gorm"
	"github.com/shopicano/shopicano-backend/models"
)

type CouponRepositoryImpl struct {
}

var couponRepository CouponRepository

func NewCouponRepository() CouponRepository {
	if couponRepository == nil {
		couponRepository = &CouponRepositoryImpl{}
	}
	return couponRepository
}

func (cr *CouponRepositoryImpl) Create(db *gorm.DB, c *models.Coupon) error {
	return db.Table(c.TableName()).Create(c).Error
}

func (cr *CouponRepositoryImpl) Update(db *gorm.DB, c *models.Coupon) error {
	return db.Table(c.TableName()).
		Where("store_id = ? AND id = ?", c.StoreID, c.ID).
		Update(map[string]interface{}{
			"code":             c.Code,
			"is_active":        c.IsActive,
			"discount_amount":  c.DiscountAmount,
			"is_flat_discount": c.IsFlatDiscount,
			"is_user_specific": c.IsUserSpecific,
			"max_discount":     c.MaxDiscount,
			"max_usage":        c.MaxUsage,
			"discount_type":    c.DiscountType,
			"start_at":         c.StartAt,
			"end_at":           c.EndAt,
			"updated_at":       c.UpdatedAt,
		}).Error
}

func (cr *CouponRepositoryImpl) List(db *gorm.DB, storeID string, from, limit int) ([]models.Coupon, error) {
	var coupons []models.Coupon
	c := models.Coupon{}
	if err := db.Table(c.TableName()).
		Where("store_id = ?", storeID).
		Order("updated_at DESC").
		Offset(from).
		Limit(limit).
		Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (cr *CouponRepositoryImpl) Search(db *gorm.DB, storeID, query string, from, limit int) ([]models.Coupon, error) {
	var coupons []models.Coupon
	c := models.Coupon{}
	if err := db.Table(c.TableName()).
		Where("store_id = ? AND code LIKE ?", storeID, "%"+query+"%").
		Order("updated_at DESC").
		Offset(from).
		Limit(limit).
		Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (cr *CouponRepositoryImpl) Delete(db *gorm.DB, storeID, couponID string) error {
	c := models.Coupon{}
	return db.Table(c.TableName()).
		Where("store_id = ? AND id = ?", storeID, couponID).
		Delete(&c).Error
}

func (cr *CouponRepositoryImpl) Get(db *gorm.DB, storeID, couponID string) (*models.Coupon, error) {
	c := models.Coupon{}
	if err := db.Table(c.TableName()).
		Where("store_id = ? AND id = ?", storeID, couponID).
		Find(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (cr *CouponRepositoryImpl) GetByCode(db *gorm.DB, storeID, couponCode string) (*models.Coupon, error) {
	c := models.Coupon{}
	if err := db.Table(c.TableName()).
		Where("store_id = ? AND code = ?", storeID, couponCode).
		Delete(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (cr *CouponRepositoryImpl) AddUser(db *gorm.DB, cf *models.CouponFor) error {
	if err := db.Table(cf.TableName()).
		Create(cf).Error; err != nil {
		return err
	}
	return nil
}

func (cr *CouponRepositoryImpl) RemoveUser(db *gorm.DB, cf *models.CouponFor) error {
	if err := db.Table(cf.TableName()).
		Delete(cf, "coupon_id = ? AND user_id = ?", cf.CouponID, cf.UserID).Error; err != nil {
		return err
	}
	return nil
}
