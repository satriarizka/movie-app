run:
	go run cmd/api/main.go

build:
	go build -o bin/main cmd/api/main.go

test:
	go test -v ./...

docker-build:
	docker build -t movie-app .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: run build test docker-build docker-up docker-down migrate-create