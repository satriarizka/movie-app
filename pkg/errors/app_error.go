package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Implementasi interface error bawaan Go
func (e *AppError) Error() string {
	return e.Message
}

// Helper Functions untuk membuat error umum

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}
