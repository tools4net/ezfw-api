// @title           EZFW API
// @version         1.0
// @description     This is the API for managing EZFW (Easy Firewall/Proxy) configurations for Xray and Sing-box.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    https://github.com/tools4net/ezfw/issues
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API Key for authentication. Prepend with "Bearer " if it's a bearer token, otherwise just the key. For this API, just the key is expected.

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/tools4net/ezfw/backend/internal/store/sqlite"
	// "github.com/tools4net/ezfw/backend/internal/config" // Placeholder for config
)

func main() {
	// // Load configuration (e.g., from .env file or environment variables)
	// cfg, err := config.LoadConfig(".") // Assuming config loader is in internal/config
	// if err != nil {
	//  log.Fatalf("could not load config: %v", err)
	// }

	// Determine data directory (e.g., relative to executable or from env var)
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		// Default to a 'data' directory in the current working directory of the executable
		// This might need adjustment based on how/where the binary is run.
		// For Docker, this path will be inside the container.
		dataDir = "./data"
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory %s: %v", dataDir, err)
	}
	dbPath := filepath.Join(dataDir, "proxypanel.db")
	log.Printf("Using database at: %s", dbPath)

	// Initialize SQLite store
	dbStore, err := sqlite.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite store: %v", err)
	}
	defer dbStore.Close() // Ensure DB is closed when main exits

}
