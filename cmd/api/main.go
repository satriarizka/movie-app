package main

import (
	"context"
	"fmt"

	// IMPORT DOCS (Wajib pakai underscore)
	_ "movie-app/docs"

	"movie-app/internal/config"
	"movie-app/internal/delivery/http/handler"
	"movie-app/internal/delivery/http/middleware"
	"movie-app/internal/delivery/http/route"
	"movie-app/internal/delivery/worker" // Import Worker Package Baru
	"movie-app/internal/repository"
	"movie-app/internal/usecase"
	"movie-app/pkg/database"
	"movie-app/pkg/logger"
	"movie-app/pkg/mailer"
	"movie-app/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Bioskop API
// @version         1.0
// @description     API Server for Movie Booking Application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Config & Logger
	cfg := config.LoadConfig()
	logger.InitLogger()

	// 2. Database
	db := database.ConnectDB(cfg)

	// 3. Validators
	val := validator.NewValidator()

	// Setup Mailer
	mailService := mailer.NewMailer(cfg)

	// 4. Layers Initialization
	// Repository
	userRepo := repository.NewUserRepository(db)
	studioRepo := repository.NewStudioRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	transRepo := repository.NewTransactionRepository(db)
	reportRepo := repository.NewReportRepository(db)
	promoRepo := repository.NewPromoRepository(db)

	// UseCase
	authUC := usecase.NewAuthUseCase(userRepo, cfg)
	studioUC := usecase.NewStudioUseCase(studioRepo)
	movieUC := usecase.NewMovieUseCase(movieRepo)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, movieRepo, studioRepo)
	ticketUC := usecase.NewTicketUseCase(ticketRepo, scheduleRepo, studioRepo, promoRepo)
	transUC := usecase.NewTransactionUseCase(transRepo, mailService)
	reportUC := usecase.NewReportUseCase(reportRepo)
	promoUC := usecase.NewPromoUseCase(promoRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authUC, val)
	studioHandler := handler.NewStudioHandler(studioUC, val)
	movieHandler := handler.NewMovieHandler(movieUC, val)
	scheduleHandler := handler.NewScheduleHandler(scheduleUC, val)
	ticketHandler := handler.NewTicketHandler(ticketUC, val)
	transHandler := handler.NewTransactionHandler(transUC, val)
	reportHandler := handler.NewReportHandler(reportUC)
	promoHandler := handler.NewPromoHandler(promoUC)

	// 5. Worker / Scheduler Initialization
	// Logic background job dipisah ke package worker agar main.go bersih
	bgWorker := worker.NewScheduler(transUC)
	logger.Log.Info("Starting background scheduler...")
	bgWorker.Start()

	// 6. Router Setup
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogger())
	r.Use(gin.Recovery())

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "UP"})
		})
		route.SetupRoutes(api, authHandler, studioHandler, movieHandler, scheduleHandler, ticketHandler, transHandler, reportHandler, promoHandler, cfg)
	}

	// 7. Server Setup
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Run Server in Goroutine
	go func() {
		logger.Log.Info(fmt.Sprintf("Server running on port %s", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Listen error", zap.Error(err))
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	// Stop Scheduler dulu
	bgWorker.Stop()

	// Shutdown HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Log.Info("Server exiting")
}
