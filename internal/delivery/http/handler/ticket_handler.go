package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	_ "movie-app/internal/delivery/http/dto/response"
	_ "movie-app/internal/domain"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"movie-app/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketUC usecase.TicketUseCase
	val      *validator.CustomValidator
}

func NewTicketHandler(ticketUC usecase.TicketUseCase, val *validator.CustomValidator) *TicketHandler {
	return &TicketHandler{ticketUC, val}
}

// GetAvailableSeats godoc
// @Summary      Get available seats
// @Description  Check which seats are booked or free for a schedule
// @Tags         Ticketing
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule UUID"
// @Success      200  {object}  utils.APIResponse{data=[]response.SeatAvailabilityResponse}
// @Router       /tickets/schedules/{id}/seats [get]
// @Security     BearerAuth
func (h *TicketHandler) GetAvailableSeats(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	seats, err := h.ticketUC.GetAvailableSeats(scheduleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Available seats", seats)
}

// BookTicket godoc
// @Summary      Book tickets
// @Description  Book seats for a specific schedule
// @Tags         Ticketing
// @Accept       json
// @Produce      json
// @Param        request body request.BookTicketRequest true "Booking Data"
// @Success      201  {object}  utils.APIResponse{data=domain.Transaction}
// @Failure      409  {object}  utils.APIResponse "Conflict / Double Booking"
// @Router       /tickets/book [post]
// @Security     BearerAuth
func (h *TicketHandler) BookTicket(c *gin.Context) {
	// 1. Ambil User ID dari Token (set via Middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID, _ := uuid.Parse(userIDStr.(string))

	// 2. Parse Body
	var req request.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	// 3. Call UseCase
	transaction, err := h.ticketUC.BookTicket(userID, req)
	if err != nil {
		// Kemungkinan besar conflict (double booking)
		utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Booking successful, waiting for payment", transaction)
}

// GetUserHistory godoc
// @Summary      Get booking history
// @Description  Get all transaction history for current user
// @Tags         Ticketing
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]domain.Transaction}
// @Router       /tickets/me [get]
// @Security     BearerAuth
func (h *TicketHandler) GetUserHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	history, err := h.ticketUC.GetUserHistory(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User booking history", history)
}
