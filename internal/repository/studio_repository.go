package repository

import (
	"movie-app/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudioRepository interface {
	Create(studio *domain.Studio) error
	Update(studio *domain.Studio) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Studio, error)
	FindAll(page int, limit int) ([]domain.Studio, int64, error)
	GetSeatsByStudioID(studioID uuid.UUID) ([]domain.Seat, error)
}

type studioRepository struct {
	db *gorm.DB
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &studioRepository{db}
}

func (r *studioRepository) Create(studio *domain.Studio) error {
	return r.db.Create(studio).Error
}

func (r *studioRepository) Update(studio *domain.Studio) error {
	return r.db.Save(studio).Error
}

func (r *studioRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Studio{}, id).Error
}

func (r *studioRepository) FindByID(id uuid.UUID) (*domain.Studio, error) {
	var studio domain.Studio
	err := r.db.First(&studio, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &studio, nil
}

func (r *studioRepository) FindAll(page int, limit int) ([]domain.Studio, int64, error) {
	var studios []domain.Studio
	var total int64

	if err := r.db.Model(&domain.Studio{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Limit(limit).Offset(offset).Find(&studios).Error

	return studios, total, err
}

func (r *studioRepository) GetSeatsByStudioID(studioID uuid.UUID) ([]domain.Seat, error) {
	var seats []domain.Seat
	err := r.db.Where("studio_id = ?", studioID).Order("row_code, seat_number").Find(&seats).Error
	return seats, err
}
