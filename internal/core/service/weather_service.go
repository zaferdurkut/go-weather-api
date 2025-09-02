package service

import (
	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
)

// WeatherServiceInterface defines the interface for the core weather business logic.
// It returns a pure domain entity or an error.
type WeatherServiceInterface interface {
	GetWeatherByCity(city string) (*entity.Weather, error)
	GetWeatherOverviewByLatLong(lon float32, lat float32) (*entity.WeatherOverview, error)
}

// WeatherService handles weather business logic.
// This is the CORE business logic in Hexagonal Architecture.
type WeatherService struct {
	weatherRepo repository.WeatherRepository
}

// NewWeatherService creates a new weather service.
func NewWeatherService(weatherRepo repository.WeatherRepository) *WeatherService {
	return &WeatherService{
		weatherRepo: weatherRepo,
	}
}

// GetWeatherByCity retrieves weather information for a given city.
// It returns the core domain model or an error if the data cannot be fetched.
func (s *WeatherService) GetWeatherByCity(city string) (*entity.Weather, error) {
	weather, err := s.weatherRepo.GetWeatherByCity(city)
	if err != nil {
		return nil, err
	}

	return weather, nil
}

func (s *WeatherService) GetWeatherOverviewByLatLong(lon float32, lat float32) (*entity.WeatherOverview, error) {

	weatherOverview, err := s.weatherRepo.GetWeatherOverviewByLatLong(lon, lat)

	if err != nil {
		return nil, err
	}

	return weatherOverview, nil

}
