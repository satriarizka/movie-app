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

type MovieHandler struct {
	movieUC usecase.MovieUseCase
	val     *validator.CustomValidator
}

func NewMovieHandler(movieUC usecase.MovieUseCase, val *validator.CustomValidator) *MovieHandler {
	return &MovieHandler{movieUC, val}
}

// Create godoc
// @Summary      Create new movie
// @Description  Add a new movie (Admin only)
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        request body request.CreateMovieRequest true "Movie Data"
// @Success      201  {object}  utils.APIResponse{data=response.MovieResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /movies [post]
// @Security     BearerAuth
func (h *MovieHandler) Create(c *gin.Context) {
	var req request.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}
	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	movie, err := h.movieUC.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := h.mapResponse(movie)
	utils.SuccessResponse(c, http.StatusCreated, "Movie created", res)
}

// GetAll godoc
// @Summary      Get all movies
// @Description  Get list of movies (Public)
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        page   query    int     false  "Page number" default(1)
// @Param        limit  query    int     false  "Limit per page" default(10)
// @Success      200    {object} utils.APIResponse{data=[]response.MovieResponse}
// @Router       /movies [get]
func (h *MovieHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search") // Ambil query params ?search=...

	movies, meta, err := h.movieUC.GetAll(page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var res []response.MovieResponse
	for _, m := range movies {
		res = append(res, h.mapResponse(&m))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "List of movies",
		"data":    res,
		"meta":    meta,
	})
}

func (h *MovieHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	movie, err := h.movieUC.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Movie found", h.mapResponse(movie))
}

func (h *MovieHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	var req request.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid body", err.Error())
		return
	}

	movie, err := h.movieUC.Update(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Movie updated", h.mapResponse(movie))
}

func (h *MovieHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID", nil)
		return
	}

	if err := h.movieUC.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Movie deleted", nil)
}

// Helper untuk mapping agar kode tidak berulang
func (h *MovieHandler) mapResponse(m *domain.Movie) response.MovieResponse {
	return response.MovieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		Genre:       m.Genre,
		PosterURL:   m.PosterURL,
	}
}
