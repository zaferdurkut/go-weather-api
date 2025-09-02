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

type WeatherOverviewData struct {
	Lat             float32 `json:"lat" example:"38.4"`
	Lon             float32 `json:"lon" example:"38.4"`
	TZ              string  `json:"tz" example:"+02:00"`
	Date            string  `json:"date" example:"2023-04-27"`
	Units           string  `json:"units" example:"metric"`
	WeatherOverview string  `json:"weather_overview" example:"clear sky"`
}

// WeatherResponse is the generic response wrapper for the weather API.
// It's used for both successful and failed responses.
type WeatherResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    *WeatherData `json:"data,omitempty"`
	Error   string       `json:"error,omitempty" example:"city not found"`
}

type WeatherOverviewResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    *WeatherOverviewData `json:"data,omitempty"`
	Error   string               `json:"error,omitempty" example:"lat lon not found"`
}
