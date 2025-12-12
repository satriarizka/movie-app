package usecase

import (
	"errors"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/delivery/http/dto/response"
	"movie-app/internal/domain"
	"movie-app/internal/enums"
	"movie-app/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TicketUseCase interface {
	GetAvailableSeats(scheduleID uuid.UUID) ([]response.SeatAvailabilityResponse, error)
	BookTicket(userID uuid.UUID, req request.BookTicketRequest) (*domain.Transaction, error)
	GetUserHistory(userID uuid.UUID) ([]domain.Transaction, error)
}

type ticketUseCase struct {
	ticketRepo   repository.TicketRepository
	scheduleRepo repository.ScheduleRepository
	studioRepo   repository.StudioRepository
	promoRepo    repository.PromoRepository
}

func NewTicketUseCase(
	tRepo repository.TicketRepository,
	sRepo repository.ScheduleRepository,
	stRepo repository.StudioRepository,
	pRepo repository.PromoRepository,
) TicketUseCase {
	// FIX 1: Masukkan pRepo ke struct return
	return &ticketUseCase{
		ticketRepo:   tRepo,
		scheduleRepo: sRepo,
		studioRepo:   stRepo,
		promoRepo:    pRepo,
	}
}

func (uc *ticketUseCase) GetAvailableSeats(scheduleID uuid.UUID) ([]response.SeatAvailabilityResponse, error) {
	// ... (Logika GetAvailableSeats tidak berubah, copy dari sebelumnya) ...
	schedule, err := uc.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	allSeats, err := uc.studioRepo.GetSeatsByStudioID(schedule.StudioID)
	if err != nil {
		return nil, err
	}

	bookedTickets, err := uc.ticketRepo.GetBookedSeats(scheduleID)
	if err != nil {
		return nil, err
	}

	bookedMap := make(map[uuid.UUID]bool)
	for _, t := range bookedTickets {
		bookedMap[t.SeatID] = true
	}

	var result []response.SeatAvailabilityResponse
	for _, seat := range allSeats {
		isBooked := bookedMap[seat.ID]
		result = append(result, response.SeatAvailabilityResponse{
			ID:         seat.ID,
			RowCode:    seat.RowCode,
			SeatNumber: seat.SeatNumber,
			IsBooked:   isBooked,
		})
	}

	return result, nil
}

func (uc *ticketUseCase) BookTicket(userID uuid.UUID, req request.BookTicketRequest) (*domain.Transaction, error) {
	// 1. Validasi Jadwal
	scheduleID, _ := uuid.Parse(req.ScheduleID)
	schedule, err := uc.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	// 2. Hitung Harga Dasar
	totalAmount := schedule.Price * float64(len(req.SeatIDs))

	// FIX 2: Jangan buat transaction struct dulu. Buat slice tiket dulu.
	var tickets []domain.Ticket

	// 3. Siapkan Tiket ke dalam Slice
	for _, seatIDStr := range req.SeatIDs {
		seatID, _ := uuid.Parse(seatIDStr)
		ticket := domain.Ticket{
			ScheduleID: scheduleID,
			SeatID:     seatID,
		}
		tickets = append(tickets, ticket)
	}

	// 4. Logic Promo
	var discountAmount float64 = 0
	var promoID *uuid.UUID = nil

	if req.PromoCode != "" {
		promo, err := uc.promoRepo.FindByCode(req.PromoCode)
		if err != nil {
			return nil, errors.New("promo code invalid or expired")
		}

		promoID = &promo.ID

		if promo.DiscountType == enums.DiscountTypePercentage {
			discountAmount = totalAmount * (promo.DiscountValue / 100)
		} else {
			discountAmount = promo.DiscountValue
		}

		if discountAmount > totalAmount {
			discountAmount = totalAmount
		}
	}
	finalAmount := totalAmount - discountAmount

	// 5. Build Transaction Struct (SEKALI SAJA DI SINI)
	transaction := &domain.Transaction{
		UserID:         userID,
		TotalAmount:    totalAmount,    // Harga Asli
		DiscountAmount: discountAmount, // Potongan
		FinalAmount:    finalAmount,    // Harga Akhir
		PromoID:        promoID,
		Status:         enums.TransactionPending,
		PaymentMethod:  "",
		Tickets:        tickets, // Masukkan slice tiket yang sudah dibuat
	}

	// 6. Simpan (Atomic Transaction)
	if err := uc.ticketRepo.CreateBooking(transaction); err != nil {
		return nil, errors.New("some seats are already booked")
	}

	// Set ExpiresAt untuk Response
	transaction.ExpiresAt = transaction.CreatedAt.Add(15 * time.Minute)

	return transaction, nil
}

func (uc *ticketUseCase) GetUserHistory(userID uuid.UUID) ([]domain.Transaction, error) {
	return uc.ticketRepo.GetByUserID(userID)
}
