package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
	"weather-api/internal/infrastructure/support"
	"weather-api/pkg/circuitbreaker"
)

type OpenWeatherAdapter struct {
	client         *http.Client
	apiKey         string
	baseURL        string
	circuitBreaker *circuitbreaker.CircuitBreaker
}

type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`

	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`

	Name string `json:"name"`

	// Field for capturing error messages from the API
	Message string `json:"message"`
}

// NewOpenWeatherAdapter creates a new OpenWeatherAdapter.
// nolint: unused
func NewOpenWeatherAdapter(apiKey string) *OpenWeatherAdapter {
	return &OpenWeatherAdapter{
		client: &http.Client{
			Timeout: 10 * time.Second},
		apiKey:         apiKey,
		baseURL:        "https://api.openweathermap.org/data/2.5/weather",
		circuitBreaker: circuitbreaker.NewCircuitBreaker("openweather-api"),
	}
}

func (a *OpenWeatherAdapter) GetWeatherByCity(city string) (*entity.Weather, error) {
	ctx := context.Background()

	result, err := a.circuitBreaker.Execute(ctx, func() (interface{}, error) {
		return a.fetchWeatherData(city)
	})

	if err != nil {
		return nil, err // Pass the error up, including custom error types
	}

	weather, ok := result.(*entity.Weather)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from circuit breaker")
	}

	return weather, nil
}

// fetchWeatherData makes the actual HTTP request to OpenWeather API
func (a *OpenWeatherAdapter) fetchWeatherData(city string) (*entity.Weather, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", a.baseURL, city, a.apiKey)

	resp, err := a.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() // Properly handle close error

	// Decode the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to decode the error message from the API response
		var apiResp OpenWeatherResponse
		_ = json.Unmarshal(body, &apiResp)

		// If the API returns a 404, we return our custom ErrNotFound
		if resp.StatusCode == http.StatusNotFound {
			msg := fmt.Sprintf("city '%s' not found", city)
			if apiResp.Message != "" {
				msg = apiResp.Message // Use the more specific message from the API if available
			}
			return nil, support.NewErrNotFound(msg)
		}

		// For all other errors, return a generic error
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp OpenWeatherResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode successful response: %w", err)
	}

	description := ""
	if len(apiResp.Weather) > 0 {
		description = apiResp.Weather[0].Description
	}

	weather := &entity.Weather{
		City:        apiResp.Name,
		Temperature: apiResp.Main.Temp,
		Description: description,
		Humidity:    apiResp.Main.Humidity,
		WindSpeed:   apiResp.Wind.Speed,
		Timestamp:   time.Now(),
	}

	return weather, nil
}

var _ repository.WeatherRepository = (*OpenWeatherAdapter)(nil)
