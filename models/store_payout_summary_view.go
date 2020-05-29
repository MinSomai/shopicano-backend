package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type StorePayoutSummaryView struct {
	StoreID        string `json:"store_id"`
	TotalEarnings  int64  `json:"total_earnings"`
	TotalRequested int64  `json:"total_requested"`
	TotalPaid      int64  `json:"total_paid"`
	TotalAvailable int64  `json:"total_available"`
}

func (sps *StorePayoutSummaryView) TableName() string {
	return "store_payout_summaries"
}

func (sps *StorePayoutSummaryView) CreateView(tx *gorm.DB) error {
	sql := fmt.Sprintf("CREATE OR REPLACE VIEW %s AS SELECT ps.store_id, sfs.total_earnings AS total_earnings, "+
		"SUM(ps.amount) AS total_requested, SUM(ps.amount) FILTER (WHERE ps.status = 'payout_completed') AS total_paid, "+
		"(total_earnings - SUM(ps.amount)) AS total_available FROM payout_sends AS ps "+
		"JOIN store_finance_summaries AS sfs ON ps.store_id = sfs.store_id WHERE ps.status != 'payout_failed' "+
		"GROUP BY ps.store_id, sfs.total_earnings;", sps.TableName())
	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (sps *StorePayoutSummaryView) DropView(tx *gorm.DB) error {
	sql := fmt.Sprintf("DROP VIEW %s", sps.TableName())

	if err := tx.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
