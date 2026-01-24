.PHONY: help dev-up dev-down dev-logs backend-run frontend-run test lint clean install migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev-up          - Start all development services (Docker Compose)"
	@echo "  make dev-down        - Stop all development services"
	@echo "  make dev-logs        - Show logs from all services"
	@echo "  make backend-run     - Run backend services"
	@echo "  make frontend-run    - Run frontend development server"
	@echo "  make install         - Install all dependencies"
	@echo "  make migrate         - Run database migrations"
	@echo "  make test            - Run all tests"
	@echo "  make lint            - Run linters"
	@echo "  make clean           - Clean build artifacts and dependencies"
	@echo "  make build           - Build all services"

# Docker Compose commands
dev-up:
	@echo "Starting development environment..."
	docker-compose up -d
	@echo "Services started successfully!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "MinIO Console: http://localhost:9001"
	@echo "MinIO API: http://localhost:9000"
	@echo "MailHog UI: http://localhost:8025"
	@echo "ClickHouse: localhost:8123"
	@echo "NATS: localhost:4222"

dev-down:
	@echo "Stopping development environment..."
	docker-compose down

dev-restart:
	@echo "Restarting development environment..."
	docker-compose restart

dev-logs:
	docker-compose logs -f

dev-logs-postgres:
	docker-compose logs -f postgres

dev-logs-redis:
	docker-compose logs -f redis

dev-logs-minio:
	docker-compose logs -f minio

# Installation
install: install-backend install-frontend

install-backend:
	@echo "Installing backend dependencies..."
	cd backend && go mod download && go mod tidy

install-frontend:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Backend commands
backend-run:
	@echo "Running backend API gateway..."
	cd backend && go run cmd/api-gateway/main.go

backend-ws:
	@echo "Running WebSocket server..."
	cd backend && go run cmd/ws-server/main.go

backend-test:
	@echo "Running backend tests..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...

backend-lint:
	@echo "Running backend linter..."
	cd backend && golangci-lint run

backend-build:
	@echo "Building backend services..."
	cd backend && go build -o bin/api-gateway cmd/api-gateway/main.go
	cd backend && go build -o bin/ws-server cmd/ws-server/main.go

# Frontend commands
frontend-run:
	@echo "Running frontend development server..."
	cd frontend && npm run dev

frontend-build:
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-preview:
	@echo "Preview frontend build..."
	cd frontend && npm run preview

frontend-test:
	@echo "Running frontend tests..."
	cd frontend && npm run test

frontend-lint:
	@echo "Running frontend linter..."
	cd frontend && npm run lint

frontend-format:
	@echo "Formatting frontend code..."
	cd frontend && npm run format

# Database commands
migrate:
	@echo "Running database migrations..."
	cd backend && go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	cd backend && go run cmd/migrate/main.go down

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	cd backend && go run cmd/migrate/main.go create $$name

# Testing
test: backend-test frontend-test

test-coverage:
	@echo "Running tests with coverage..."
	cd backend && go test -race -coverprofile=coverage.out -covermode=atomic ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html

# Linting
lint: backend-lint frontend-lint

# Code generation
generate-proto:
	@echo "Generating protobuf code..."
	cd backend/idl/protobuf && protoc --go_out=. --go-grpc_out=. *.proto

generate-thrift:
	@echo "Generating Kitex code..."
	cd backend && kitex -module hertzboard ./idl/thrift/*.thrift

# Cleaning
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/.svelte-kit
	rm -rf frontend/build
	rm -rf backend/coverage.out
	rm -rf backend/coverage.html

clean-all: clean
	@echo "Removing dependencies..."
	rm -rf backend/vendor
	rm -rf frontend/node_modules
	docker-compose down -v

# Building
build: backend-build frontend-build

# Docker build
docker-build:
	@echo "Building Docker images..."
	docker build -t hertzboard-api-gateway -f deploy/docker/Dockerfile.api-gateway .
	docker build -t hertzboard-ws-server -f deploy/docker/Dockerfile.ws-server .
	docker build -t hertzboard-frontend -f deploy/docker/Dockerfile.frontend .

# Development
dev: dev-up
	@echo "Starting development servers in parallel..."
	@make -j2 backend-run frontend-run

# Quality checks
check: lint test

# Initialize project
init:
	@echo "Initializing project..."
	cp .env.example .env
	@make install
	@make dev-up
	@echo "Waiting for services to be ready..."
	sleep 10
	@make migrate
	@echo "Project initialized successfully!"
