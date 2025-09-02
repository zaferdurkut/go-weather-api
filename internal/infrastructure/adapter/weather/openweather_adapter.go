package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/domain/repository"
	"weather-api/internal/infrastructure/config"
	"weather-api/internal/infrastructure/support"
	"weather-api/pkg/circuitbreaker"
)

type OpenWeatherAdapter struct {
	client         *http.Client
	apiKey         string
	baseURL        string
	circuitBreaker *circuitbreaker.CircuitBreaker
	maxAttempts    int
	initialBackoff time.Duration
	maxBackoff     time.Duration
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

type OpenWeatherOverviewResponse struct {
	Lat             float32 `json:"lat"`
	Lon             float32 `json:"lon"`
	TZ              string  `json:"tz"`
	Date            string  `json:"date"`
	Units           string  `json:"units"`
	WeatherOverview string  `json:"weather_overview"`
}

// NewOpenWeatherAdapter creates a new OpenWeatherAdapter.
// nolint: unused
func NewOpenWeatherAdapterWithConfig(cfg config.WeatherConfig) *OpenWeatherAdapter {
	return &OpenWeatherAdapter{
		client:         &http.Client{Timeout: cfg.HTTPTimeout},
		apiKey:         cfg.APIKey,
		baseURL:        cfg.BaseURL,
		circuitBreaker: circuitbreaker.NewCircuitBreaker("openweather-api"),
		maxAttempts:    cfg.RetryMaxAttempts,
		initialBackoff: cfg.RetryInitialBackoff,
		maxBackoff:     cfg.RetryMaxBackoff,
	}
}

// NewOpenWeatherAdapter creates a new OpenWeatherAdapter.
// nolint: unused
func NewOpenWeatherAdapter(apiKey string) *OpenWeatherAdapter {
	return &OpenWeatherAdapter{
		client: &http.Client{
			Timeout: 10 * time.Second},
		apiKey:         apiKey,
		baseURL:        "https://api.openweathermap.org",
		circuitBreaker: circuitbreaker.NewCircuitBreaker("openweather-api"),
		maxAttempts:    2,
		initialBackoff: 200 * time.Millisecond,
		maxBackoff:     2 * time.Second,
	}
}

func (a *OpenWeatherAdapter) doGetWithRetry(url string) (*http.Response, error) {
	var attempt int
	backoff := a.initialBackoff
	for {
		resp, err := a.client.Get(url)
		if err != nil {
			// Wrap timeout-style errors with a clear prefix
			var ne net.Error
			if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &ne) && ne.Timeout()) {
				return nil, fmt.Errorf("timeout: %w", err)
			}
			return nil, err
		}
		if resp.StatusCode < 500 {
			return resp, nil
		}
		attempt++
		if attempt >= a.maxAttempts {
			return resp, nil
		}
		if backoff > a.maxBackoff {
			backoff = a.maxBackoff
		}
		time.Sleep(backoff)
		backoff *= 2
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

func (a *OpenWeatherAdapter) GetWeatherOverviewByLatLong(lon float32, lat float32) (*entity.WeatherOverview, error) {
	ctx := context.Background()

	result, err := a.circuitBreaker.Execute(ctx, func() (interface{}, error) {
		return a.fetchWeatherOverviewData(lon, lat)
	})

	if err != nil {
		return nil, err // Pass the error up, including custom error types
	}

	weatherOverview, ok := result.(*entity.WeatherOverview)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from circuit breaker")
	}

	return weatherOverview, nil
}

// fetchWeatherData makes the actual HTTP request to OpenWeather API
func (a *OpenWeatherAdapter) fetchWeatherData(city string) (*entity.Weather, error) {
	url := fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s&units=metric", a.baseURL, city, a.apiKey)

	resp, err := a.doGetWithRetry(url)
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

// fetchWeatherData makes the actual HTTP request to OpenWeather API
func (a *OpenWeatherAdapter) fetchWeatherOverviewData(lon float32, lat float32) (*entity.WeatherOverview, error) {
	url := fmt.Sprintf("%s/data/3.0/onecall/overview?appid=%s&lat=%f&lon=%f", a.baseURL, a.apiKey, lon, lat)

	resp, err := a.doGetWithRetry(url)
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
			msg := fmt.Sprintf("lon '%f' , lat '%f' not found", lon, lat)
			if apiResp.Message != "" {
				msg = apiResp.Message // Use the more specific message from the API if available
			}
			return nil, support.NewErrNotFound(msg)
		}

		// For all other errors, return a generic error
		return nil, fmt.Errorf("fetchWeatherOverviewData API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp OpenWeatherOverviewResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("fetchWeatherOverviewData failed to decode successful response: %w", err)
	}

	weatherOverview := &entity.WeatherOverview{
		Lat:             apiResp.Lat,
		Lon:             apiResp.Lon,
		TZ:              apiResp.TZ,
		Date:            apiResp.Date,
		Units:           apiResp.Units,
		WeatherOverview: apiResp.WeatherOverview,
	}

	return weatherOverview, nil
}

var _ repository.WeatherRepository = (*OpenWeatherAdapter)(nil)
