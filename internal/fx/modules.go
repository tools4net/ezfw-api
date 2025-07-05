package fx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tools4net/ezfw/backend/internal/api/handlers"
	"github.com/tools4net/ezfw/backend/internal/api/routes"
	"github.com/tools4net/ezfw/backend/internal/auth"
	"github.com/tools4net/ezfw/backend/internal/store"
	"github.com/tools4net/ezfw/backend/internal/store/sqlite"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Config holds application configuration
type Config struct {
	DataDir string
	Port    string
	DBPath  string
}

// ConfigModule provides configuration dependencies
var ConfigModule = fx.Module(
	"config",
	fx.Provide(NewConfig),
)

// LoggingModule provides logging dependencies
var LoggingModule = fx.Module(
	"logging",
	fx.Provide(NewLogger),
)

// StoreModule provides store dependencies
var StoreModule = fx.Module(
	"store",
	fx.Provide(NewStore),
	fx.Invoke(RegisterStoreHooks),
)

// AuthModule provides authentication dependencies
var AuthModule = fx.Module(
	"auth",
	fx.Provide(NewAuthMiddleware),
)

// HandlersModule provides handler dependencies
var HandlersModule = fx.Module(
	"handlers",
	fx.Provide(NewConfigHandler),
	fx.Provide(NewAgentHandlers),
	fx.Provide(NewAgentTokenHandlers),
)

// RouterModule provides router dependencies
var RouterModule = fx.Module(
	"router",
	fx.Provide(NewRouter),
)

// ServerModule provides HTTP server dependencies
var ServerModule = fx.Module(
	"server",
	fx.Provide(NewHTTPServer),
	fx.Invoke(RegisterServerHooks),
)

// NewConfig creates a new configuration from environment variables
func NewConfig() *Config {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := filepath.Join(dataDir, "proxypanel.db")

	return &Config{
		DataDir: dataDir,
		Port:    port,
		DBPath:  dbPath,
	}
}

// NewLogger creates a new zap logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// NewStore creates a new SQLite store
func NewStore(config *Config, logger *zap.Logger) (store.Store, error) {
	// Create data directory
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory %s: %w", config.DataDir, err)
	}

	logger.Info("Using database", zap.String("path", config.DBPath))

	// Initialize SQLite store
	dbStore, err := sqlite.NewSQLiteStore(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SQLite store: %w", err)
	}

	return dbStore, nil
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(store store.Store, logger *zap.Logger) *auth.AuthMiddleware {
	return auth.NewAuthMiddleware(store, logger)
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(store store.Store) *handlers.ConfigHandler {
	return handlers.NewConfigHandler(store)
}

// NewAgentHandlers creates a new agent handlers instance
func NewAgentHandlers(store store.Store, logger *zap.Logger) *handlers.AgentHandlers {
	return handlers.NewAgentHandlers(store, logger)
}

// NewAgentTokenHandlers creates a new agent token handlers instance
func NewAgentTokenHandlers(store store.Store, logger *zap.Logger) *handlers.AgentTokenHandlers {
	return handlers.NewAgentTokenHandlers(store, logger)
}

// NewRouter creates a new HTTP router
func NewRouter(
	store store.Store,
	logger *zap.Logger,
	authMiddleware *auth.AuthMiddleware,
	agentHandlers *handlers.AgentHandlers,
	agentTokenHandlers *handlers.AgentTokenHandlers,
) (http.Handler, huma.API) {
	router, humaAPI := routes.SetupRouterV2(store, authMiddleware, agentHandlers, agentTokenHandlers)
	logger.Info("Router setup completed with V2 authentication and agent support")
	return router, humaAPI
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(config *Config, router http.Handler, logger *zap.Logger) *http.Server {
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	logger.Info("HTTP server configured", zap.String("port", config.Port))
	return server
}

// RegisterServerHooks registers server lifecycle hooks
func RegisterServerHooks(lc fx.Lifecycle, server *http.Server, config *Config, logger *zap.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Print API information
				fmt.Printf("Starting EZFW API server on port %s\n", config.Port)
				fmt.Printf("OpenAPI documentation available at: http://localhost:%s/docs\n", config.Port)
				fmt.Printf("Health check available at: http://localhost:%s/api/v2/health\n", config.Port)
				fmt.Printf("Root endpoint available at: http://localhost:%s/\n", config.Port)

				logger.Info("Starting HTTP server", zap.String("addr", server.Addr))

				// Start server in a goroutine
				go func() {
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						logger.Fatal("Server failed to start", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Stopping HTTP server")
				return server.Shutdown(ctx)
			},
		},
	)
}

// RegisterStoreHooks registers store lifecycle hooks
func RegisterStoreHooks(lc fx.Lifecycle, store store.Store, logger *zap.Logger) {
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("Closing database store")
				if closer, ok := store.(interface{ Close() error }); ok {
					return closer.Close()
				}
				return nil
			},
		},
	)
}
