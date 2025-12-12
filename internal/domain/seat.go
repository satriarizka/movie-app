package domain

import "github.com/google/uuid"

type Seat struct {
	BaseModel
	StudioID   uuid.UUID `gorm:"type:uuid;not null" json:"studio_id"`
	RowCode    string    `gorm:"type:varchar(5);not null" json:"row_code"` // A, B, C
	SeatNumber int       `gorm:"type:int;not null" json:"seat_number"`     // 1, 2, 3

	// Relations
	Studio Studio `gorm:"foreignKey:StudioID" json:"-"`
}
