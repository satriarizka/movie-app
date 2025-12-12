package domain

import (
	"time"

	"movie-app/internal/enums"

	"github.com/google/uuid"
)

type Transaction struct {
	BaseModel
	UserID        uuid.UUID               `gorm:"type:uuid;not null" json:"user_id"`
	TotalAmount   float64                 `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status        enums.TransactionStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaymentMethod string                  `gorm:"type:varchar(50)" json:"payment_method"`
	ExpiresAt     time.Time               `gorm:"-" json:"expires_at"`

	// --- Tambahan Field Promo ---
	PromoID        *uuid.UUID `gorm:"type:uuid" json:"promo_id"` // Pointer karena bisa null
	DiscountAmount float64    `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	FinalAmount    float64    `gorm:"type:decimal(10,2);not null" json:"final_amount"` // Total setelah diskon

	// Relations
	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tickets []Ticket `gorm:"foreignKey:TransactionID" json:"tickets,omitempty"`
	Promo   *Promo   `gorm:"foreignKey:PromoID" json:"promo,omitempty"`

	ReminderSent bool `gorm:"default:false" json:"reminder_sent"`
}
