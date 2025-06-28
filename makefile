# Project metadata
APP_NAME := nalo-workspace
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Default target
.PHONY: all
all: clean build

# Build the binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# Build for production (with optimizations)
.PHONY: build-prod
build-prod:
	@echo "Building $(APP_NAME) for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# Run the application
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run $(CMD_DIR)/main.go

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Running $(APP_NAME) with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode
.PHONY: test-watch
test-watch:
	@echo "Running tests in watch mode..."
	@if command -v gotestsum > /dev/null; then \
		gotestsum --watch; \
	else \
		echo "Gotestsum not found. Installing..."; \
		go install gotest.tools/gotestsum@latest; \
		gotestsum --watch; \
	fi

# Lint using golangci-lint
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "Golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Format code using gofmt
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -s -w .

# Tidy up go.mod
.PHONY: tidy
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy
	$(GOMOD) verify

# Generate swagger docs
.PHONY: swagger
swagger:
	@echo "Generating swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g $(CMD_DIR)/main.go -o ./docs; \
	else \
		echo "Swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g $(CMD_DIR)/main.go -o ./docs; \
	fi

# Docker commands
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 3000:3000 $(APP_NAME)

.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install development dependencies
.PHONY: install-dev-deps
install-dev-deps:
	@echo "Installing development dependencies..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install gotest.tools/gotestsum@latest

# Run the seed script to create admin user
.PHONY: seed
seed:
	@echo "Creating admin user..."
	$(GOCMD) run ./cmd/seed

# Development setup
.PHONY: dev-setup
dev-setup: install-dev-deps seed

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  build-prod     - Build for production"
	@echo "  run            - Run the application"
	@echo "  dev            - Run with hot reload"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-watch     - Run tests in watch mode"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy go.mod"
	@echo "  swagger        - Generate swagger docs"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start with Docker Compose"
	@echo "  docker-compose-down  - Stop Docker Compose"
	@echo "  clean          - Clean build artifacts"
	@echo "  install-dev-deps - Install development dependencies"
	@echo "  seed           - Create admin user"
	@echo "  dev-setup      - Setup development environment"
	@echo "  help           - Show this help"

