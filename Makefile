.PHONY: help build run dev test clean migrate

help:
	@echo "Jimpitan Backend Go - Available Commands"
	@echo "========================================"
	@echo "make build       - Build the backend binary"
	@echo "make run         - Build and run the backend"
	@echo "make dev         - Run with hot reload (requires air)"
	@echo "make test        - Run tests"
	@echo "make clean       - Clean build artifacts"
	@echo "make deps        - Download dependencies"
	@echo "make migrate     - Run database migrations"

build:
	@echo "Building Jimpitan backend..."
	@cd cmd/server && go build -o ../../backend-go-server

run: build
	@echo "Starting Jimpitan backend..."
	./backend-go-server

dev:
	@echo "Starting Jimpitan backend in dev mode..."
	@command -v air >/dev/null 2>&1 || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -f backend-go-server
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

migrate:
	@echo "Running database migrations..."
	@mysql -h $(DB_HOST) -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < migrations/001_initial_schema.sql
	@mysql -h $(DB_HOST) -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < migrations/002_add_indexes.sql
	@echo "Migrations completed!"
