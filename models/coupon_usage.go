package models

type CouponUsage struct {
	ID       string `json:"id" gorm:"column:id;primary_key"`
	CouponID string `json:"coupon_id" gorm:"column:coupon_id;index"`
	UserID   string `json:"user_id" gorm:"column:user_id;index"`
}

func (du *CouponUsage) TableName() string {
	return "coupon_usages"
}
