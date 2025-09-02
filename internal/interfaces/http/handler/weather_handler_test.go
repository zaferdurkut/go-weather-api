package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-api/internal/core/domain/entity"
	"weather-api/internal/dto"
	"weather-api/internal/infrastructure/support"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWeatherService is a mock implementation for testing
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeatherByCity(city string) (*entity.Weather, error) {
	args := m.Called(city)
	if w := args.Get(0); w != nil {
		return w.(*entity.Weather), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWeatherService) GetWeatherOverviewByLatLong(lon float32, lat float32) (*entity.WeatherOverview, error) {
	args := m.Called(lon, lat)
	if w := args.Get(0); w != nil {
		return w.(*entity.WeatherOverview), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestWeatherHandler_GetWeatherByCity_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	handler := NewWeatherHandler(mockService)

	expectedWeather := &entity.Weather{
		City:        "Istanbul",
		Temperature: 25.5,
		Description: "sunny",
		Humidity:    60,
		WindSpeed:   10.5,
	}

	mockService.On("GetWeatherByCity", "Istanbul").Return(expectedWeather, nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: "Istanbul"}}

	// Act
	handler.GetWeatherByCity(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.WeatherResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Success)
	assert.NotNil(t, response.Data)
	assert.Equal(t, "Istanbul", response.Data.City)

	mockService.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByCity_EmptyCity(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	handler := NewWeatherHandler(mockService)

	// Create test request with invalid city (fails binding)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: ""}}

	// Act
	handler.GetWeatherByCity(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.WeatherResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.Success)
	assert.NotEmpty(t, response.Error)
}

func TestWeatherHandler_GetWeatherByCity_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	handler := NewWeatherHandler(mockService)

	mockService.On("GetWeatherByCity", "InvalidCity").Return(nil, support.NewErrNotFound("city not found"))

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "city", Value: "InvalidCity"}}

	// Act
	handler.GetWeatherByCity(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["success"])

	mockService.AssertExpectations(t)
}

func TestWeatherHandler_HealthCheck(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockWeatherService)
	handler := NewWeatherHandler(mockService)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Act
	handler.HealthCheck(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "weather-api", response["service"])
}
