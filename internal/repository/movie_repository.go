package repository

import (
	"movie-app/internal/domain"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(movie *domain.Movie) error
	Update(movie *domain.Movie) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Movie, error)
	// Update signature: Tambah parameter 'search'
	FindAll(page int, limit int, search string) ([]domain.Movie, int64, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db}
}

// ... (Create, Update, Delete, FindByID SAMA SEPERTI SEBELUMNYA, tidak perlu diubah) ...
func (r *movieRepository) Create(movie *domain.Movie) error {
	return r.db.Create(movie).Error
}

func (r *movieRepository) Update(movie *domain.Movie) error {
	return r.db.Save(movie).Error
}

func (r *movieRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Movie{}, id).Error
}

func (r *movieRepository) FindByID(id uuid.UUID) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.db.First(&movie, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

// === BAGIAN YANG DIUPDATE ===
func (r *movieRepository) FindAll(page int, limit int, search string) ([]domain.Movie, int64, error) {
	var movies []domain.Movie
	var total int64

	// 1. Inisialisasi query db
	query := r.db.Model(&domain.Movie{})

	// 2. Logic Search (Title OR Genre)
	if search != "" {
		searchLower := "%" + strings.ToLower(search) + "%"
		// ILIKE untuk case-insensitive di Postgres
		query = query.Where("LOWER(title) LIKE ? OR LOWER(genre) LIKE ?", searchLower, searchLower)
	}

	// 3. Hitung Total Data (setelah difilter)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Pagination & Execution
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}
