package router

import (
	"github.com/gin-gonic/gin"
	"weather-api/internal/interfaces/http/handler"
	"weather-api/internal/interfaces/http/middleware"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(weatherHandler *handler.WeatherHandler) *gin.Engine {
	// Create a new router without any default middleware
	router := gin.New()

	// Apply CORS middleware to all incoming requests. This should be one of the first middleware.
	router.Use(middleware.CORS())

	// Add other essential middleware
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Use our custom logger middleware for all requests
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", weatherHandler.HealthCheck)

	// Weather endpoints
	weatherGroup := router.Group("/weather")
	{
		weatherGroup.GET("/:city", weatherHandler.GetWeatherByCity)
	}

	return router
}
