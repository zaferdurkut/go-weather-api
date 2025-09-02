package entity

import "time"

// Weather is the core domain model for weather information.
// It is independent of any presentation or database-specific details.
type Weather struct {
	City        string
	Temperature float64
	Description string
	Humidity    int
	WindSpeed   float64
	Timestamp   time.Time
}
