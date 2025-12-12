package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	BaseModel
	StudioID  uuid.UUID `gorm:"type:uuid;not null" json:"studio_id"`
	MovieID   uuid.UUID `gorm:"type:uuid;not null" json:"movie_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`

	// Relations (Preload)
	Studio Studio `gorm:"foreignKey:StudioID" json:"studio,omitempty"`
	Movie  Movie  `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}
