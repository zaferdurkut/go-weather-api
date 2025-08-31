package entity

import "time"

type Weather struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Description string    `json:"description"`
	Humidity    int       `json:"humidity"`
	WindSpeed   float64   `json:"wind_speed"`
	Timestamp   time.Time `json:"timestamp"`
}

// WeatherRequest represents a weather request
type WeatherRequest struct {
	City string `json:"city" validate:"required"`
}

// WeatherResponse represents a weather response
type WeatherResponse struct {
	Success bool     `json:"success"`
	Data    *Weather `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
}
