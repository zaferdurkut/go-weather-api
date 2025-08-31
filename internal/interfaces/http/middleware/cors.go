package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS is a gin middleware for enabling Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	// For a production environment, you should be more restrictive.
	// Example: AllowOrigins: []string{"https://www.your-frontend.com"}
	return cors.New(cors.Config{
		// Allow all origins for development purposes. 
		AllowAllOrigins: true,

		// Allowed methods
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// Allowed headers
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// Expose headers (optional)
		ExposeHeaders: []string{"Content-Length"},

		// Allow credentials (e.g., cookies)
		AllowCredentials: true,

		// MaxAge indicates how long the results of a preflight request can be cached.
		MaxAge: 12 * time.Hour,
	})
}
