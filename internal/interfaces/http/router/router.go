package router

import (
	"weather-api/internal/interfaces/http/handler"
	"weather-api/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(weatherHandler *handler.WeatherHandler, logger *zap.Logger, swaggerBasePath string) *gin.Engine {
	// Create a new router without any default middleware
	router := gin.New()

	// Apply CORS middleware to all incoming requests. This should be one of the first middleware.
	router.Use(middleware.CORS())

	// Request ID must run early to populate context and response header
	router.Use(middleware.RequestID())

	// Add other essential middleware
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Use structured logger middleware for all requests
	router.Use(middleware.Logger(logger))

	// Health check endpoint
	router.GET("/health", weatherHandler.HealthCheck)

	// Weather endpoints
	weatherGroup := router.Group("/weather")
	{
		weatherGroup.GET("/:city", weatherHandler.GetWeatherByCity)
		weatherGroup.GET("/overview", weatherHandler.GetWeatherOverviewByLatLong)
	}

	// Swagger endpoint
	// The URL for the swagger UI is http://localhost:8080/swagger/index.html
	if swaggerBasePath == "" {
		swaggerBasePath = "/swagger"
	}
	router.GET(swaggerBasePath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
