package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"movie-app/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transUC usecase.TransactionUseCase
	val     *validator.CustomValidator
}

func NewTransactionHandler(transUC usecase.TransactionUseCase, val *validator.CustomValidator) *TransactionHandler {
	return &TransactionHandler{transUC, val}
}

// PayTransaction godoc
// @Summary      Pay transaction
// @Description  Pay for a pending transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id       path    string  true  "Transaction UUID"
// @Param        request  body    request.PayTransactionRequest true "Payment Method"
// @Success      200      {object} utils.APIResponse
// @Failure      400      {object} utils.APIResponse "Already paid / Validation error"
// @Router       /transactions/{id}/pay [post]
// @Security     BearerAuth
func (h *TransactionHandler) PayTransaction(c *gin.Context) {
	// 1. Ambil User ID dari Token
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	// 2. Ambil Transaction ID dari URL
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	// 3. Bind Request
	var req request.PayTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	// 4. Proses Pembayaran
	if err := h.transUC.PayTransaction(userID, transactionID, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment successful", nil)
}

// CancelTransaction godoc
// @Summary      Cancel transaction
// @Description  Cancel a pending transaction manually
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction UUID"
// @Success      200  {object}  utils.APIResponse
// @Router       /transactions/{id}/cancel [post]
// @Security     BearerAuth
func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	// 1. Ambil User ID
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	// 2. Ambil Transaction ID
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	// 3. Panggil UseCase
	if err := h.transUC.CancelTransaction(userID, transactionID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction cancelled successfully", nil)
}

func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	transactions, err := h.transUC.GetUserTransactions(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "My transactions", transactions)
}
