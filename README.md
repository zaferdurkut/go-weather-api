# Weather API

A robust Go Weather API built with Hexagonal Architecture (Ports and Adapters Pattern), featuring Circuit Breaker pattern, comprehensive testing, and RESTful API design.

## 🏗️ Architecture

This project follows **Hexagonal Architecture** principles, ensuring clean separation of concerns and high testability:

```
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── core/                       # Core Business Logic
│   │   ├── domain/
│   │   │   ├── entity/             # Domain entities (Weather, WeatherRequest, WeatherResponse)
│   │   │   └── repository/         # Repository interfaces (Ports)
│   │   └── service/                # Business logic services
│   ├── infrastructure/             # External Dependencies
│   │   ├── adapter/
│   │   │   └── weather/            # OpenWeather API adapter
│   │   └── config/                 # Configuration management
│   └── interfaces/                 # Interface Adapters
│       └── http/
│           ├── handler/            # HTTP request handlers
│           └── router/             # Route definitions
├── pkg/
│   └── circuitbreaker/             # Circuit Breaker implementation
└── go.mod
```

## ✨ Features

- **🏛️ Hexagonal Architecture**: Clean separation between business logic and external dependencies
- **⚡ Circuit Breaker Pattern**: Fault tolerance for external API calls using Sony gobreaker
- **🌐 RESTful API**: Clean HTTP endpoints with proper status codes
- **🌤️ OpenWeather Integration**: Real-time weather data from OpenWeather API
- **🧪 Comprehensive Testing**: 100% service layer coverage, 92.3% adapter coverage, 85.7% handler coverage
- **⚙️ Configuration Management**: Environment-based configuration with .env support
- **🔧 Dependency Injection**: Loose coupling between components
- **📦 Go Modules**: Modern dependency management

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2 or higher
- OpenWeather API key ([Get one here](https://openweathermap.org/api))
- Docker and Docker Compose (for containerized deployment)

### Installation

#### Option 1: Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-weather-api
   ```

2. **Install dependencies**
   ```bash
   make deps
   # or
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env and add your OpenWeather API key
   export OPENWEATHER_API_KEY=your_api_key_here
   ```

4. **Run the application**
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

#### Option 2: Docker Deployment

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-weather-api
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env and add your OpenWeather API key
   export OPENWEATHER_API_KEY=your_api_key_here
   ```

3. **Run with Docker Compose**
   ```bash
   # Production
   make docker-run
   
   # Development with hot reload
   make docker-dev
   ```

The server will start on `http://localhost:8080` (production) or `http://localhost:8081` (development)

## 📘 API Documentation (Swagger)

- **View in Browser**
  - Local: [Swagger UI](http://localhost:8080/swagger/index.html)
  - Docker Dev: [Swagger UI (Dev)](http://localhost:8081/swagger/index.html)

- **Regenerate Swagger Docs**
  1. Install generator (once):
     ```bash
     go install github.com/swaggo/swag/cmd/swag@latest
     ```
  2. Generate/update docs into the `docs/` folder:
     ```bash
     swag init -g cmd/server/main.go -o docs
     ```

- **Notes**
  - Swagger UI is served at route prefix `/swagger/*any`.
  - Keep the import `_ "weather-api/docs"` in `cmd/server/main.go` so the UI can load the generated spec.

## 📡 API Endpoints

### Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "service": "weather-api"
}
```

### Get Weather by City
```http
GET /weather/{city}
```
**Example:**
```bash
curl http://localhost:8080/weather/Istanbul
```

**Response:**
```json
{
  "success": true,
  "data": {
    "city": "Istanbul",
    "temperature": 25.5,
    "description": "clear sky",
    "humidity": 60,
    "wind_speed": 10.5,
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "city not found"
}
```

## 🧪 Testing

### Run All Tests
```bash
make test
# or
go test -v ./...
```

### Test Coverage
```bash
go test -v -cover ./...
```

**Current Coverage:**
- **Service Layer**: 100% ✅
- **Adapter Layer**: 92.3% ✅
- **Handler Layer**: 85.7% ✅

### Test Categories

1. **Service Tests**: Business logic testing with mock repositories
2. **Adapter Tests**: OpenWeather API integration testing with mock HTTP server
3. **Handler Tests**: HTTP request/response testing with mock services

## 🛠️ Development

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run all tests
make run           # Run the application
make clean         # Clean build artifacts
make deps          # Download dependencies

# Docker commands
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose (production)
make docker-dev    # Run development environment with hot reload
make docker-stop   # Stop Docker containers
make docker-clean  # Clean Docker containers and images
```

### Project Structure

- **`cmd/server/`**: Application entry point and dependency injection
- **`internal/core/`**: Business logic and domain models
- **`internal/infrastructure/`**: External service adapters and configuration
- **`internal/interfaces/`**: HTTP handlers and routing
- **`pkg/circuitbreaker/`**: Reusable circuit breaker implementation

## 🔧 Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (`debug`, `release`, `test`) | `debug` |
| `READ_TIMEOUT` | Server read timeout | `10s` |
| `WRITE_TIMEOUT` | Server write timeout | `15s` |
| `IDLE_TIMEOUT` | Server idle timeout | `60s` |
| `OPENWEATHER_API_KEY` | OpenWeather API key | Required |
| `OPENWEATHER_BASE_URL` | OpenWeather API base URL | `https://api.openweathermap.org` |
| `OPENWEATHER_HTTP_TIMEOUT` | OpenWeather HTTP client timeout | `10s` |
| `OPENWEATHER_RETRY_MAX_ATTEMPTS` | Retry attempts for adapter | `2` |
| `OPENWEATHER_RETRY_INITIAL_BACKOFF` | Initial backoff duration | `200ms` |
| `OPENWEATHER_RETRY_MAX_BACKOFF` | Max backoff duration | `2s` |
| `SWAGGER_BASE_PATH` | Swagger UI base path | `/swagger` |

### Docker Configuration

The application includes Docker support with the following features:

- **Multi-stage builds** for optimized production images
- **Development environment** with hot reload using Air
- **Health checks** for container monitoring
- **Non-root user** for security
- **Volume mounting** for development
- **Network isolation** with custom bridge network

## 🏛️ Architecture Benefits

1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Easy to swap implementations (e.g., different weather APIs)
4. **Scalability**: Modular design allows for easy extension
5. **Reliability**: Circuit breaker pattern provides fault tolerance

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [OpenWeather API](https://openweathermap.org/api) for weather data
- [Gin](https://github.com/gin-gonic/gin) for HTTP framework
- [Sony gobreaker](https://github.com/sony/gobreaker) for circuit breaker implementation
- [Testify](https://github.com/stretchr/testify) for testing utilities