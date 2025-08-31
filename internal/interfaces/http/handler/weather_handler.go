package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/internal/core/service"
)

// WeatherHandler handles HTTP requests for weather endpoints
// This is an ADAPTER in Hexagonal Architecture
type WeatherHandler struct {
	weatherService service.WeatherServiceInterface
}

// NewWeatherHandler creates a new weather handler
func NewWeatherHandler(weatherService service.WeatherServiceInterface) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

// GetWeatherByCity handles GET /weather/:city requests
func (h *WeatherHandler) GetWeatherByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "city parameter is required",
		})
		return
	}

	response, err := h.weatherService.GetWeatherByCity(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles GET /health requests
func (h *WeatherHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "weather-api",
	})
}
