package usecase

import (
	"errors"
	"fmt"
	"movie-app/internal/delivery/http/dto/request"
	"movie-app/internal/domain"
	"movie-app/internal/enums"
	"movie-app/internal/repository"
	"movie-app/pkg/mailer"
	"time"

	"github.com/google/uuid"
)

type TransactionUseCase interface {
	PayTransaction(userID uuid.UUID, transactionID uuid.UUID, req request.PayTransactionRequest) error
	CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error
	AutoCancelExpiredTransactions() error
	GetUserTransactions(userID uuid.UUID) ([]domain.Transaction, error)
	SendUpcomingScheduleReminders() error
}

type transactionUseCase struct {
	transRepo repository.TransactionRepository
	mailer    *mailer.Mailer
}

func NewTransactionUseCase(transRepo repository.TransactionRepository, mailer *mailer.Mailer) TransactionUseCase {
	return &transactionUseCase{transRepo, mailer}
}

func (uc *transactionUseCase) PayTransaction(userID uuid.UUID, transactionID uuid.UUID, req request.PayTransactionRequest) error {
	// 1. Cari Transaksi
	transaction, err := uc.transRepo.FindByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// 2. Validasi Kepemilikan (Security Check)
	if transaction.UserID != userID {
		return errors.New("unauthorized access to this transaction")
	}

	// 3. Validasi Status (Hanya 'pending' yang boleh dibayar)
	if transaction.Status != enums.TransactionPending {
		return errors.New("transaction is not pending (already paid or cancelled)")
	}

	if err := uc.transRepo.UpdateStatus(transactionID, enums.TransactionPaid, req.PaymentMethod); err != nil {
		return err
	}

	// --- LOGIC EMAIL NOTIFIKASI ---
	go func() {
		trx, err := uc.transRepo.FindByID(transactionID)
		if err != nil {
			fmt.Printf("ERROR Email: Transaction not found: %v\n", err)
			return
		}

		subject := "Booking Confirmed!"
		body := fmt.Sprintf(`
            <h1>Payment Successful</h1>
            <p>Hi %s, terima kasih sudah memesan tiket.</p>
            <p>Film: <b>%s</b></p>
            <p>Total: Rp %.2f</p>
        `, trx.User.Name, trx.Tickets[0].Schedule.Movie.Title, trx.FinalAmount)

		// Cek Error Send
		if err := uc.mailer.Send(trx.User.Email, subject, body); err != nil {
			fmt.Printf("ERROR Email: Failed to send email to %s: %v\n", trx.User.Email, err)
		} else {
			fmt.Printf("SUCCESS: Email sent to %s\n", trx.User.Email)
		}
	}()

	return nil
}

func (uc *transactionUseCase) CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error {
	// A. Cari Transaksi
	transaction, err := uc.transRepo.FindByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// B. Validasi Kepemilikan (User A tidak boleh cancel punya User B)
	if transaction.UserID != userID {
		return errors.New("unauthorized access to this transaction")
	}

	// C. Validasi Status (Hanya 'pending' yang boleh dicancel)
	// Kalau sudah 'paid', biasanya harus lewat proses Refund (beda fitur)
	if transaction.Status != enums.TransactionPending {
		return errors.New("cannot cancel transaction (already paid or cancelled)")
	}

	// D. Update Status jadi CANCELLED
	// Payment method dikosongkan atau biarkan string kosong
	return uc.transRepo.UpdateStatus(transactionID, enums.TransactionCancel, "")
}

func (uc *transactionUseCase) GetUserTransactions(userID uuid.UUID) ([]domain.Transaction, error) {
	return uc.transRepo.GetByUserID(userID)
}

func (uc *transactionUseCase) AutoCancelExpiredTransactions() error {
	// 1. Tentukan batas waktu (Misal: 15 Menit yang lalu)
	expiryTime := time.Now().Add(-15 * time.Minute)

	// 2. Cari transaksi yang bandel (belum bayar lewat dari 15 menit)
	expiredTransactions, err := uc.transRepo.GetExpiredPendingTransactions(expiryTime)
	if err != nil {
		return err
	}

	// 3. Loop dan Cancel satu per satu
	for _, tx := range expiredTransactions {
		// Update status ke Cancelled
		// Kita abaikan error per item agar satu gagal tidak menghentikan yang lain
		_ = uc.transRepo.UpdateStatus(tx.ID, enums.TransactionCancel, "")

		// (Optional) Log ke terminal
		// fmt.Printf("Auto-cancelling transaction: %s\n", tx.ID)
	}

	return nil
}

func (uc *transactionUseCase) SendUpcomingScheduleReminders() error {
	// Range waktu: film yang mulai 1 jam dari sekarang s/d 2 jam dari sekarang
	now := time.Now()
	startWindow := now.Add(1 * time.Hour)
	endWindow := now.Add(2 * time.Hour)

	transactions, err := uc.transRepo.GetUpcomingPaidTransactions(startWindow, endWindow)
	if err != nil {
		return err
	}

	for _, trx := range transactions {
		// Kirim Email
		subject := "Reminder: Film Anda Segera Mulai!"
		body := fmt.Sprintf(`
            <h1>Siap-siap nonton!</h1>
            <p>Hi %s, film <b>%s</b> akan dimulai dalam 1 jam lagi.</p>
            <p>Segera datang ke bioskop ya!</p>
        `, trx.User.Name, trx.Tickets[0].Schedule.Movie.Title)

		if err := uc.mailer.Send(trx.User.Email, subject, body); err == nil {
			// Jika sukses kirim, tandai di DB agar tidak kirim lagi
			_ = uc.transRepo.MarkReminderSent(trx.ID)
		}
	}
	return nil
}
