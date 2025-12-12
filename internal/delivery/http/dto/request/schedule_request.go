package request

import "time"

type CreateScheduleRequest struct {
	StudioID  string    `json:"studio_id" validate:"required,uuid"`
	MovieID   string    `json:"movie_id" validate:"required,uuid"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Price     float64   `json:"price" validate:"required,min=0"`
}

type UpdateScheduleRequest struct {
	StudioID  string    `json:"studio_id" validate:"omitempty,uuid"`
	MovieID   string    `json:"movie_id" validate:"omitempty,uuid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time" validate:"omitempty,gtfield=StartTime"`
	Price     float64   `json:"price" validate:"omitempty,min=0"`
}
