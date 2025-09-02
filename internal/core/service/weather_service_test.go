package service

import (
	"errors"
	"testing"

	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
)

// MockWeatherRepository is a mock implementation for testing
type MockWeatherRepository struct {
	weather  *entity.Weather
	overview *entity.WeatherOverview
	err      error
}

func (m *MockWeatherRepository) GetWeatherByCity(city string) (*entity.Weather, error) {
	return m.weather, m.err
}

func (m *MockWeatherRepository) GetWeatherOverviewByLatLong(lon float32, lat float32) (*entity.WeatherOverview, error) {
	return m.overview, m.err
}

func TestWeatherService_GetWeatherByCity_Success(t *testing.T) {
	// Arrange
	mockRepo := &MockWeatherRepository{
		weather: &entity.Weather{
			City:        "Istanbul",
			Temperature: 25.5,
			Description: "sunny",
			Humidity:    60,
			WindSpeed:   10.5,
		},
		err: nil,
	}

	service := NewWeatherService(mockRepo)

	// Act
	weather, err := service.GetWeatherByCity("Istanbul")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if weather == nil {
		t.Fatalf("Expected weather, got nil")
	}
	if weather.City != "Istanbul" {
		t.Errorf("Expected city=Istanbul, got %s", weather.City)
	}
	if weather.Temperature != 25.5 {
		t.Errorf("Expected temperature=25.5, got %f", weather.Temperature)
	}
}

func TestWeatherService_GetWeatherByCity_Error(t *testing.T) {
	// Arrange
	mockRepo := &MockWeatherRepository{
		weather: nil,
		err:     repository.ErrCityNotFound,
	}

	service := NewWeatherService(mockRepo)

	// Act
	weather, err := service.GetWeatherByCity("InvalidCity")

	// Assert
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if !errors.Is(err, repository.ErrCityNotFound) {
		t.Errorf("Expected ErrCityNotFound, got %v", err)
	}
	if weather != nil {
		t.Error("Expected nil weather on error")
	}
}

func TestWeatherService_GetWeatherByCity_EmptyCity(t *testing.T) {
	// Arrange
	mockRepo := &MockWeatherRepository{
		weather: nil,
		err:     repository.ErrInvalidCity,
	}

	service := NewWeatherService(mockRepo)

	// Act
	weather, err := service.GetWeatherByCity("")

	// Assert
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if !errors.Is(err, repository.ErrInvalidCity) {
		t.Errorf("Expected ErrInvalidCity, got %v", err)
	}
	if weather != nil {
		t.Error("Expected nil weather on error")
	}
}
