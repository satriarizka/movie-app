package usecase

import (
	"errors"
	"math"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/domain"
	"movie-app/internal/repository"
	"movie-app/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type ScheduleUseCase interface {
	Create(req request.CreateScheduleRequest) (*domain.Schedule, error)
	GetByID(id uuid.UUID) (*domain.Schedule, error)
	GetAll(page int, limit int) ([]domain.Schedule, *utils.PaginationMeta, error)
	Update(id uuid.UUID, req request.UpdateScheduleRequest) (*domain.Schedule, error)
	Delete(id uuid.UUID) error
	// Update bisa Anda tambahkan sendiri nanti sbg latihan
}

type scheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	movieRepo    repository.MovieRepository
	studioRepo   repository.StudioRepository
}

// Butuh repo lain untuk validasi ID
func NewScheduleUseCase(
	sRepo repository.ScheduleRepository,
	mRepo repository.MovieRepository,
	stRepo repository.StudioRepository,
) ScheduleUseCase {
	return &scheduleUseCase{sRepo, mRepo, stRepo}
}

func (uc *scheduleUseCase) Create(req request.CreateScheduleRequest) (*domain.Schedule, error) {
	// 1. Parsing UUID
	studioID, _ := uuid.Parse(req.StudioID)
	movieID, _ := uuid.Parse(req.MovieID)

	// 2. Validasi Exist (Apakah Studio & Movie ada?)
	_, err := uc.studioRepo.FindByID(studioID)
	if err != nil {
		return nil, errors.New("studio not found")
	}
	_, err = uc.movieRepo.FindByID(movieID)
	if err != nil {
		return nil, errors.New("movie not found")
	}

	// 3. Validasi Konflik Jadwal
	isOverlap, err := uc.scheduleRepo.CheckOverlap(studioID, req.StartTime, req.EndTime, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if isOverlap {
		return nil, errors.New("schedule overlaps with existing showtime")
	}

	// 4. Dynamic Pricing (Prime Time Logic)
	finalPrice := req.Price
	hour := req.StartTime.Hour()
	weekday := req.StartTime.Weekday()

	// pembuatan otomatis menambahkan 10000, jika terjadi di saturday dan sunday
	isPrimeTime := hour >= 17 || weekday == time.Saturday || weekday == time.Sunday
	if isPrimeTime {
		finalPrice += 10000
	}

	// 5. Simpan ke Database
	schedule := &domain.Schedule{
		StudioID:  studioID,
		MovieID:   movieID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Price:     finalPrice,
	}

	if err := uc.scheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	// 6. Ambil data lengkap (Reload dari DB agar Preload Studio & Movie muncul)
	// --- Baris yang error dihapus, diganti logic fetch ulang yang proper ---
	createdSchedule, err := uc.scheduleRepo.FindByID(schedule.ID)
	if err != nil {
		return nil, err
	}

	return createdSchedule, nil
}

func (uc *scheduleUseCase) Update(id uuid.UUID, req request.UpdateScheduleRequest) (*domain.Schedule, error) {
	// 1. Cek apakah data jadwal ada?
	schedule, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	// 2. Parsing UUID & Validasi Input jika ada perubahan
	// Kalau user kirim StudioID baru, validasi
	if req.StudioID != "" {
		sID, _ := uuid.Parse(req.StudioID)
		if _, err := uc.studioRepo.FindByID(sID); err != nil {
			return nil, errors.New("studio not found")
		}
		schedule.StudioID = sID
	}

	// Kalau user kirim MovieID baru, validasi
	if req.MovieID != "" {
		mID, _ := uuid.Parse(req.MovieID)
		if _, err := uc.movieRepo.FindByID(mID); err != nil {
			return nil, errors.New("movie not found")
		}
		schedule.MovieID = mID
	}

	// 3. Update Waktu & Cek Konflik
	// Kita perlu tau waktu "baru" untuk pengecekan konflik
	newStart := schedule.StartTime
	newEnd := schedule.EndTime

	if !req.StartTime.IsZero() {
		newStart = req.StartTime
	}
	if !req.EndTime.IsZero() {
		newEnd = req.EndTime
	}

	// Cek konflik dengan jadwal lain (kecuali dirinya sendiri 'id')
	isOverlap, err := uc.scheduleRepo.CheckOverlap(schedule.StudioID, newStart, newEnd, id)
	if err != nil {
		return nil, err
	}
	if isOverlap {
		return nil, errors.New("schedule overlaps with existing showtime")
	}

	// Apply perubahan waktu
	schedule.StartTime = newStart
	schedule.EndTime = newEnd

	// 4. Update Harga (Optional)
	if req.Price > 0 {
		schedule.Price = req.Price
	}

	// 5. Simpan Perubahan
	if err := uc.scheduleRepo.Update(schedule); err != nil {
		return nil, err
	}

	// 6. Return data terbaru
	return schedule, nil
}

func (uc *scheduleUseCase) GetByID(id uuid.UUID) (*domain.Schedule, error) {
	return uc.scheduleRepo.FindByID(id)
}

func (uc *scheduleUseCase) GetAll(page int, limit int) ([]domain.Schedule, *utils.PaginationMeta, error) {
	schedules, total, err := uc.scheduleRepo.FindAll(page, limit)
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
	return schedules, meta, nil
}

func (uc *scheduleUseCase) Delete(id uuid.UUID) error {
	_, err := uc.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("schedule not found")
	}
	return uc.scheduleRepo.Delete(id)
}
