package usecase

import (
	"movie-app/internal/domain"
	"movie-app/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PromoUseCase interface {
	CreatePromo(code string, discountType string, value float64, validUntil time.Time) (*domain.Promo, error)
	GetAllPromos() ([]domain.Promo, error)
	UpdatePromo(id uuid.UUID, code string, discountType string, value float64, validUntil time.Time) (*domain.Promo, error)
	DeletePromo(id uuid.UUID) error
}

type promoUseCase struct {
	promoRepo repository.PromoRepository
}

func NewPromoUseCase(promoRepo repository.PromoRepository) PromoUseCase {
	return &promoUseCase{promoRepo}
}

func (uc *promoUseCase) CreatePromo(code string, discountType string, value float64, validUntil time.Time) (*domain.Promo, error) {
	promo := &domain.Promo{
		Code:          code,
		DiscountType:  discountType,
		DiscountValue: value,
		ValidUntil:    validUntil,
	}
	if err := uc.promoRepo.Create(promo); err != nil {
		return nil, err
	}
	return promo, nil
}

func (uc *promoUseCase) GetAllPromos() ([]domain.Promo, error) {
	return uc.promoRepo.FindAll()
}

func (uc *promoUseCase) UpdatePromo(id uuid.UUID, code string, discountType string, value float64, validUntil time.Time) (*domain.Promo, error) {
	// Cek dulu datanya ada atau tidak
	promo, err := uc.promoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update field jika tidak kosong (Partial Update logic sederhana)
	if code != "" {
		promo.Code = code
	}
	if discountType != "" {
		promo.DiscountType = discountType
	}
	if value > 0 {
		promo.DiscountValue = value
	}
	if !validUntil.IsZero() {
		promo.ValidUntil = validUntil
	}

	if err := uc.promoRepo.Update(promo); err != nil {
		return nil, err
	}
	return promo, nil
}

func (uc *promoUseCase) DeletePromo(id uuid.UUID) error {
	// Cek existensi
	if _, err := uc.promoRepo.FindByID(id); err != nil {
		return err
	}
	return uc.promoRepo.Delete(id)
}
