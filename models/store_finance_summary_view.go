package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type StoreFinanceSummaryView struct {
	StoreID         string `json:"store_id"`
	TotalIncome     int64  `json:"total_income"`
	TotalEarnings   int64  `json:"total_earnings"`
	TotalCommission int64  `json:"total_commission"`
}

func (sfs *StoreFinanceSummaryView) TableName() string {
	return "store_finance_summaries"
}

func (sfs *StoreFinanceSummaryView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT odv.store_id, SUM(odv.actual_earnings) AS total_income, "+
		"SUM(odv.seller_earnings) AS total_earnings, SUM(odv.platform_earnings) AS total_commission "+
		"FROM order_details_views AS odv WHERE odv.payment_status = 'payment_completed' GROUP BY odv.store_id;", sfs.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (sfs *StoreFinanceSummaryView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", sfs.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
