package usecase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/repository"
)

type ReportUseCase interface {
	GetTopMovies(limit int) ([]response.TopMovieResponse, error)
	GetRevenueReport(mode string) ([]response.DailyRevenueResponse, error)
	GenerateRevenueCSV(mode string) ([]byte, error)
}

type reportUseCase struct {
	reportRepo repository.ReportRepository
}

func NewReportUseCase(reportRepo repository.ReportRepository) ReportUseCase {
	return &reportUseCase{reportRepo}
}

func (uc *reportUseCase) GetTopMovies(limit int) ([]response.TopMovieResponse, error) {
	if limit <= 0 {
		limit = 5 // Default top 5
	}
	return uc.reportRepo.GetTopMovies(limit)
}

func (uc *reportUseCase) GetRevenueReport(mode string) ([]response.DailyRevenueResponse, error) {
	// Mode: 'day' atau 'month'
	if mode != "month" {
		mode = "day"
	}
	return uc.reportRepo.GetRevenueReport(mode)
}

// Implementasi Generate CSV
func (uc *reportUseCase) GenerateRevenueCSV(mode string) ([]byte, error) {
	// 1. Ambil Data dari Repo
	data, err := uc.GetRevenueReport(mode)
	if err != nil {
		return nil, err
	}

	// 2. Buat Buffer untuk menampung CSV
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	// 3. Tulis Header CSV
	if err := w.Write([]string{"Date/Period", "Transaction Count", "Total Revenue"}); err != nil {
		return nil, err
	}

	// 4. Tulis Data Baris per Baris
	for _, item := range data {
		record := []string{
			item.Date,
			fmt.Sprintf("%d", item.Count),
			fmt.Sprintf("%.2f", item.TotalAmount),
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush() // Pastikan semua data tertulis ke buffer

	if err := w.Error(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
