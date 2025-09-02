package dto

import "time"

// WeatherData defines the structure of the weather data returned to the client.
type WeatherData struct {
	City        string    `json:"city" example:"London"`
	Temperature float64   `json:"temperature" example:"15.5"`
	Description string    `json:"description" example:"scattered clouds"`
	Humidity    int       `json:"humidity" example:"80"`
	WindSpeed   float64   `json:"wind_speed" example:"4.5"`
	Timestamp   time.Time `json:"timestamp"`
}

// WeatherResponse is the generic response wrapper for the weather API.
// It's used for both successful and failed responses.
type WeatherResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    *WeatherData `json:"data,omitempty"`
	Error   string       `json:"error,omitempty" example:"city not found"`
}
