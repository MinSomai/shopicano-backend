package models

import (
	"fmt"
	"time"
)

type PayoutSendStatus string

func (pss PayoutSendStatus) IsValid() bool {
	for _, v := range []PayoutSendStatus{PayoutSendStatusPending, PayoutSendStatusConfirmed,
		PayoutSendStatusProcessing, PayoutSendStatusFailed, PayoutSendStatusCompleted} {
		if pss == v {
			return true
		}
	}
	return false
}

const (
	PayoutSendStatusPending    PayoutSendStatus = "payout_pending"
	PayoutSendStatusConfirmed  PayoutSendStatus = "payout_confirmed"
	PayoutSendStatusProcessing PayoutSendStatus = "payout_processing"
	PayoutSendStatusFailed     PayoutSendStatus = "payout_failed"
	PayoutSendStatusCompleted  PayoutSendStatus = "payout_completed"
)

type PayoutSend struct {
	ID                     string           `json:"id" gorm:"column:id;primary_key"`
	StoreID                string           `json:"store_id" gorm:"column:store_id;index"`
	InitiatedByUserID      string           `json:"initiated_by_user_id" gorm:"column:initiated_by_user_id;index"`
	IsMarketplaceInitiated bool             `json:"is_marketplace_initiated" gorm:"column:is_marketplace_initiated;index"`
	Status                 PayoutSendStatus `json:"status" gorm:"column:status;index"`
	Amount                 int64            `json:"amount" gorm:"column:amount;index"`
	FailureReason          string           `json:"failure_reason" gorm:"column:failure_reason"`
	Note                   string           `json:"note" gorm:"column:note"`
	Highlights             string           `json:"highlights" gorm:"column:highlights"`
	PayoutMethodID         string           `json:"payout_method_id" gorm:"column:payout_method_id"`
	PayoutMethodDetails    string           `json:"payout_method_details" gorm:"column:payout_method_details"`
	CreatedAt              time.Time        `json:"created_at" gorm:"column:created_at;index"`
	UpdatedAt              time.Time        `json:"updated_at" gorm:"column:updated_at"`
}

func (pst *PayoutSend) TableName() string {
	return "payout_sends"
}

func (pst *PayoutSend) ForeignKeys() []string {
	s := Store{}
	u := User{}
	pom := PayoutMethod{}

	return []string{
		fmt.Sprintf("store_id;%s(id);RESTRICT;RESTRICT", s.TableName()),
		fmt.Sprintf("initiated_by_user_id;%s(id);RESTRICT;RESTRICT", u.TableName()),
		fmt.Sprintf("payout_method_id;%s(id);RESTRICT;RESTRICT", pom.TableName()),
	}
}

type PayoutSendDetails struct {
	ID                     string           `json:"id"`
	StoreID                string           `json:"store_id"`
	InitiatedByUserID      string           `json:"initiated_by_user_id"`
	IsMarketplaceInitiated bool             `json:"is_marketplace_initiated"`
	Status                 PayoutSendStatus `json:"status"`
	Amount                 int64            `json:"amount"`
	FailureReason          string           `json:"failure_reason"`
	Note                   string           `json:"note"`
	Highlights             string           `json:"highlights"`
	PayoutMethodID         string           `json:"payout_method_id"`
	PayoutMethodName       string           `json:"payout_method_name"`
	PayoutMethodInputs     string           `json:"payout_method_inputs"`
	PayoutMethodDetails    string           `json:"payout_method_details"`
	CreatedAt              time.Time        `json:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at"`
}

func (psd *PayoutSendDetails) TableName() string {
	ps := PayoutSend{}
	return ps.TableName()
}
