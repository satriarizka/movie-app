# 🎬 Movie Booking Backend API

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Framework](https://img.shields.io/badge/Framework-Gin-00ADD8?style=flat&logo=go)
![Database](https://img.shields.io/badge/Database-PostgreSQL-336791?style=flat&logo=postgresql)
![Swagger](https://img.shields.io/badge/Docs-Swagger-85EA2D?style=flat&logo=swagger)

A robust, production-ready RESTful API for a Movie Ticket Booking System built with **Go (Golang)** implementing **Clean Architecture**.

## 🚀 Key Features

### 👤 User Management
- **Authentication**: Secure JWT-based Auth (Register & Login).
- **Authorization**: Role-based access control (User vs Admin).
- **History**: View personal booking history.

### 🎥 Movie & Schedule (Master Data)
- **Manage Movies**: CRUD operations for movies (Title, Genre, Duration).
- **Manage Studios**: Studio capacity and layout management.
- **Scheduling**: Dynamic screening schedules with conflict detection.

### 🎫 Booking System
- **Real-time Availability**: Check seat status (Available/Booked) instantly.
- **Concurrency Safe**: Atomic transactions to prevent double-booking.
- **Promo Codes**: Apply fixed or percentage-based discounts.

### 💳 Transactions & Payments
- **Payment Simulation**: Support various payment methods.
- **Auto Cancellation**: Background worker cancels unpaid tickets after 15 minutes.
- **Email Notifications**:
  - Immediate booking confirmation.
  - Automatic reminders 1 hour before the movie starts.

### 📊 Reporting (Admin)
- **Revenue Reports**: View daily or monthly revenue.
- **Export Data**: Download reports as CSV files.
- **Top Movies**: Analytics for best-selling movies.

## 🛠️ Tech Stack

- **Core**: Go (Golang)
- **Framework**: Gin Gonic
- **Database**: PostgreSQL & GORM (ORM)
- **Migration**: Golang-Migrate
- **Configuration**: Viper
- **Logging**: Zap Logger
- **Documentation**: Swaggo (Swagger)
- **Email**: Gomail & Mailpit (SMTP Mock)

## 📂 Project Structure

This project follows the **Standard Go Project Layout** and **Clean Architecture**.

```text
movie-app/
├── cmd/
│   └── api/            # Main entry point (main.go)
├── database/
│   └── migrations/     # SQL Migration files (.sql)
├── docs/               # Auto-generated Swagger documentation
├── internal/           # Private application logic
│   ├── config/         # Configuration loader (Viper)
│   ├── constants/      # Global constants (Context keys, Headers)
│   ├── delivery/       # Transport Layer
│   │   ├── http/       # REST API Handlers, Routes, Middleware, DTOs
│   │   └── worker/     # Background jobs (Scheduler/Cron)
│   ├── domain/         # Enterprise Entities & GORM Structs
│   ├── enums/          # Constants/Enumerations (Roles, Status)
│   ├── repository/     # Database Access Layer (PostgreSQL)
│   └── usecase/        # Business Logic Layer
└── pkg/                # Public shared packages
    ├── database/       # DB Connection setup
    ├── logger/         # Logging utility (Zap)
    ├── mailer/         # Email sender (SMTP/Gomail)
    ├── utils/          # Helper functions (Response, JWT)
    └── validator/      # Input validation logic

## ⚡ Installation & Setup

### Prerequisites
- Go 1.20+ installed
- PostgreSQL installed and running
- Docker (optional, for Mailpit)

1. Clone Repository
```
git clone [https://github.com/username/movie-app.git](https://github.com/username/movie-app.git)
cd movie-app
```

2. Configure Environment
Create a .env file in the root directory:
```
APP_NAME=MovieApp
APP_ENV=development
APP_PORT=8080

# Database Config
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=bioskop_db
DB_PORT=5432
DB_SSLMODE=disable

# JWT Config
JWT_SECRET=YOUR_SUPER_SECRET_KEY
JWT_EXP_TIME=24h

# SMTP Config (Mailpit Local)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=
```

3. Run Mailpit (For Email Testing)
Using Docker to run a local SMTP server:
```
docker run -d -p 8025:8025 -p 1025:1025 axllent/mailpit
```
Access Mailpit UI at: http://localhost:8025

4. Run Application
The application handles database migration automatically on startup.
```
go mod tidy
go run cmd/api/main.go
```

## 📚 API Documentation
Once the server is running, access the comprehensive API documentation via Swagger UI:
```
http://localhost:8080/swagger/index.html
```

## 🗄️ Database Schema (ERD)
![ERD](./erd/movie_app.png)

## 🧪 Testing

### Manual Testing
Use the Swagger UI to test endpoints directly.

1. Create an Admin account via DB or Registration.
2. Login to get the Bearer Token.
3. Use the Authorize button in Swagger to authenticate.

### Background Worker Testing
To test the Auto-Cancel and Reminder features:

1. Book a ticket but do not pay -> Check if status changes to cancelled after 15 mins.
2. Book a ticket for a movie starting in 1 hour -> Check Mailpit for the reminder email.