package middleware

import (
	"github.com/gin-contrib/cors"
)

// Add CORS headers for json responses.
func SetCors() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}

	return config
}
