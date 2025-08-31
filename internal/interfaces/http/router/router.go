package router

import (
	"github.com/gin-gonic/gin"
	"weather-api/internal/interfaces/http/handler"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(weatherHandler *handler.WeatherHandler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", weatherHandler.HealthCheck)

	// Weather endpoints
	weatherGroup := router.Group("/weather")
	{
		weatherGroup.GET("/:city", weatherHandler.GetWeatherByCity)
	}

	return router
}
