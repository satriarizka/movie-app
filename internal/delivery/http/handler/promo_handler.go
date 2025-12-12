package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	_ "movie-app/internal/domain"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PromoHandler struct {
	promoUC usecase.PromoUseCase
}

func NewPromoHandler(promoUC usecase.PromoUseCase) *PromoHandler {
	return &PromoHandler{promoUC}
}

// Create godoc
// @Summary      Create new promo
// @Description  Add a new promo code (Admin only)
// @Tags         Promos
// @Accept       json
// @Produce      json
// @Param        request body request.CreatePromoRequest true "Promo Data"
// @Success      201  {object}  utils.APIResponse{data=domain.Promo}
// @Failure      400  {object}  utils.APIResponse
// @Router       /promos [post]
// @Security     BearerAuth
func (h *PromoHandler) Create(c *gin.Context) {
	var req request.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	promo, err := h.promoUC.CreatePromo(req.Code, req.DiscountType, req.DiscountValue, req.ValidUntil)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Promo created", promo)
}

// GetAll godoc
// @Summary      Get all promos
// @Description  Get list of active promos (Admin only)
// @Tags         Promos
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]domain.Promo}
// @Router       /promos [get]
// @Security     BearerAuth
func (h *PromoHandler) GetAll(c *gin.Context) {
	promos, err := h.promoUC.GetAllPromos()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "List promos", promos)
}

// Update godoc
// @Summary      Update promo
// @Description  Update promo details (Admin only)
// @Tags         Promos
// @Accept       json
// @Produce      json
// @Param        id       path    string  true  "Promo UUID"
// @Param        request  body    request.UpdatePromoRequest true "Update Data"
// @Success      200      {object} utils.APIResponse{data=domain.Promo}
// @Router       /promos/{id} [put]
// @Security     BearerAuth
func (h *PromoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	var req request.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	promo, err := h.promoUC.UpdatePromo(id, req.Code, req.DiscountType, req.DiscountValue, req.ValidUntil)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Promo updated", promo)
}

// Delete godoc
// @Summary      Delete promo
// @Description  Delete/Remove a promo (Admin only)
// @Tags         Promos
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Promo UUID"
// @Success      200  {object}  utils.APIResponse
// @Router       /promos/{id} [delete]
// @Security     BearerAuth
func (h *PromoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	if err := h.promoUC.DeletePromo(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Promo deleted", nil)
}
