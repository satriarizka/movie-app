package repository

import (
	"movie-app/internal/domain"
	"movie-app/internal/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	// GetBookedSeats mengambil daftar kursi yang SUDAH laku untuk jadwal tertentu
	GetBookedSeats(scheduleID uuid.UUID) ([]domain.Ticket, error)
	GetByUserID(userID uuid.UUID) ([]domain.Transaction, error) // History
	// CreateBooking melakukan insert Transaction & Tickets dalam 1 db transaction
	CreateBooking(tx *domain.Transaction) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) GetBookedSeats(scheduleID uuid.UUID) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	// Ganti domain.TransactionCancel menjadi enums.TransactionCancelled
	err := r.db.Joins("JOIN transactions ON transactions.id = tickets.transaction_id").
		Where("tickets.schedule_id = ? AND transactions.status != ?", scheduleID, enums.TransactionCancel).
		Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetByUserID(userID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Preload("Tickets.Seat").Preload("Tickets.Schedule.Movie").Preload("Tickets.Schedule.Studio").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *ticketRepository) CreateBooking(transaction *domain.Transaction) error {
	// GORM Transaction: Atomic Operation
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Header Transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 2. Create Detail Tickets (otomatis karena relasi HasMany)
		// GORM cukup pintar, jika struct transaction punya field Tickets terisi,
		// dia akan insert ke tabel tickets juga.
		return nil
	})
}
