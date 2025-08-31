.PHONY: build test run clean help docker-build docker-run docker-stop docker-clean

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  run           - Run the application (requires OPENWEATHER_API_KEY)"
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

# Run the application
run:
	@if [ -z "$(OPENWEATHER_API_KEY)" ]; then \
		echo "Error: OPENWEATHER_API_KEY environment variable is required"; \
		echo "Please set it: export OPENWEATHER_API_KEY=your_api_key"; \
		exit 1; \
	fi
	go run cmd/server/main.go

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