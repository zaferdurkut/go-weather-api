package server

import (
	"log"
	"net/http"

	"weather-api/internal/core/service"
	"weather-api/internal/infrastructure/adapter/weather"
	"weather-api/internal/infrastructure/config"
	"weather-api/internal/interfaces/http/handler"
	"weather-api/internal/interfaces/http/router"
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
	weatherAdapter := weather.NewOpenWeatherAdapter(cfg.Weather.APIKey)

	// Initialize services
	weatherService := service.NewWeatherService(weatherAdapter)

	// Initialize handlers
	weatherHandler := handler.NewWeatherHandler(weatherService)

	// Setup router
	r := router.SetupRouter(weatherHandler)

	return &Container{
		Router: r,
		Config: cfg,
	}
}
