package handler

import (
	"net/http"

	"weather-api/internal/core/service"
	"weather-api/internal/dto"
	"weather-api/internal/infrastructure/support"

	"github.com/gin-gonic/gin"
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
	// Bind and validate path parameter using URI binding
	type cityURI struct {
		City string `uri:"city" binding:"required,alphaunicode,min=2"`
	}
	var params cityURI
	if err := c.ShouldBindUri(&params); err != nil {
		writeError(c, support.NewErrBadRequest(err.Error()))
		return
	}

	// Call the core service, which returns a pure domain model or an error.
	weather, err := h.weatherService.GetWeatherByCity(params.City)
	if err != nil {
		writeError(c, err)
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

// GetWeatherOverviewByLatLong godoc
// @Summary      Get weather Overview by Lat Lon
// @Description  Retrieves the current weather overview information for a given lat lon.
// @Tags         Weather
// @Accept       json
// @Produce      json
// @Param        lat  query      number  true  "Lat"
// @Param        lon  query      number  true  "Lon"
// @Success      200  {object}  dto.WeatherOverviewResponse  "Successfully retrieved weather data"
// @Failure      400  {object}  dto.WeatherOverviewResponse  "Invalid request (e.g., city name is missing)"
// @Failure      404  {object}  dto.WeatherOverviewResponse  "Weather data not found for the specified city"
// @Failure      500  {object}  dto.WeatherOverviewResponse  "Internal server error"
// @Router       /weather/overview [get]
func (h *WeatherHandler) GetWeatherOverviewByLatLong(c *gin.Context) {
	// Bind and validate query parameters with ranges
	var input struct {
		Lon float32 `form:"lon" binding:"required,gte=-180,lte=180"`
		Lat float32 `form:"lat" binding:"required,gte=-90,lte=90"`
	}

	if err := c.ShouldBindQuery(&input); err != nil {
		writeError(c, support.NewErrBadRequest(err.Error()))
		return
	}

	// Call the core service, which returns a pure domain model or an error.
	weatherOverview, err := h.weatherService.GetWeatherOverviewByLatLong(input.Lon, input.Lat)
	if err != nil {
		writeError(c, err)
		return
	}

	// If successful, map the domain model to the response DTO.
	response := dto.WeatherOverviewResponse{
		Success: true,
		Data: &dto.WeatherOverviewData{
			Lat:             weatherOverview.Lat,
			Lon:             weatherOverview.Lon,
			TZ:              weatherOverview.TZ,
			Date:            weatherOverview.Date,
			Units:           weatherOverview.Units,
			WeatherOverview: weatherOverview.WeatherOverview,
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
