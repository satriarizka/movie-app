package usecase

import (
	"errors"
	"movie-app/internal/config"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/domain"
	"movie-app/internal/enums"
	"movie-app/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(req request.RegisterRequest) (*domain.User, error)
	RegisterAdmin(req request.RegisterRequest) (*domain.User, error)
	Login(req request.LoginRequest) (*response.AuthResponse, error)
	GetProfile(user_id uuid.UUID) (*domain.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthUseCase(userRepo repository.UserRepository, cfg *config.Config) AuthUseCase {
	return &authUseCase{userRepo, cfg}
}

func (uc *authUseCase) Register(req request.RegisterRequest) (*domain.User, error) {
	// 1. Cek apakah email sudah ada
	existingUser, err := uc.userRepo.FindByEmail(req.Email)
	if err == nil || existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Buat Object User
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     enums.RoleUser,
	}

	// 4. Simpan ke DB
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUseCase) RegisterAdmin(req request.RegisterRequest) (*domain.User, error) {
	// ⚠️ WARNING: Method ini hanya untuk Development/Seeding awal.
	// Sebaiknya dihapus atau diproteksi key khusus saat production.

	existingUser, _ := uc.userRepo.FindByEmail((req.Email))
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     enums.RoleAdmin,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *authUseCase) GetProfile(user_id uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(user_id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *authUseCase) Login(req request.LoginRequest) (*response.AuthResponse, error) {
	// 1. Cari user by Email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Cek Password (Bandingkan Hash)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	token, err := uc.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{Token: token}, nil
}

func (uc *authUseCase) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 Jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.JWTSecret))
}
