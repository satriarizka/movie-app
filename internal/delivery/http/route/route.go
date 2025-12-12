package route

import (
	"movie-app/internal/config"
	"movie-app/internal/delivery/http/handler"
	"movie-app/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, authHandler *handler.AuthHandler, studioHandler *handler.StudioHandler, movieHandler *handler.MovieHandler, scheduleHandler *handler.ScheduleHandler, ticketHandler *handler.TicketHandler, transactionHandler *handler.TransactionHandler, reportHandler *handler.ReportHandler, promoHandler *handler.PromoHandler, cfg *config.Config) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/register-admin", authHandler.RegisterAdmin)
		auth.POST("/login", authHandler.Login)

		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/me", authHandler.GetProfile)
		}
	}

	// studio routes
	studio := r.Group("/studios")
	// all route studio membutuhkan login admin
	studio.Use(middleware.AuthMiddleware(cfg))
	{
		// Public (User biasa bisa lihat)
		studio.GET("", studioHandler.GetAll)
		studio.GET("/:id", studioHandler.GetByID)

		// Admin only (Create, Update, Delete)
		admin := studio.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("", studioHandler.Create)
			admin.PUT("/:id", studioHandler.Update)
			admin.DELETE("/:id", studioHandler.Delete)
		}
	}

	// Movie Routes
	movies := r.Group("/movies")

	// public (bisa dilihat tanpa login)
	movies.GET("", movieHandler.GetAll)
	movies.GET("/:id", movieHandler.GetByID)

	// Admin Only (Create, Update, Delete)
	moviesAdmin := movies.Group("/")
	moviesAdmin.Use(middleware.AuthMiddleware(cfg))
	moviesAdmin.Use(middleware.AdminMiddleware())
	{
		moviesAdmin.POST("", movieHandler.Create)
		moviesAdmin.PUT("/:id", movieHandler.Update)
		moviesAdmin.DELETE("/:id", movieHandler.Delete)
	}

	// Schedule route
	schedules := r.Group("/schedules")
	// Public boleh lihat jadwal
	schedules.GET("", scheduleHandler.GetAll)

	// Admin Only
	schedulesAdmin := schedules.Group("/")
	schedulesAdmin.Use(middleware.AuthMiddleware(cfg))
	schedulesAdmin.Use(middleware.AdminMiddleware())
	{
		schedulesAdmin.POST("", scheduleHandler.Create)
		schedulesAdmin.PUT("/:id", scheduleHandler.Update)
		schedulesAdmin.DELETE("/:id", scheduleHandler.Delete)
	}

	// ticket & booking route
	tickets := r.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware(cfg)) // User harus login
	{
		tickets.GET("/schedules/:id/seats", ticketHandler.GetAvailableSeats)

		tickets.POST("/book", ticketHandler.BookTicket)

		tickets.GET("/me", ticketHandler.GetUserHistory)
	}

	// Transaction & payment route
	transactions := r.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(cfg))
	{
		transactions.GET("/me", transactionHandler.GetUserTransactions)
		transactions.POST("/:id/pay", transactionHandler.PayTransaction)
		transactions.POST("/:id/cancel", transactionHandler.CancelTransaction)
	}

	// Report route (admin only)
	reports := r.Group("/reports")
	reports.Use(middleware.AuthMiddleware(cfg))
	reports.Use(middleware.AdminMiddleware())
	{
		reports.GET("/revenue", reportHandler.GetRevenueReport)
		reports.GET("/revenue/export", reportHandler.ExportRevenueCSV)
		reports.GET("/top-movies", reportHandler.GetTopMovies)
	}

	// Promo route (Admin)
	promos := r.Group("/promos")
	promos.Use(middleware.AuthMiddleware(cfg))
	promos.Use(middleware.AdminMiddleware())
	{
		promos.POST("", promoHandler.Create)
		promos.GET("", promoHandler.GetAll)
		promos.PUT("/:id", promoHandler.Update)
		promos.DELETE("/:id", promoHandler.Delete)
	}
}
