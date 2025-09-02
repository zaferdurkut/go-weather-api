package main

import (
	_ "weather-api/docs" // This line is necessary for swag to find your docs!
	"weather-api/internal/server"
)

// @title Go Weather API
// @version 1.0
// @description A simple weather API service built with Go, Gin, and Hexagonal Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @BasePath /
func main() {
	server.Run()
}
