package response

import "github.com/google/uuid"

type MovieResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Genre       string    `json:"genre"`
	PosterURL   string    `json:"poster_url"`
}
