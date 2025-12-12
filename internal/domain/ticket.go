package domain

import "github.com/google/uuid"

type Ticket struct {
	BaseModel
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	ScheduleID    uuid.UUID `gorm:"type:uuid;not null" json:"schedule_id"`
	SeatID        uuid.UUID `gorm:"type:uuid;not null" json:"seat_id"`

	// Relations
	Seat     Seat     `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"-"`
}
