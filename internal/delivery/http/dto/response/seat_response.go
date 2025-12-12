package response

import "github.com/google/uuid"

type SeatAvailabilityResponse struct {
	ID         uuid.UUID `json:"id"`
	RowCode    string    `json:"row_code"`
	SeatNumber int       `json:"seat_number"`
	IsBooked   bool      `json:"is_booked"` // True jika sudah ada yang punya
}
