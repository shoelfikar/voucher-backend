package entity

import (
	"time"

	"gorm.io/gorm"
)

// Voucher represents a voucher in the system
type Voucher struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	VoucherCode     string         `gorm:"uniqueIndex;not null;size:50" json:"voucher_code"`
	DiscountPercent float64        `gorm:"not null;check:discount_percent >= 1 AND discount_percent <= 100" json:"discount_percent"`
	ExpiryDate      time.Time      `gorm:"not null;type:date" json:"expiry_date"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Voucher entity
func (Voucher) TableName() string {
	return "vouchers"
}
