package handler

import (
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/usecase"
	"movie-app/pkg/utils"
	"movie-app/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	val         *validator.CustomValidator
}

func NewAuthHandler(authUC usecase.AuthUseCase, val *validator.CustomValidator) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUC,
		val:         val,
	}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.RegisterRequest true "User Data"
// @Success      201  {object}  utils.APIResponse{data=response.UserResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	user, err := h.authUseCase.Register(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	// Mapping ke Response agar password tidak ikut
	userResponse := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", userResponse)
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	user, err := h.authUseCase.RegisterAdmin(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	userResponse := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin registered successfully", userResponse)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 1. Ambil user_id dari context (diset oleh Middleware Auth)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// 2. Konversi ke UUID
	// JWT Claims biasanya menyimpan string, jadi kita parse dulu
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID format", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user UUID", nil)
		return
	}

	// 3. Panggil UseCase
	user, err := h.authUseCase.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	// 4. Return Response
	userResponse := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved", userResponse)
}

// Login godoc
// @Summary      Login User
// @Description  Login with email and password to get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "Login Credentials"
// @Success      200  {object}  utils.APIResponse{data=response.AuthResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.val.Validate(req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", validator.FormatError(err))
		return
	}

	authRes, err := h.authUseCase.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", authRes)
}
