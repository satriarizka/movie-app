package utils

import (
	"movie-app/pkg/errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {

	// Cek apakah errornya adalah tipe AppError buatan kita
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, APIResponse{
			Status:  false,
			Message: appErr.Message,
			Errors:  nil, // Atau appErr.Err.Error() jika mau debug
		})
		return
	}

	// Jika error biasa
	c.JSON(statusCode, APIResponse{
		Status:  false,
		Message: message,
		Errors:  err,
	})
}

// PaginationMeta untuk response list
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}
