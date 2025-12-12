package repository

import (
	"fmt"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/enums"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetTopMovies(limit int) ([]response.TopMovieResponse, error)
	GetRevenueReport(groupBy string) ([]response.DailyRevenueResponse, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetTopMovies(limit int) ([]response.TopMovieResponse, error) {
	var results []response.TopMovieResponse

	// Query Join 4 Tabel: Transactions -> Tickets -> Schedules -> Movies
	// Hitung jumlah tiket per film
	err := r.db.Table("tickets").
		Select("movies.id as movie_id, movies.title, COUNT(tickets.id) as total_sold, SUM(schedules.price) as total_sales").
		Joins("JOIN transactions ON transactions.id = tickets.transaction_id").
		Joins("JOIN schedules ON schedules.id = tickets.schedule_id").
		Joins("JOIN movies ON movies.id = schedules.movie_id").
		Where("transactions.status = ?", enums.TransactionPaid).
		Group("movies.id, movies.title").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

func (r *reportRepository) GetRevenueReport(groupBy string) ([]response.DailyRevenueResponse, error) {
	var results []response.DailyRevenueResponse
	var dateFormat string

	// Tentukan format tanggal PostgreSQL berdasarkan grouping
	if groupBy == "month" {
		dateFormat = "YYYY-MM" // Format Tahun-Bulan (2025-12)
	} else {
		dateFormat = "YYYY-MM-DD" // Default Harian (2025-12-11)
	}

	// Query Dynamic
	querySelect := fmt.Sprintf("TO_CHAR(created_at, '%s') as date, SUM(final_amount) as total_amount, COUNT(id) as count", dateFormat)

	err := r.db.Table("transactions").
		Select(querySelect).
		Where("status = ?", enums.TransactionPaid). // Hanya yang sudah lunas
		Group("date").
		Order("date DESC").
		Scan(&results).Error

	return results, err
}
