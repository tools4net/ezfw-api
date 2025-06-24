package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	apiKeyHeader = "X-API-Key"
	envAPIKey    = "EZFW_API_KEY"
)

// APIKeyAuthMiddleware checks for a valid API key.
// It's a simplified authentication for now.
func APIKeyAuthMiddleware() gin.HandlerFunc {
	// Get the expected API key from environment variable ONCE during middleware setup
	expectedAPIKey := os.Getenv(envAPIKey)

	if expectedAPIKey == "" {
		// Log this or handle it more gracefully if the API should still run in a "no-auth" mode.
		// For now, if no key is set, it means any key will be rejected (unless it's also empty).
		// This effectively disables the API if the key isn't configured.
		// Consider logging a warning: log.Println("Warning: EZFW_API_KEY is not set. API will be inaccessible.")
	}

	return func(c *gin.Context) {
		// Allow OPTIONS requests to pass through for CORS preflight
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		providedKey := c.GetHeader(apiKeyHeader)

		if expectedAPIKey == "" { // If no API key is configured on the server
			// This case means the admin hasn't set up an API key.
			// Depending on policy, you might allow access or deny it.
			// For heightened security, if no key is set, deny all.
			// Or, for ease of initial setup, you might allow if no key is set (less secure).
			// Let's choose to deny if not set, as it's safer.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API access requires configuration (server API key not set)."})
			c.Abort()
			return
		}

		if providedKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required. Provide it in the '" + apiKeyHeader + "' header."})
			c.Abort()
			return
		}

		if !strings.EqualFold(providedKey, expectedAPIKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key."})
			c.Abort()
			return
		}

		c.Next()
	}
}
