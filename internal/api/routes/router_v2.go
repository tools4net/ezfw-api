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

	// Register V2 Configuration Routes (user authentication required)
	registerV2ConfigRoutes(api, configHandler, authMiddleware)

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

func registerV2ConfigRoutes(api huma.API, configHandler *handlers.ConfigHandler, authMiddleware *auth.AuthMiddleware) {
	// =============================================================================
	// SingBox Configuration Routes
	// =============================================================================

	// Create SingBox Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "create-singbox-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/singbox",
			Summary:     "Create SingBox Configuration",
			Description: "Creates a new SingBox proxy configuration",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.CreateSingBoxConfig,
	)

	// List SingBox Configurations
	huma.Register(
		api, huma.Operation{
			OperationID: "list-singbox-configs-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/singbox",
			Summary:     "List SingBox Configurations",
			Description: "Retrieves a paginated list of SingBox configurations",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.ListSingBoxConfigs,
	)

	// Get SingBox Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "get-singbox-config-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/singbox/{config_id}",
			Summary:     "Get SingBox Configuration",
			Description: "Retrieves a specific SingBox configuration by ID",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GetSingBoxConfig,
	)

	// Update SingBox Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "update-singbox-config-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/configs/singbox/{config_id}",
			Summary:     "Update SingBox Configuration",
			Description: "Updates an existing SingBox configuration",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.UpdateSingBoxConfig,
	)

	// Delete SingBox Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "delete-singbox-config-v2",
			Method:      http.MethodDelete,
			Path:        "/api/v2/configs/singbox/{config_id}",
			Summary:     "Delete SingBox Configuration",
			Description: "Deletes a SingBox configuration",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.DeleteSingBoxConfig,
	)

	// Generate SingBox Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "generate-singbox-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/singbox/{config_id}/generate",
			Summary:     "Generate SingBox Configuration",
			Description: "Generates a SingBox configuration file from stored configuration",
			Tags:        []string{"SingBox Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GenerateSingBoxConfig,
	)

	// =============================================================================
	// Xray Configuration Routes
	// =============================================================================

	// Create Xray Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "create-xray-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/xray",
			Summary:     "Create Xray Configuration",
			Description: "Creates a new Xray proxy configuration",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.CreateXrayConfig,
	)

	// List Xray Configurations
	huma.Register(
		api, huma.Operation{
			OperationID: "list-xray-configs-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/xray",
			Summary:     "List Xray Configurations",
			Description: "Retrieves a paginated list of Xray configurations",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.ListXrayConfigs,
	)

	// Get Xray Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "get-xray-config-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/xray/{config_id}",
			Summary:     "Get Xray Configuration",
			Description: "Retrieves a specific Xray configuration by ID",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GetXrayConfig,
	)

	// Update Xray Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "update-xray-config-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/configs/xray/{config_id}",
			Summary:     "Update Xray Configuration",
			Description: "Updates an existing Xray configuration",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.UpdateXrayConfig,
	)

	// Delete Xray Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "delete-xray-config-v2",
			Method:      http.MethodDelete,
			Path:        "/api/v2/configs/xray/{config_id}",
			Summary:     "Delete Xray Configuration",
			Description: "Deletes an Xray configuration",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.DeleteXrayConfig,
	)

	// Generate Xray Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "generate-xray-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/xray/{config_id}/generate",
			Summary:     "Generate Xray Configuration",
			Description: "Generates an Xray configuration file from stored configuration",
			Tags:        []string{"Xray Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GenerateXrayConfig,
	)

	// =============================================================================
	// HAProxy Configuration Routes
	// =============================================================================

	// Create HAProxy Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "create-haproxy-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/haproxy",
			Summary:     "Create HAProxy Configuration",
			Description: "Creates a new HAProxy load balancer configuration",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.CreateHAProxyConfig,
	)

	// List HAProxy Configurations
	huma.Register(
		api, huma.Operation{
			OperationID: "list-haproxy-configs-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/haproxy",
			Summary:     "List HAProxy Configurations",
			Description: "Retrieves a paginated list of HAProxy configurations",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.ListHAProxyConfigs,
	)

	// Get HAProxy Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "get-haproxy-config-v2",
			Method:      http.MethodGet,
			Path:        "/api/v2/configs/haproxy/{config_id}",
			Summary:     "Get HAProxy Configuration",
			Description: "Retrieves a specific HAProxy configuration by ID",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GetHAProxyConfig,
	)

	// Update HAProxy Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "update-haproxy-config-v2",
			Method:      http.MethodPut,
			Path:        "/api/v2/configs/haproxy/{config_id}",
			Summary:     "Update HAProxy Configuration",
			Description: "Updates an existing HAProxy configuration",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.UpdateHAProxyConfig,
	)

	// Delete HAProxy Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "delete-haproxy-config-v2",
			Method:      http.MethodDelete,
			Path:        "/api/v2/configs/haproxy/{config_id}",
			Summary:     "Delete HAProxy Configuration",
			Description: "Deletes an HAProxy configuration",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.DeleteHAProxyConfig,
	)

	// Generate HAProxy Configuration
	huma.Register(
		api, huma.Operation{
			OperationID: "generate-haproxy-config-v2",
			Method:      http.MethodPost,
			Path:        "/api/v2/configs/haproxy/{config_id}/generate",
			Summary:     "Generate HAProxy Configuration",
			Description: "Generates an HAProxy configuration file from stored configuration",
			Tags:        []string{"HAProxy Configurations"},
			Security: []map[string][]string{
				{"BearerAuth": {}},
			},
			Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
		}, configHandler.GenerateHAProxyConfig,
	)
}
