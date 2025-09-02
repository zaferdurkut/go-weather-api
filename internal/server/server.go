package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	log.Printf("Weather API: http://localhost%s/weather", serverAddr)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Listen for termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Printf("Shutdown signal received, shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
