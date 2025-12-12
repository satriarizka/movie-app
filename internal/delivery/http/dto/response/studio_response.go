package response

import "github.com/google/uuid"

type StudioResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
}
