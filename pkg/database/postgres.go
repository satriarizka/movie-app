package database

import (
	"fmt"
	"log"
	"movie-app/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Import library migrasi
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	// 1. Buat DSN untuk GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// 2. Koneksi GORM (Untuk aplikasi)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Jalankan Auto Migration (Menggunakan golang-migrate)
	runMigration(cfg)

	log.Println("Connected to PostgreSQL successfully")
	return db
}

func runMigration(cfg *config.Config) {
	// Buat URL koneksi database khusus untuk library migrate
	// Format: postgres://user:pass@host:port/dbname?sslmode=disable
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Inisialisasi migrasi membaca dari folder "migrations"
	m, err := migrate.New(
		"file://./database/migrations", // Pastikan folder ini ada di root project
		databaseURL,
	)
	if err != nil {
		log.Fatal("Migration initialization failed:", err)
	}

	// Eksekusi Migrasi (UP)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration executed successfully")
}
