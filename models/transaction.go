package models

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MSISDN      string         `json:"msisdn" gorm:"not null;index"`
	TrxID       string         `json:"trx_id" gorm:"not null;unique"`
	TrxDate     time.Time      `json:"trx_date" gorm:"not null;index"`
	Item        string         `json:"item" gorm:"not null"`
	VoucherCode string         `json:"voucher_code" gorm:"not null"`
	Status      int            `json:"status" gorm:"not null;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Transaction) TableName() string {
	return "transactions"
}
