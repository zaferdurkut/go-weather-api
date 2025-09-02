package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/internal/core/domain/entity"
	"weather-api/internal/core/service"
	"weather-api/internal/dto"
	"weather-api/internal/infrastructure/support"
)

// WeatherHandler handles HTTP requests for weather endpoints.
// Its responsibility is to receive requests, call the appropriate service,
// and map the results (domain models or errors) to DTOs for the response.
type WeatherHandler struct {
	weatherService service.WeatherServiceInterface
}

// NewWeatherHandler creates a new weather handler.
func NewWeatherHandler(weatherService service.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

// GetWeatherByCity godoc
// @Summary      Get weather by city
// @Description  Retrieves the current weather information for a given city name.
// @Tags         Weather
// @Accept       json
// @Produce      json
// @Param        city  path      string  true  "City name"
// @Success      200  {object}  dto.WeatherResponse  "Successfully retrieved weather data"
// @Failure      400  {object}  dto.WeatherResponse  "Invalid request (e.g., city name is missing)"
// @Failure      404  {object}  dto.WeatherResponse  "Weather data not found for the specified city"
// @Failure      500  {object}  dto.WeatherResponse  "Internal server error"
// @Router       /weather/{city} [get]
func (h *WeatherHandler) GetWeatherByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, dto.WeatherResponse{
			Success: false,
			Error:   "city parameter is required",
		})
		return
	}

	// Call the core service, which returns a pure domain model or an error.
	weather, err := h.weatherService.GetWeatherByCity(city)
	if err != nil {
		// If there is an error, map it to the appropriate HTTP status and DTO.
		// Here we can check for specific error types.
		var notFoundErr *support.ErrNotFound
		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, dto.WeatherResponse{
				Success: false,
				Error:   err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.WeatherResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
		return
	}

	// If successful, map the domain model to the response DTO.
	response := dto.WeatherResponse{
		Success: true,
		Data: &dto.WeatherData{
			City:        weather.City,
			Temperature: weather.Temperature,
			Description: weather.Description,
			Humidity:    weather.Humidity,
			WindSpeed:   weather.WindSpeed,
			Timestamp:   weather.Timestamp,
		},
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck godoc
// @Summary      Service Health Check
// @Description  Checks if the weather service is up and running.
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string "Healthy response"
// @Router       /health [get]
func (h *WeatherHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "weather-api",
	})
}
