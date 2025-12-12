package request

type CreateMovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    int    `json:"duration" validate:"required,min=1"` // Menit
	Genre       string `json:"genre" validate:"required"`
	PosterURL   string `json:"poster_url" validate:"required,url"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration" validate:"omitempty,min=1"`
	Genre       string `json:"genre"`
	PosterURL   string `json:"poster_url" validate:"omitempty,url"`
}
