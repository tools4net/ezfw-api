package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tools4net/ezfw/backend/internal/auth"
	"github.com/tools4net/ezfw/backend/internal/models"
	"github.com/tools4net/ezfw/backend/internal/store"
	"go.uber.org/zap"
)

// AgentHandlers contains handlers for agent communication endpoints
type AgentHandlers struct {
	store  store.Store
	logger *zap.Logger
}

// NewAgentHandlers creates a new instance of agent handlers
func NewAgentHandlers(store store.Store, logger *zap.Logger) *AgentHandlers {
	return &AgentHandlers{
		store:  store,
		logger: logger,
	}
}

// Agent Heartbeat Models
type AgentHeartbeatRequest struct {
	Body struct {
		Status      string                 `json:"status" validate:"required" example:"healthy"` // healthy, degraded, error
		Version     string                 `json:"version" validate:"required" example:"v1.2.3"`
		Uptime      int64                  `json:"uptime" validate:"required" example:"3600"` // seconds
		// SystemInfo  *models.SystemInfo     `json:"system_info,omitempty"` // TODO: Define SystemInfo model
		OSInfo      *models.OSInfo         `json:"os_info,omitempty"`
		Metrics     map[string]interface{} `json:"metrics,omitempty" example:"{\"cpu_usage\": 45.2, \"memory_usage\": 67.8}"`
		Services    []ServiceStatus        `json:"services,omitempty"`
		Timestamp   time.Time              `json:"timestamp" validate:"required" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

type AgentHeartbeatResponse struct {
	Body struct {
		Status      string                 `json:"status" example:"acknowledged"`
		Message     string                 `json:"message" example:"Heartbeat received successfully"`
		Commands    []AgentCommand         `json:"commands,omitempty"`
		Config      map[string]interface{} `json:"config,omitempty"`
		Timestamp   time.Time              `json:"timestamp" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

// Service Configuration Models
type ServiceConfigurationRequest struct {
	ServiceTypes []string `query:"service_types" example:"xray,singbox"`
	Version      string   `query:"version" example:"latest"`
}

type ServiceConfigurationResponse struct {
	Body struct {
		Configurations []ServiceConfiguration `json:"configurations"`
		Version        string                  `json:"version" example:"v1.2.3"`
		Timestamp      time.Time               `json:"timestamp" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

// Service Status Report Models
type ServiceStatusReportRequest struct {
	Body struct {
		Reports   []ServiceStatusReport `json:"reports" validate:"required"`
		Timestamp time.Time             `json:"timestamp" validate:"required" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

type ServiceStatusReportResponse struct {
	Body struct {
		Status    string    `json:"status" example:"processed"`
		Message   string    `json:"message" example:"Status reports processed successfully"`
		Processed int       `json:"processed" example:"5"`
		Timestamp time.Time `json:"timestamp" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

// Supporting Models
type ServiceStatus struct {
	ServiceID   string                 `json:"service_id" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	ServiceType string                 `json:"service_type" example:"xray"`
	Status      string                 `json:"status" example:"running"` // running, stopped, error, starting, stopping
	PID         *int                   `json:"pid,omitempty" example:"1234"`
	Port        *int                   `json:"port,omitempty" example:"8080"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	LastRestart *time.Time             `json:"last_restart,omitempty" example:"2023-01-01T12:00:00Z"`
}

type AgentCommand struct {
	ID          string                 `json:"id" example:"cmd-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Type        string                 `json:"type" example:"restart_service"` // restart_service, update_config, stop_service, start_service
	Target      string                 `json:"target,omitempty" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Timeout     int                    `json:"timeout,omitempty" example:"30"` // seconds
	CreatedAt   time.Time              `json:"created_at" example:"2023-01-01T13:00:00Z"`
}

type ServiceConfiguration struct {
	ServiceID     string                 `json:"service_id" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	ServiceType   string                 `json:"service_type" example:"xray"`
	Configuration map[string]interface{} `json:"configuration"`
	Version       string                 `json:"version" example:"v1.0.0"`
	Checksum      string                 `json:"checksum" example:"sha256:abcdef123456..."`
	UpdatedAt     time.Time              `json:"updated_at" example:"2023-01-01T13:00:00Z"`
}

type ServiceStatusReport struct {
	ServiceID   string                 `json:"service_id" validate:"required" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Status      string                 `json:"status" validate:"required" example:"running"`
	PID         *int                   `json:"pid,omitempty" example:"1234"`
	Port        *int                   `json:"port,omitempty" example:"8080"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	Errors      []string               `json:"errors,omitempty"`
	LastRestart *time.Time             `json:"last_restart,omitempty" example:"2023-01-01T12:00:00Z"`
	Timestamp   time.Time              `json:"timestamp" validate:"required" example:"2023-01-01T13:00:00Z"`
}

// RegisterAgentRoutes registers all agent communication routes
func (h *AgentHandlers) RegisterAgentRoutes(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Agent heartbeat endpoint
	huma.Register(api, huma.Operation{
		OperationID: "agent-heartbeat",
		Method:      http.MethodPost,
		Path:        "/api/v2/agent/heartbeat",
		Summary:     "Agent Heartbeat",
		Description: "Endpoint for agents to report their health status and receive commands",
		Tags:        []string{"Agent Communication"},
		Security: []map[string][]string{
			{"AgentToken": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireAgentAuth},
	}, h.AgentHeartbeat)

	// Service configurations endpoint
	huma.Register(api, huma.Operation{
		OperationID: "agent-get-service-configurations",
		Method:      http.MethodGet,
		Path:        "/api/v2/agent/services/configurations",
		Summary:     "Get Service Configurations",
		Description: "Endpoint for agents to fetch service configurations",
		Tags:        []string{"Agent Communication"},
		Security: []map[string][]string{
			{"AgentToken": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireAgentAuth},
	}, h.GetServiceConfigurations)

	// Service status reports endpoint
	huma.Register(api, huma.Operation{
		OperationID: "agent-service-status-reports",
		Method:      http.MethodPost,
		Path:        "/api/v2/agent/services/status_reports",
		Summary:     "Submit Service Status Reports",
		Description: "Endpoint for agents to submit service status reports",
		Tags:        []string{"Agent Communication"},
		Security: []map[string][]string{
			{"AgentToken": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireAgentAuth},
	}, h.SubmitServiceStatusReports)
}

// AgentHeartbeat handles agent heartbeat requests
func (h *AgentHandlers) AgentHeartbeat(ctx context.Context, req *AgentHeartbeatRequest) (*AgentHeartbeatResponse, error) {
	// Get agent context from authentication middleware
	agentCtx := auth.GetAgentContext(ctx)
	if agentCtx == nil {
		h.logger.Error("Agent context not found in heartbeat request")
		return nil, huma.Error400BadRequest("Invalid agent context")
	}

	h.logger.Info("Received agent heartbeat",
		zap.String("node_id", agentCtx.NodeID),
		zap.String("status", req.Body.Status),
		zap.String("version", req.Body.Version),
	)

	// Update node's agent info
	nodeUpdate := &models.NodeUpdateV2{
		// Note: AgentInfo updates would need to be handled separately
		// as NodeUpdateV2 doesn't include AgentInfo field
	}

	// Note: SystemInfo and OSInfo updates are handled separately
	// NodeUpdateV2 doesn't include these fields directly
	// They would need to be updated through separate API calls if needed

	// Update the node
	_, err := h.store.UpdateNode(ctx, agentCtx.NodeID, nodeUpdate)
	if err != nil {
		h.logger.Error("Failed to update node from heartbeat",
			zap.String("node_id", agentCtx.NodeID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to process heartbeat")
	}

	// TODO: Process service status updates
	// TODO: Generate agent commands based on pending operations
	// TODO: Return configuration updates if needed

	response := &AgentHeartbeatResponse{}
	response.Body.Status = "acknowledged"
	response.Body.Message = "Heartbeat received successfully"
	response.Body.Timestamp = time.Now()
	response.Body.Commands = []AgentCommand{} // TODO: Implement command generation
	response.Body.Config = map[string]interface{}{} // TODO: Implement config updates

	return response, nil
}

// GetServiceConfigurations handles requests for service configurations
func (h *AgentHandlers) GetServiceConfigurations(ctx context.Context, req *ServiceConfigurationRequest) (*ServiceConfigurationResponse, error) {
	// Get agent context from authentication middleware
	agentCtx := auth.GetAgentContext(ctx)
	if agentCtx == nil {
		h.logger.Error("Agent context not found in configuration request")
		return nil, huma.Error400BadRequest("Invalid agent context")
	}

	h.logger.Info("Agent requesting service configurations",
		zap.String("node_id", agentCtx.NodeID),
		zap.Strings("service_types", req.ServiceTypes),
	)

	// Get all service instances for this node
	services, err := h.store.ListServiceInstances(ctx, agentCtx.NodeID, 100, 0)
	if err != nil {
		h.logger.Error("Failed to get service instances",
			zap.String("node_id", agentCtx.NodeID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to retrieve service configurations")
	}

	// Filter services by requested types if specified
	var filteredServices []*models.ServiceInstanceV2
	if len(req.ServiceTypes) > 0 {
		serviceTypeMap := make(map[string]bool)
		for _, serviceType := range req.ServiceTypes {
			serviceTypeMap[serviceType] = true
		}

		for _, service := range services {
			if serviceTypeMap[service.ServiceType] {
				filteredServices = append(filteredServices, service)
			}
		}
	} else {
		filteredServices = services
	}

	// Convert to configuration format
	configurations := make([]ServiceConfiguration, 0, len(filteredServices))
	for _, service := range filteredServices {
		// TODO: Generate actual configuration based on service type
		// This would involve calling the configuration generation engine
		config := ServiceConfiguration{
			ServiceID:     service.ID,
			ServiceType:   service.ServiceType,
			Configuration: map[string]interface{}{"placeholder": "config"}, // TODO: Generate real config
			Version:       "v1.0.0", // TODO: Use actual version
			Checksum:      "sha256:placeholder", // TODO: Calculate real checksum
			UpdatedAt:     service.UpdatedAt,
		}
		configurations = append(configurations, config)
	}

	response := &ServiceConfigurationResponse{}
	response.Body.Configurations = configurations
	response.Body.Version = req.Version
	response.Body.Timestamp = time.Now()

	return response, nil
}

// SubmitServiceStatusReports handles service status report submissions
func (h *AgentHandlers) SubmitServiceStatusReports(ctx context.Context, req *ServiceStatusReportRequest) (*ServiceStatusReportResponse, error) {
	// Get agent context from authentication middleware
	agentCtx := auth.GetAgentContext(ctx)
	if agentCtx == nil {
		h.logger.Error("Agent context not found in status report request")
		return nil, huma.Error400BadRequest("Invalid agent context")
	}

	h.logger.Info("Agent submitting service status reports",
		zap.String("node_id", agentCtx.NodeID),
		zap.Int("report_count", len(req.Body.Reports)),
	)

	processed := 0
	for _, report := range req.Body.Reports {
		// Validate that the service belongs to this node
		service, err := h.store.GetServiceInstance(ctx, agentCtx.NodeID, report.ServiceID)
		if err != nil {
			h.logger.Warn("Service not found for status report",
				zap.String("node_id", agentCtx.NodeID),
				zap.String("service_id", report.ServiceID),
				zap.Error(err),
			)
			continue
		}

		// Update service instance with new status
		update := &models.ServiceInstanceUpdateV2{
			Status: &report.Status,
			// TODO: Add fields for PID, port, metrics, etc.
		}

		_, err = h.store.UpdateServiceInstance(ctx, agentCtx.NodeID, service.ID, update)
		if err != nil {
			h.logger.Error("Failed to update service status",
				zap.String("node_id", agentCtx.NodeID),
				zap.String("service_id", report.ServiceID),
				zap.Error(err),
			)
			continue
		}

		processed++
	}

	response := &ServiceStatusReportResponse{}
	response.Body.Status = "processed"
	response.Body.Message = "Status reports processed successfully"
	response.Body.Processed = processed
	response.Body.Timestamp = time.Now()

	return response, nil
}