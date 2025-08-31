package service

import (
	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
)

// WeatherServiceInterface defines the interface for weather service
type WeatherServiceInterface interface {
	GetWeatherByCity(city string) (*entity.WeatherResponse, error)
}

// WeatherService handles weather business logic
// This is the CORE business logic in Hexagonal Architecture
type WeatherService struct {
	weatherRepo repository.WeatherRepository
}

// NewWeatherService creates a new weather service
func NewWeatherService(weatherRepo repository.WeatherRepository) *WeatherService {
	return &WeatherService{
		weatherRepo: weatherRepo,
	}
}

// GetWeatherByCity retrieves weather information for a given city
func (s *WeatherService) GetWeatherByCity(city string) (*entity.WeatherResponse, error) {
	weather, err := s.weatherRepo.GetWeatherByCity(city)
	if err != nil {
		return &entity.WeatherResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &entity.WeatherResponse{
		Success: true,
		Data:    weather,
	}, nil
}
