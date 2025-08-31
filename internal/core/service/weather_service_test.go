package service

import (
	"testing"

	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
)

// MockWeatherRepository is a mock implementation for testing
type MockWeatherRepository struct {
	weather *entity.Weather
	err     error
}

func (m *MockWeatherRepository) GetWeatherByCity(city string) (*entity.Weather, error) {
	return m.weather, m.err
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
	response, err := service.GetWeatherByCity("Istanbul")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}

	if response.Data.City != "Istanbul" {
		t.Errorf("Expected city=Istanbul, got %s", response.Data.City)
	}

	if response.Data.Temperature != 25.5 {
		t.Errorf("Expected temperature=25.5, got %f", response.Data.Temperature)
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
	response, err := service.GetWeatherByCity("InvalidCity")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}

	if response.Error == "" {
		t.Error("Expected error message, got empty string")
	}

	if response.Data != nil {
		t.Error("Expected data to be nil when error occurs")
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
	response, err := service.GetWeatherByCity("")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}
}
