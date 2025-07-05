package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tools4net/ezfw/backend/internal/api/handlers"
	"github.com/tools4net/ezfw/backend/internal/auth"
	"github.com/tools4net/ezfw/backend/internal/store"
)

// SetupRouterV2 initializes the Huma API with V2 authentication and agent support
func SetupRouterV2(
	dbStore store.Store,
	authMiddleware *auth.AuthMiddleware,
	agentHandlers *handlers.AgentHandlers,
	agentTokenHandlers *handlers.AgentTokenHandlers,
) (http.Handler, huma.API) {
	// Create Chi router
	router := chi.NewRouter()

	// Add Chi middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Create Huma API configuration
	config := huma.DefaultConfig("EZFW API V2", "2.0.0")
	config.Info.Description = "This is the V2 API for managing EZFW (Easy Firewall/Proxy) configurations with multi-service architecture, dual authentication, and agent communication support."
	config.Info.Contact = &huma.Contact{
		Name:  "API Support",
		URL:   "https://github.com/tools4net/ezfw/issues",
		Email: "support@example.com",
	}
	config.Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	config.Servers = []*huma.Server{
		{URL: "http://localhost:8080", Description: "Development server"},
	}

	// Add security schemes for dual authentication
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Clerk JWT token for user authentication",
		},
		"AgentToken": {
			Type:        "apiKey",
			In:          "header",
			Name:        "X-Agent-Token",
			Description: "Agent authentication token for node communication",
		},
	}

	// Create Huma API
	api := humachi.New(router, config)

	// Create handlers
	nodeHandler := handlers.NewNodeHandler(dbStore)
	serviceHandler := handlers.NewServiceHandler(dbStore)
	configHandler := handlers.NewConfigHandler(dbStore)

	// Basic root health check
	router.Get(
		"/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "EZFW API V2 is running!", "version": "2.0.0"}`))
		},
	)

	// =============================================================================
	// V2 API Routes (New architecture with dual authentication)
	// =============================================================================

	// V2 Health check with dual auth support
	huma.Register(
		api, huma.Operation{
			OperationID: "health-check-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/health",
			Summary:     "Health Check (V2)",
			Description: "Check if the API v2 is healthy (supports both user and agent authentication)",
			Tags:        []string{"Health", "V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
				{"AgentToken": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireAuth},
		}, configHandler.HealthCheck,
	)

	// Register V2 Node Management Routes (user authentication required)
	registerV2NodeRoutes(api, nodeHandler, authMiddleware)

	// Register V2 Service Instance Routes (user authentication required)
	registerV2ServiceRoutes(api, serviceHandler, authMiddleware)

	// Register Agent Communication Routes (agent authentication required)
	agentHandlers.RegisterAgentRoutes(api, authMiddleware)

	// Register Agent Token Management Routes (user authentication required)
	agentTokenHandlers.RegisterAgentTokenRoutes(api, authMiddleware)

	return router, api
}

// =============================================================================
// V2 Route Registration Functions
// =============================================================================

func registerV2NodeRoutes(api huma.API, nodeHandler *handlers.NodeHandler, authMiddleware *auth.AuthMiddleware) {
	// Create Node V2
	huma.Register(
		api, huma.Operation{
			OperationID: "create-node-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/nodes",
			Summary:     "Create Node (V2)",
			Description: "Creates a new node with enhanced features",
			Tags:        []string{"Nodes V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, nodeHandler.CreateNode,
	)

	// List Nodes V2
	huma.Register(
		api, huma.Operation{
			OperationID: "list-nodes-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/nodes",
			Summary:     "List Nodes (V2)",
			Description: "Retrieves a list of nodes with enhanced filtering",
			Tags:        []string{"Nodes V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, nodeHandler.ListNodes,
	)

	// Get Node V2
	huma.Register(
		api, huma.Operation{
			OperationID: "get-node-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/nodes/{id}",
			Summary:     "Get Node (V2)",
			Description: "Retrieves a specific node by ID with enhanced details",
			Tags:        []string{"Nodes V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, nodeHandler.GetNode,
	)

	// Update Node V2
	huma.Register(
		api, huma.Operation{
			OperationID: "update-node-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/nodes/{id}",
			Summary:     "Update Node (V2)",
			Description: "Updates an existing node with enhanced features",
			Tags:        []string{"Nodes V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, nodeHandler.UpdateNode,
	)

	// Delete Node V2
	huma.Register(
		api, huma.Operation{
			OperationID: "delete-node-v2",
			Method:      http.MethodDelete,
			Path:        "/api/v2/nodes/{id}",
			Summary:     "Delete Node (V2)",
			Description: "Deletes a node and all associated resources",
			Tags:        []string{"Nodes V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, nodeHandler.DeleteNode,
	)
}

func registerV2ServiceRoutes(api huma.API, serviceHandler *handlers.ServiceHandler, authMiddleware *auth.AuthMiddleware) {
	// Create Service Instance V2
	huma.Register(
		api, huma.Operation{
			OperationID: "create-service-instance-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/nodes/{node_id}/services",
			Summary:     "Create Service Instance (V2)",
			Description: "Creates a new service instance on a node with enhanced configuration",
			Tags:        []string{"Service Instances V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.CreateService,
	)

	// List Service Instances V2
	huma.Register(
		api, huma.Operation{
			OperationID: "list-service-instances-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/nodes/{node_id}/services",
			Summary:     "List Service Instances (V2)",
			Description: "Retrieves service instances for a node with enhanced details",
			Tags:        []string{"Service Instances V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.ListServices,
	)

	// Get Service Instance V2
	huma.Register(
		api, huma.Operation{
			OperationID: "get-service-instance-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}",
			Summary:     "Get Service Instance (V2)",
			Description: "Retrieves a specific service instance with enhanced details",
			Tags:        []string{"Service Instances V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.GetService,
	)

	// Update Service Instance V2
	huma.Register(
		api, huma.Operation{
			OperationID: "update-service-instance-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}",
			Summary:     "Update Service Instance (V2)",
			Description: "Updates a service instance with enhanced configuration options",
			Tags:        []string{"Service Instances V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.UpdateService,
	)

	// Delete Service Instance V2
	huma.Register(
		api, huma.Operation{
			OperationID: "delete-service-instance-v2",
			Method:      http.MethodDelete,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}",
			Summary:     "Delete Service Instance (V2)",
			Description: "Deletes a service instance and cleans up associated resources",
			Tags:        []string{"Service Instances V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.DeleteService,
	)

	// Get Service Configuration V2
	huma.Register(
		api, huma.Operation{
			OperationID: "get-service-config-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}/config",
			Summary:     "Get Service Configuration (V2)",
			Description: "Retrieves the configuration for a specific service instance",
			Tags:        []string{"Service Configuration V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.GetServiceConfig,
	)

	// Update Service Configuration V2
	huma.Register(
		api, huma.Operation{
			OperationID: "update-service-config-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}/config",
			Summary:     "Update Service Configuration (V2)",
			Description: "Updates the configuration for a specific service instance",
			Tags:        []string{"Service Configuration V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.UpdateServiceConfig,
	)

	// Generate Service Configuration V2
	huma.Register(
		api, huma.Operation{
			OperationID: "generate-service-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/nodes/{node_id}/services/{service_id}/config/generate",
			Summary:     "Generate Service Configuration (V2)",
			Description: "Generates a new configuration for a specific service instance based on parameters",
			Tags:        []string{"Service Configuration V2"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, serviceHandler.GenerateServiceConfig,
	)
}
