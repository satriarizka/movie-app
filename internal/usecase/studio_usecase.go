package usecase

import (
	"errors"
	"math"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/domain"
	"movie-app/internal/repository"
	"movie-app/pkg/utils"

	"github.com/google/uuid"
)

type StudioUseCase interface {
	Create(req request.CreateStudioRequest) (*domain.Studio, error)
	Update(id uuid.UUID, req request.UpdateStudioRequest) (*domain.Studio, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*domain.Studio, error)
	GetAll(page int, limit int) ([]domain.Studio, *utils.PaginationMeta, error)
}

type studioUseCase struct {
	studioRepo repository.StudioRepository
}

func NewStudioUseCase(studioRepo repository.StudioRepository) StudioUseCase {
	return &studioUseCase{studioRepo}
}

func (uc *studioUseCase) Create(req request.CreateStudioRequest) (*domain.Studio, error) {
	studio := &domain.Studio{
		Name:     req.Name,
		Capacity: req.Capacity,
	}

	// --- LOGIC AUTO GENERATE SEATS ---
	seatsPerRow := 10 // Konfigurasi: 1 baris isi 10 kursi
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"}

	var seats []domain.Seat

	// Algoritma pembagian kursi
	currentRowIdx := 0
	currentSeatNum := 1

	for i := 0; i < req.Capacity; i++ {
		// Buat kursi baru
		seats = append(seats, domain.Seat{
			RowCode:    rows[currentRowIdx],
			SeatNumber: currentSeatNum,
			// StudioID akan diisi otomatis oleh GORM saat create studio
		})

		currentSeatNum++

		// Jika sudah mencapai batas per baris, pindah ke baris berikutnya (A -> B)
		if currentSeatNum > seatsPerRow {
			currentSeatNum = 1
			currentRowIdx++
			// Safety check agar tidak index out of bound jika kapasitas kegedean
			if currentRowIdx >= len(rows) {
				currentRowIdx = 0 // Atau handle error jika perlu
			}
		}
	}

	// Masukkan seats ke dalam object studio
	studio.Seats = seats

	// --- END LOGIC ---

	if err := uc.studioRepo.Create(studio); err != nil {
		return nil, err
	}
	return studio, nil
}

func (uc *studioUseCase) Update(id uuid.UUID, req request.UpdateStudioRequest) (*domain.Studio, error) {
	studio, err := uc.studioRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("studio not found")
	}

	if req.Name != "" {
		studio.Name = req.Name
	}
	if req.Capacity > 0 {
		studio.Capacity = req.Capacity
	}

	if err := uc.studioRepo.Update(studio); err != nil {
		return nil, err
	}
	return studio, nil
}

func (uc *studioUseCase) Delete(id uuid.UUID) error {
	// Cek dulu apakah ada
	_, err := uc.studioRepo.FindByID(id)
	if err != nil {
		return errors.New("studio not found")
	}
	return uc.studioRepo.Delete(id)
}

func (uc *studioUseCase) GetByID(id uuid.UUID) (*domain.Studio, error) {
	studio, err := uc.studioRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("studio not found")
	}
	return studio, nil
}

func (uc *studioUseCase) GetAll(page int, limit int) ([]domain.Studio, *utils.PaginationMeta, error) {
	studios, total, err := uc.studioRepo.FindAll(page, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	meta := &utils.PaginationMeta{
		CurrentPage: page,
		TotalPage:   totalPages,
		TotalItems:  total,
		Limit:       limit,
	}

	return studios, meta, nil
}
