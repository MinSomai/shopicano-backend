package models

import "fmt"

type CouponUsage struct {
	CouponID string `json:"coupon_id" gorm:"column:coupon_id;primary_key"`
	UserID   string `json:"user_id" gorm:"column:user_id;primary_key"`
}

func (cu *CouponUsage) TableName() string {
	return "coupon_usages"
}

func (cu *CouponUsage) ForeignKeys() []string {
	c := Coupon{}
	u := User{}

	return []string{
		fmt.Sprintf("coupon_id;%s(id);RESTRICT;RESTRICT", c.TableName()),
		fmt.Sprintf("user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
	}
}
