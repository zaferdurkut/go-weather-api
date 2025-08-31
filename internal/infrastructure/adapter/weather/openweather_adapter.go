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
}

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
		return nil, fmt.Errorf("circuit breaker error: %w", err)
	}

	weather, ok := result.(*entity.Weather)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
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

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
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
