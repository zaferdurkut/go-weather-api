package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather-api/pkg/circuitbreaker"

	"github.com/stretchr/testify/assert"
)

func TestOpenWeatherAdapter_GetWeatherByCity_Success(t *testing.T) {
	// Arrange - Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Query().Get("q"), "Istanbul")
		assert.Contains(t, r.URL.Query().Get("appid"), "test-api-key")
		assert.Equal(t, "metric", r.URL.Query().Get("units"))

		// Mock response
		response := OpenWeatherResponse{
			Main: struct {
				Temp     float64 `json:"temp"`
				Humidity int     `json:"humidity"`
			}{
				Temp:     25.5,
				Humidity: 60,
			},
			Weather: []struct {
				Description string `json:"description"`
			}{
				{Description: "clear sky"},
			},
			Wind: struct {
				Speed float64 `json:"speed"`
			}{
				Speed: 10.5,
			},
			Name: "Istanbul",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create adapter with mock server URL
	adapter := &OpenWeatherAdapter{
		client:         &http.Client{Timeout: 10 * time.Second},
		apiKey:         "test-api-key",
		baseURL:        mockServer.URL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("test-openweather-api"),
		maxAttempts:    1,
		initialBackoff: 50 * time.Millisecond,
		maxBackoff:     100 * time.Millisecond,
	}

	// Act
	weather, err := adapter.GetWeatherByCity("Istanbul")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, weather)
	assert.Equal(t, "Istanbul", weather.City)
	assert.Equal(t, 25.5, weather.Temperature)
	assert.Equal(t, "clear sky", weather.Description)
	assert.Equal(t, 60, weather.Humidity)
	assert.Equal(t, 10.5, weather.WindSpeed)
}

func TestOpenWeatherAdapter_GetWeatherByCity_NotFound(t *testing.T) {
	// Arrange - Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock 404 response
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"cod":"404","message":"city not found"}`))
	}))
	defer mockServer.Close()

	// Create adapter with mock server URL
	adapter := &OpenWeatherAdapter{
		client:         &http.Client{Timeout: 10 * time.Second},
		apiKey:         "test-api-key",
		baseURL:        mockServer.URL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("test-openweather-api"),
	}

	// Act
	weather, err := adapter.GetWeatherByCity("InvalidCity")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Contains(t, err.Error(), "not found")
}

func TestOpenWeatherAdapter_GetWeatherByCity_InvalidResponse(t *testing.T) {
	// Arrange - Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock invalid JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": "json"`))
	}))
	defer mockServer.Close()

	// Create adapter with mock server URL
	adapter := &OpenWeatherAdapter{
		client:         &http.Client{Timeout: 10 * time.Second},
		apiKey:         "test-api-key",
		baseURL:        mockServer.URL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("test-openweather-api"),
	}

	// Act
	weather, err := adapter.GetWeatherByCity("Istanbul")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Contains(t, err.Error(), "failed to decode successful response")
}

func TestOpenWeatherAdapter_GetWeatherByCity_Timeout(t *testing.T) {
	// Arrange - Create mock server that delays response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"Istanbul"}`))
	}))
	defer mockServer.Close()

	// Create adapter with short timeout
	adapter := &OpenWeatherAdapter{
		client:         &http.Client{Timeout: 1 * time.Second},
		apiKey:         "test-api-key",
		baseURL:        mockServer.URL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("test-openweather-api"),
		maxAttempts:    1,
	}

	// Act
	weather, err := adapter.GetWeatherByCity("Istanbul")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, weather)
	assert.Contains(t, err.Error(), "timeout")
}

func TestOpenWeatherAdapter_GetWeatherByCity_EmptyWeatherArray(t *testing.T) {
	// Arrange - Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response with empty weather array
		response := OpenWeatherResponse{
			Main: struct {
				Temp     float64 `json:"temp"`
				Humidity int     `json:"humidity"`
			}{
				Temp:     25.5,
				Humidity: 60,
			},
			Weather: []struct {
				Description string `json:"description"`
			}{}, // Empty array
			Wind: struct {
				Speed float64 `json:"speed"`
			}{
				Speed: 10.5,
			},
			Name: "Istanbul",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create adapter with mock server URL
	adapter := &OpenWeatherAdapter{
		client:         &http.Client{Timeout: 10 * time.Second},
		apiKey:         "test-api-key",
		baseURL:        mockServer.URL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("test-openweather-api"),
	}

	// Act
	weather, err := adapter.GetWeatherByCity("Istanbul")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, weather)
	assert.Equal(t, "Istanbul", weather.City)
	assert.Equal(t, "", weather.Description) // Should be empty string
	assert.Equal(t, 25.5, weather.Temperature)
}
