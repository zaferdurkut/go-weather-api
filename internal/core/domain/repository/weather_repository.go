package repository

import (
	"errors"
	"weather-api/internal/core/domain/entity"
)

type WeatherRepository interface {
	GetWeatherByCity(city string) (*entity.Weather, error)
}

var (
	ErrCityNotFound = errors.New("city not found")
	ErrInvalidCity  = errors.New("invalid city name")
	ErrAPIError     = errors.New("external API error")
)
