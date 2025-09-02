package handler

import (
	"net/http"

	"weather-api/internal/dto"
	"weather-api/internal/infrastructure/support"

	"github.com/gin-gonic/gin"
)

// writeError maps known error types to HTTP status codes and writes a consistent response envelope.
func writeError(c *gin.Context, err error) {
	// Attach error to context so logging middleware can record it for non-4xx as well
	_ = c.Error(err)

	switch e := err.(type) {
	case *support.ErrBadRequest:
		c.JSON(http.StatusBadRequest, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	case *support.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	case *support.ErrForbidden:
		c.JSON(http.StatusForbidden, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	case *support.ErrNotFound:
		c.JSON(http.StatusNotFound, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	case *support.ErrTimeout:
		c.JSON(http.StatusGatewayTimeout, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	case *support.ErrUpstream:
		// Map 502/503 if provided, fallback to 502
		status := e.StatusCode
		if status != http.StatusBadGateway && status != http.StatusServiceUnavailable {
			status = http.StatusBadGateway
		}
		c.JSON(status, dto.WeatherResponse{Success: false, Error: e.Error()})
		return
	}

	// Default: internal error
	c.JSON(http.StatusInternalServerError, dto.WeatherResponse{Success: false, Error: err.Error()})
}
