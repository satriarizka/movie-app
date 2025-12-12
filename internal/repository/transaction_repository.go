package repository

import (
	"movie-app/internal/domain"
	"movie-app/internal/enums"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(id uuid.UUID) (*domain.Transaction, error)
	UpdateStatus(id uuid.UUID, status enums.TransactionStatus, paymentMethod string) error
	GetByUserID(userID uuid.UUID) ([]domain.Transaction, error)
	GetExpiredPendingTransactions(threshold time.Time) ([]domain.Transaction, error)
	GetUpcomingPaidTransactions(startTime, endTime time.Time) ([]domain.Transaction, error)
	MarkReminderSent(id uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*domain.Transaction, error) {
	var transaction domain.Transaction

	// FIX: Tambahkan Preload("User") agar data user (email/nama) ikut terambil
	// Preload Tiket & Kursi tetap ada agar data tiket tidak hilang
	err := r.db.
		Preload("User").
		Preload("Tickets.Seat").
		Preload("Tickets.Schedule.Movie").
		First(&transaction, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateStatus(id uuid.UUID, status enums.TransactionStatus, paymentMethod string) error {
	return r.db.Model(&domain.Transaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         status,
		"payment_method": paymentMethod,
	}).Error
}

func (r *transactionRepository) GetByUserID(userID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.Preload("Tickets.Seat").
		Preload("Tickets.Schedule.Movie").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetExpiredPendingTransactions(threshold time.Time) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	// Query: Status = Pending AND CreatedAt < (Waktu Sekarang - 15 menit)
	err := r.db.Where("status = ? AND created_at < ?", enums.TransactionPending, threshold).
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetUpcomingPaidTransactions(startTime, endTime time.Time) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	// Join ke Schedule -> Preload User (untuk email) & Movie (untuk judul film)
	err := r.db.
		Preload("User").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio").
		Joins("JOIN tickets ON tickets.transaction_id = transactions.id").
		Joins("JOIN schedules ON schedules.id = tickets.schedule_id").
		Where("transactions.status = ? AND transactions.reminder_sent = ?", enums.TransactionPaid, false).
		Where("schedules.start_time BETWEEN ? AND ?", startTime, endTime).
		Distinct("transactions.id"). // Mencegah duplikat karena join tickets
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) MarkReminderSent(id uuid.UUID) error {
	// Tandai bahwa reminder sudah dikirim agar user tidak dispam email
	return r.db.Model(&domain.Transaction{}).Where("id = ?", id).Update("reminder_sent", true).Error
}
