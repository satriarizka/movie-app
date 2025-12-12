package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"movie-app/pkg/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudioHandler struct {
	studioUC usecase.StudioUseCase
	val      *validator.CustomValidator
}

func NewStudioHandler(studioUC usecase.StudioUseCase, val *validator.CustomValidator) *StudioHandler {
	return &StudioHandler{studioUC, val}
}

// Create godoc
// @Summary      Create new studio
// @Description  Create a new studio (Admin only)
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        request body request.CreateStudioRequest true "Studio Data"
// @Success      201  {object}  utils.APIResponse{data=response.StudioResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /studios [post]
// @Security     BearerAuth
func (h *StudioHandler) Create(c *gin.Context) {
	var req request.CreateStudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	studio, err := h.studioUC.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := response.StudioResponse{ID: studio.ID, Name: studio.Name, Capacity: studio.Capacity}
	utils.SuccessResponse(c, http.StatusCreated, "Studio created", res)
}

// GetAll godoc
// @Summary      Get all studios
// @Description  Get list of studios with pagination
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        page   query    int     false  "Page number" default(1)
// @Param        limit  query    int     false  "Limit per page" default(10)
// @Success      200    {object} utils.APIResponse{data=[]response.StudioResponse}
// @Failure      500    {object} utils.APIResponse
// @Router       /studios [get]
// @Security     BearerAuth
func (h *StudioHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	studios, meta, err := h.studioUC.GetAll(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Mapping response list
	var res []response.StudioResponse
	for _, s := range studios {
		res = append(res, response.StudioResponse{ID: s.ID, Name: s.Name, Capacity: s.Capacity})
	}

	// Kita butuh wrapper khusus untuk list dengan pagination
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "List of studios",
		"data":    res,
		"meta":    meta,
	})
}

// GetByID godoc
// @Summary      Get studio by ID
// @Description  Get details of a specific studio
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Studio UUID"
// @Success      200  {object}  utils.APIResponse{data=response.StudioResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /studios/{id} [get]
// @Security     BearerAuth
func (h *StudioHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	studio, err := h.studioUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	res := response.StudioResponse{ID: studio.ID, Name: studio.Name, Capacity: studio.Capacity}
	utils.SuccessResponse(c, http.StatusOK, "Studio found", res)
}

// Update godoc
// @Summary      Update studio
// @Description  Update studio details (Admin only)
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        id       path    string  true  "Studio UUID"
// @Param        request  body    request.UpdateStudioRequest true "Update Data"
// @Success      200      {object} utils.APIResponse{data=response.StudioResponse}
// @Failure      400      {object} utils.APIResponse
// @Router       /studios/{id} [put]
// @Security     BearerAuth
func (h *StudioHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	var req request.UpdateStudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	studio, err := h.studioUC.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := response.StudioResponse{ID: studio.ID, Name: studio.Name, Capacity: studio.Capacity}
	utils.SuccessResponse(c, http.StatusOK, "Studio updated", res)
}

// Delete godoc
// @Summary      Delete studio
// @Description  Delete a studio (Admin only)
// @Tags         Studios
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Studio UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Router       /studios/{id} [delete]
// @Security     BearerAuth
func (h *StudioHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	if err := h.studioUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Studio deleted", nil)
}
