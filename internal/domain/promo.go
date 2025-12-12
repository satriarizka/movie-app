package domain

import (
	"time"
)

type Promo struct {
	BaseModel
	Code          string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	DiscountType  string    `gorm:"type:varchar(20);not null" json:"discount_type"` // 'percentage' or 'fixed'
	DiscountValue float64   `gorm:"type:decimal(10,2);not null" json:"discount_value"`
	ValidUntil    time.Time `gorm:"not null" json:"valid_until"`
}
