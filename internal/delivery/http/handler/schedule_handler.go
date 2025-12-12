package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/domain"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"movie-app/pkg/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduleHandler struct {
	scheduleUC usecase.ScheduleUseCase
	val        *validator.CustomValidator
}

func NewScheduleHandler(scheduleUC usecase.ScheduleUseCase, val *validator.CustomValidator) *ScheduleHandler {
	return &ScheduleHandler{scheduleUC, val}
}

// Create godoc
// @Summary      Create new schedule
// @Description  Add a new movie schedule (Admin only)
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        request body request.CreateScheduleRequest true "Schedule Data"
// @Success      201  {object}  utils.APIResponse{data=response.ScheduleResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse "Conflict / Schedule Overlap"
// @Router       /schedules [post]
// @Security     BearerAuth
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req request.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body format", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	schedule, err := h.scheduleUC.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	res := h.mapResponse(schedule)
	utils.SuccessResponse(c, http.StatusCreated, "Schedule created", res)
}

// Update godoc
// @Summary      Update schedule
// @Description  Update existing schedule details (Admin only)
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        id       path    string  true  "Schedule UUID"
// @Param        request  body    request.UpdateScheduleRequest true "Update Data"
// @Success      200      {object} utils.APIResponse{data=response.ScheduleResponse}
// @Failure      400      {object} utils.APIResponse
// @Router       /schedules/{id} [put]
// @Security     BearerAuth
func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	var req request.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body format", err.Error())
		return
	}

	// Validasi input optional
	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	schedule, err := h.scheduleUC.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := h.mapResponse(schedule)
	utils.SuccessResponse(c, http.StatusOK, "Schedule updated", res)
}

// GetAll godoc
// @Summary      Get all schedules
// @Description  Get list of schedules with pagination (Public)
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        page   query    int     false  "Page number" default(1)
// @Param        limit  query    int     false  "Limit per page" default(10)
// @Success      200    {object} utils.APIResponse{data=[]response.ScheduleResponse}
// @Router       /schedules [get]
func (h *ScheduleHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	schedules, meta, err := h.scheduleUC.GetAll(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var res []response.ScheduleResponse
	for _, s := range schedules {
		res = append(res, h.mapResponse(&s))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "List of schedules",
		"data":    res,
		"meta":    meta,
	})
}

// Delete godoc
// @Summary      Delete schedule
// @Description  Delete a schedule (Admin only)
// @Tags         Schedules
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Router       /schedules/{id} [delete]
// @Security     BearerAuth
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	if err := h.scheduleUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Schedule deleted", nil)
}

func (h *ScheduleHandler) mapResponse(s *domain.Schedule) response.ScheduleResponse {
	return response.ScheduleResponse{
		ID:        s.ID,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Price:     s.Price,
		Studio: response.StudioResponse{
			ID:       s.Studio.ID,
			Name:     s.Studio.Name,
			Capacity: s.Studio.Capacity,
		},
		Movie: response.MovieResponse{
			ID:          s.Movie.ID,
			Title:       s.Movie.Title,
			Description: s.Movie.Description,
			Duration:    s.Movie.Duration,
			Genre:       s.Movie.Genre,
			PosterURL:   s.Movie.PosterURL,
		},
	}
}
