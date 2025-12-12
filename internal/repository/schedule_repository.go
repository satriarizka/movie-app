package repository

import (
	"movie-app/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(schedule *domain.Schedule) error
	Update(schedule *domain.Schedule) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Schedule, error)
	FindAll(page int, limit int) ([]domain.Schedule, int64, error)
	// CheckOverlap mengecek apakah ada jadwal lain di studio yg sama pada rentang waktu tsb
	CheckOverlap(studioID uuid.UUID, startTime, endTime time.Time, excludeID uuid.UUID) (bool, error)
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db}
}

func (r *scheduleRepository) Create(schedule *domain.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) Update(schedule *domain.Schedule) error {
	return r.db.Save(schedule).Error
}

func (r *scheduleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Schedule{}, id).Error
}

func (r *scheduleRepository) FindByID(id uuid.UUID) (*domain.Schedule, error) {
	var schedule domain.Schedule
	// Preload Studio dan Movie agar datanya lengkap saat diambil
	err := r.db.Preload("Studio").Preload("Movie").First(&schedule, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindAll(page int, limit int) ([]domain.Schedule, int64, error) {
	var schedules []domain.Schedule
	var total int64

	if err := r.db.Model(&domain.Schedule{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Studio").Preload("Movie").
		Limit(limit).Offset(offset).
		Order("start_time DESC").
		Find(&schedules).Error

	return schedules, total, err
}

func (r *scheduleRepository) CheckOverlap(studioID uuid.UUID, startTime, endTime time.Time, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Schedule{}).
		Where("studio_id = ?", studioID).
		Where("((start_time < ?) AND (end_time > ?))", endTime, startTime) // Logika Overlap

	// Jika sedang Update, jangan anggap jadwal diri sendiri sebagai bentrok
	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
