package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	Weather WeatherConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// WeatherConfig holds weather API configuration
type WeatherConfig struct {
	APIKey string
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Weather: WeatherConfig{
			APIKey: getEnv("OPENWEATHER_API_KEY", ""),
		},
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
