package worker

import (
	"movie-app/internal/usecase"
	"movie-app/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type Scheduler struct {
	transUC usecase.TransactionUseCase
	ticker  *time.Ticker
	quit    chan bool
}

func NewScheduler(transUC usecase.TransactionUseCase) *Scheduler {
	return &Scheduler{
		transUC: transUC,
		// Ticker default 1 menit, bisa diambil dari config sebenarnya
		ticker: time.NewTicker(1 * time.Minute),
		quit:   make(chan bool),
	}
}

// Start menjalankan scheduler di background
func (s *Scheduler) Start() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.runJobs()
			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop memberhentikan scheduler (Graceful Shutdown)
func (s *Scheduler) Stop() {
	s.quit <- true
}

// runJobs berisi daftar pekerjaan yang harus dilakukan
func (s *Scheduler) runJobs() {
	// Job 1: Auto Cancel
	if err := s.transUC.AutoCancelExpiredTransactions(); err != nil {
		logger.Log.Error("Scheduler: AutoCancel error", zap.Error(err))
	}

	// Job 2: Auto Reminder
	if err := s.transUC.SendUpcomingScheduleReminders(); err != nil {
		logger.Log.Error("Scheduler: Reminder error", zap.Error(err))
	}
}
