package repository

import (
	"movie-app/internal/domain"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo *domain.Promo) error
	FindByID(id uuid.UUID) (*domain.Promo, error)
	Update(promo *domain.Promo) error
	Delete(id uuid.UUID) error
	FindByCode(code string) (*domain.Promo, error)
	FindAll() ([]domain.Promo, error)
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db}
}

func (r *promoRepository) Create(promo *domain.Promo) error {
	return r.db.Create(promo).Error
}

func (r *promoRepository) FindByCode(code string) (*domain.Promo, error) {
	var promo domain.Promo
	// Cek kode dan pastikan belum expired
	err := r.db.Where("code = ? AND valid_until > NOW()", code).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) FindAll() ([]domain.Promo, error) {
	var promos []domain.Promo
	err := r.db.Find(&promos).Error
	return promos, err
}

func (r *promoRepository) FindByID(id uuid.UUID) (*domain.Promo, error) {
	var promo domain.Promo
	err := r.db.First(&promo, "id = ?", id).Error
	return &promo, err
}

func (r *promoRepository) Update(promo *domain.Promo) error {
	return r.db.Save(promo).Error
}

func (r *promoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Promo{}, id).Error
}
