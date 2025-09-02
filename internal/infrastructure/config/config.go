package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	Weather WeatherConfig
	Swagger SwaggerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	GinMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// WeatherConfig holds weather API configuration
type WeatherConfig struct {
	APIKey              string
	BaseURL             string
	HTTPTimeout         time.Duration
	RetryMaxAttempts    int
	RetryInitialBackoff time.Duration
	RetryMaxBackoff     time.Duration
}

// SwaggerConfig holds Swagger related configuration
type SwaggerConfig struct {
	BasePath string
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			GinMode:      getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", "10s"),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", "15s"),
			IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", "60s"),
		},
		Weather: WeatherConfig{
			APIKey:              getEnv("OPENWEATHER_API_KEY", ""),
			BaseURL:             getEnv("OPENWEATHER_BASE_URL", "https://api.openweathermap.org"),
			HTTPTimeout:         getEnvDuration("OPENWEATHER_HTTP_TIMEOUT", "10s"),
			RetryMaxAttempts:    getEnvInt("OPENWEATHER_RETRY_MAX_ATTEMPTS", 2),
			RetryInitialBackoff: getEnvDuration("OPENWEATHER_RETRY_INITIAL_BACKOFF", "200ms"),
			RetryMaxBackoff:     getEnvDuration("OPENWEATHER_RETRY_MAX_BACKOFF", "2s"),
		},
		Swagger: SwaggerConfig{
			BasePath: getEnv("SWAGGER_BASE_PATH", "/swagger"),
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

// getEnvDuration gets a time.Duration from env with fallback string parsed via time.ParseDuration
func getEnvDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("invalid duration for %s=%q, using fallback %s: %v", key, value, fallback, err)
		fd, _ := time.ParseDuration(fallback)
		return fd
	}
	return d
}

// getEnvInt gets an int from env with fallback
func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
		log.Printf("invalid int for %s=%q, using fallback %d", key, value, fallback)
	}
	return fallback
}
