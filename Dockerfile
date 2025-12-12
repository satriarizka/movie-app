# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies yang diperlukan untuk build (jika ada CGO)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary dengan nama 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Install CA certificates untuk HTTPS call
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dari stage builder
COPY --from=builder /app/main .
# Copy .env (opsional, sebaiknya di-inject via environment variable saat run container)
COPY .env .

EXPOSE 8080

CMD ["./main"]