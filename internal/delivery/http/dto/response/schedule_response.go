package response

import (
	"time"

	"github.com/google/uuid"
)

// Nested struct untuk response yang rapi
type ScheduleResponse struct {
	ID        uuid.UUID      `json:"id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Price     float64        `json:"price"`
	Studio    StudioResponse `json:"studio"`
	Movie     MovieResponse  `json:"movie"`
}
