package request

type CreateStudioRequest struct {
	Name     string `json:"name" validate:"required"`
	Capacity int    `json:"capacity" validate:"required,min=1"`
}

type UpdateStudioRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity" validate:"omitempty,min=1"`
}
