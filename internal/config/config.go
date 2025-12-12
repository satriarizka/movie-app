package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// App Config
	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	// Database Config
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// JWT Config
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	JWTExpTime string `mapstructure:"JWT_EXP_TIME"`

	// notifikasi config
	SMTPHost string `mapstructure:"SMTP_HOST"`
	SMTPPort int    `mapstructure:"SMTP_PORT"`
	SMTPUser string `mapstructure:"SMTP_USER"`
	SMTPPass string `mapstructure:"SMTP_PASS"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("SMTP_HOST", "localhost")
	viper.SetDefault("SMTP_PORT", 1025)

	// Jika file .env tidak ditemukan, tidak panic (karena mungkin pakai environment variables asli)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return &config
}
