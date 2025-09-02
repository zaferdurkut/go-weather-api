.PHONY: build test test-race run run-dev vet lint clean help docker-build docker-run docker-dev docker-stop docker-clean swag swagger-verify

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detector"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint if available"
	@echo "  run           - Run the application (requires OPENWEATHER_API_KEY)"
	@echo "  run-dev       - Run in debug mode (GIN_MODE=debug)"
	@echo "  swag          - Generate Swagger docs"
	@echo "  swagger-verify- Regenerate Swagger and fail if diffs exist"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-dev    - Run development environment with hot reload"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  docker-clean  - Clean Docker containers and images"

# Build the application
build:
	go build -o weather-api cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with race detector
test-race:
	go test -race -v ./...

# Run go vet
vet:
	go vet ./...

# Run golangci-lint if installed
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint"; \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 0; \
	fi

# Run the application
run:
	@if [ -z "$(OPENWEATHER_API_KEY)" ]; then \
		echo "Error: OPENWEATHER_API_KEY environment variable is required"; \
		echo "Please set it: export OPENWEATHER_API_KEY=your_api_key"; \
		exit 1; \
	fi
	go run cmd/server/main.go

# Run in debug mode (local dev)
run-dev:
	GIN_MODE=debug $(MAKE) run

# Generate Swagger docs (requires swag)
swag:
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "swag CLI is not installed. Install: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi
	swag init -g cmd/server/main.go -o docs

# Verify Swagger docs are up-to-date (fails if changes)
swagger-verify: swag
	@if ! git diff --quiet -- docs; then \
		echo "Swagger docs are out of date. Commit generated files."; \
		exit 1; \
	else \
		echo "Swagger docs are up to date."; \
	fi

# Clean build artifacts
clean:
	rm -f weather-api

# Download dependencies
deps:
	go mod tidy
	go mod download

# Docker commands
docker-build:
	docker build -t weather-api .

docker-run:
	@if [ -z "$(OPENWEATHER_API_KEY)" ]; then \
		echo "Error: OPENWEATHER_API_KEY environment variable is required"; \
		echo "Please set it: export OPENWEATHER_API_KEY=your_api_key"; \
		exit 1; \
	fi
	docker-compose up -d

docker-dev:
	@if [ -z "$(OPENWEATHER_API_KEY)" ]; then \
		echo "Error: OPENWEATHER_API_KEY environment variable is required"; \
		echo "Please set it: export OPENWEATHER_API_KEY=your_api_key"; \
		exit 1; \
	fi
	docker-compose --profile dev up -d

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down -v --rmi all
	docker system prune -f