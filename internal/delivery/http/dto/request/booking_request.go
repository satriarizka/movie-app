package request

type BookTicketRequest struct {
	ScheduleID string   `json:"schedule_id" validate:"required,uuid"`
	SeatIDs    []string `json:"seat_ids" validate:"required,min=1,dive,uuid"` // Array of Seat UUID
	PromoCode  string   `json:"promo_code"`
}
