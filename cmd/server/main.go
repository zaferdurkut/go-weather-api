package main

import (
	"fmt"
	"log"
	"net/http"

	"weather-api/internal/core/service"
	"weather-api/internal/infrastructure/adapter/weather"
	"weather-api/internal/infrastructure/config"
	"weather-api/internal/interfaces/http/handler"
	"weather-api/internal/interfaces/http/router"
)

func main() {
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
	router := router.SetupRouter(weatherHandler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Health check: http://localhost%s/health", serverAddr)
	log.Printf("Weather API: http://localhost%s/weather/:city", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
