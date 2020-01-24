package models

type CouponFor struct {
	CouponID string `json:"coupon_id" gorm:"column:coupon_id;primary_key"`
	UserID   string `json:"user_id" gorm:"column:user_id;primary_key"`
}

func (df *CouponFor) TableName() string {
	return "coupon_for"
}
