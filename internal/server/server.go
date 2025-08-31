package server

import (
	"fmt"
	"log"
	"net/http"
)

// Run initializes all dependencies and starts the HTTP server.
func Run() {
	// Build the dependency container
	container := BuildContainer()

	// Start server
	cfg := container.Config
	router := container.Router

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Health check: http://localhost%s/health", serverAddr)
	log.Printf("Weather API: http://localhost%s/weather/:city", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
