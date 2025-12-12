package response

import "github.com/google/uuid"

type AuthResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
