.PHONY: build run test clean migrate-create migrate-up migrate-down migrate-status sqlc swagger docker-build docker-up docker-down help

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := ai-chat-service
MIGRATION_DIR := ./migrations
SWAGGER_DOCS := ./docs/swagger

# Database connection string for Goose
DB_CONNECTION = "postgres://$(shell grep DB_USER .env | cut -d= -f2):$(shell grep DB_PASSWORD .env | cut -d= -f2)@$(shell grep DB_HOST .env | cut -d= -f2):$(shell grep DB_PORT .env | cut -d= -f2)/$(shell grep DB_NAME .env | cut -d= -f2)?sslmode=$(shell grep DB_SSLMODE .env | cut -d= -f2)"

# Build the application
build: sqlc swagger
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) .

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@./bin/$(APP_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -rf $(SWAGGER_DOCS)

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATION_DIR) create $$name sql

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	@goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION) up

# Run migrations down
migrate-down:
	@echo "Running migrations down..."
	@goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION) down

# Check migration status
migrate-status:
	@echo "Migration status:"
	@goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION) status

# Generate SQL code
sqlc:
	@echo "Generating SQL code..."
	@sqlc generate

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@mkdir -p $(SWAGGER_DOCS)
	@swagger generate spec -o $(SWAGGER_DOCS)/swagger.json --scan-models
	@swagger serve -F=swagger $(SWAGGER_DOCS)/swagger.json

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker compose build

docker-up:
	@echo "Starting Docker services..."
	@docker compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker compose down

# Install tools
tools:
	@echo "Installing required tools..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Show help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  migrate-create - Create a new migration file"
	@echo "  migrate-up     - Run migrations up"
	@echo "  migrate-down   - Run migrations down"
	@echo "  migrate-status - Show migration status"
	@echo "  sqlc           - Generate SQL code from query files"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker services"
	@echo "  docker-down    - Stop Docker services"
	@echo "  tools          - Install required development tools"