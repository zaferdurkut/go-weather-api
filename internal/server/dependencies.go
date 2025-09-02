package server

import (
	"log"
	"net/http"

	"weather-api/internal/core/service"
	"weather-api/internal/infrastructure/adapter/weather"
	"weather-api/internal/infrastructure/config"
	"weather-api/internal/infrastructure/support"
	"weather-api/internal/interfaces/http/handler"
	"weather-api/internal/interfaces/http/router"

	"github.com/gin-gonic/gin"
)

// Container holds all the dependencies for the application.
type Container struct {
	Router http.Handler
	Config *config.Config
}

// BuildContainer creates and wires all the application dependencies.
func BuildContainer() *Container {
	// Load configuration
	cfg := config.LoadConfig()

	// Validate required configuration
	if cfg.Weather.APIKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable is required")
	}

	// Initialize adapters
	weatherAdapter := weather.NewOpenWeatherAdapterWithConfig(cfg.Weather)

	// Initialize services
	weatherService := service.NewWeatherService(weatherAdapter)

	// Initialize handlers
	weatherHandler := handler.NewWeatherHandler(weatherService)

	// Configure Gin mode before creating the router (debug|release|test)
	if cfg.Server.GinMode != "" {
		gin.SetMode(cfg.Server.GinMode)
	}

	// Initialize structured logger
	logger, err := support.NewLogger(cfg.Server.GinMode)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	// Setup router with logger and swagger base path
	r := router.SetupRouter(weatherHandler, logger, cfg.Swagger.BasePath)

	return &Container{
		Router: r,
		Config: cfg,
	}
}
