package types

import (
	"time"

	"github.com/tools4net/ezfw/backend/internal/models"
)

// Health check response
type HealthResponse struct {
	Body struct {
		Status string `json:"status" example:"API v2 is healthy"`
	}
}

// Generic error response
type ErrorResponse struct {
	Body struct {
		Error string `json:"error" example:"Detailed error message"`
	}
}

// SingBox Config Operations

// CreateSingBoxConfigInput represents the input for creating a SingBox config
type CreateSingBoxConfigInput struct {
	Body models.SingBoxConfig `json:"singBoxConfig"`
}

// CreateSingBoxConfigResponse represents the response for creating a SingBox config
type CreateSingBoxConfigResponse struct {
	Body models.SingBoxConfig
}

// ListSingBoxConfigsInput represents the input for listing SingBox configs
type ListSingBoxConfigsInput struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset int `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
}

// ListSingBoxConfigsResponse represents the response for listing SingBox configs
type ListSingBoxConfigsResponse struct {
	Body struct {
		Configs []models.SingBoxConfig `json:"configs"`
	}
}

// GetSingBoxConfigInput represents the input for getting a SingBox config by ID
type GetSingBoxConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"SingBox Configuration ID"`
}

// GetSingBoxConfigResponse represents the response for getting a SingBox config
type GetSingBoxConfigResponse struct {
	Body models.SingBoxConfig
}

// UpdateSingBoxConfigInput represents the input for updating a SingBox config
type UpdateSingBoxConfigInput struct {
	ConfigID string                `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"SingBox Configuration ID"`
	Body     models.SingBoxConfig `json:"singBoxConfig"`
}

// UpdateSingBoxConfigResponse represents the response for updating a SingBox config
type UpdateSingBoxConfigResponse struct {
	Body models.SingBoxConfig
}

// DeleteSingBoxConfigInput represents the input for deleting a SingBox config
type DeleteSingBoxConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"SingBox Configuration ID"`
}

// DeleteSingBoxConfigResponse represents the response for deleting a SingBox config
type DeleteSingBoxConfigResponse struct {
	Body struct {
		Message string `json:"message" example:"Configuration deleted successfully"`
	}
}

// GenerateSingBoxConfigInput represents the input for generating a SingBox config
type GenerateSingBoxConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"SingBox Configuration ID"`
}

// GenerateSingBoxConfigResponse represents the response for generating a SingBox config
type GenerateSingBoxConfigResponse struct {
	Body interface{} `json:"config" doc:"Generated SingBox configuration"`
}

// Xray Config Operations (similar structure)

// CreateXrayConfigInput represents the input for creating an Xray config
type CreateXrayConfigInput struct {
	Body models.XrayConfig `json:"xrayConfig"`
}

// CreateXrayConfigResponse represents the response for creating an Xray config
type CreateXrayConfigResponse struct {
	Body models.XrayConfig
}

// ListXrayConfigsInput represents the input for listing Xray configs
type ListXrayConfigsInput struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset int `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
}

// ListXrayConfigsResponse represents the response for listing Xray configs
type ListXrayConfigsResponse struct {
	Body struct {
		Configs []models.XrayConfig `json:"configs"`
	}
}

// GetXrayConfigInput represents the input for getting an Xray config by ID
type GetXrayConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Xray Configuration ID"`
}

// GetXrayConfigResponse represents the response for getting an Xray config
type GetXrayConfigResponse struct {
	Body models.XrayConfig
}

// UpdateXrayConfigInput represents the input for updating an Xray config
type UpdateXrayConfigInput struct {
	ConfigID string           `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Xray Configuration ID"`
	Body     models.XrayConfig `json:"xrayConfig"`
}

// UpdateXrayConfigResponse represents the response for updating an Xray config
type UpdateXrayConfigResponse struct {
	Body models.XrayConfig
}

// DeleteXrayConfigInput represents the input for deleting an Xray config
type DeleteXrayConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Xray Configuration ID"`
}

// DeleteXrayConfigResponse represents the response for deleting an Xray config
type DeleteXrayConfigResponse struct {
	Body struct {
		Message string `json:"message" example:"Configuration deleted successfully"`
	}
}

// GenerateXrayConfigInput represents the input for generating an Xray config
type GenerateXrayConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Xray Configuration ID"`
}

// GenerateXrayConfigResponse represents the response for generating an Xray config
type GenerateXrayConfigResponse struct {
	Body interface{} `json:"config" doc:"Generated Xray configuration"`
}

// HAProxy Config Operations

// CreateHAProxyConfigInput represents the input for creating an HAProxy config
type CreateHAProxyConfigInput struct {
	Body models.HAProxyConfig `json:"haproxyConfig"`
}

// CreateHAProxyConfigResponse represents the response for creating an HAProxy config
type CreateHAProxyConfigResponse struct {
	Body models.HAProxyConfig
}

// ListHAProxyConfigsInput represents the input for listing HAProxy configs
type ListHAProxyConfigsInput struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset int `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
}

// ListHAProxyConfigsResponse represents the response for listing HAProxy configs
type ListHAProxyConfigsResponse struct {
	Body struct {
		Configs []models.HAProxyConfig `json:"configs"`
	}
}

// GetHAProxyConfigInput represents the input for getting an HAProxy config by ID
type GetHAProxyConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"HAProxy Configuration ID"`
}

// GetHAProxyConfigResponse represents the response for getting an HAProxy config
type GetHAProxyConfigResponse struct {
	Body models.HAProxyConfig
}

// UpdateHAProxyConfigInput represents the input for updating an HAProxy config
type UpdateHAProxyConfigInput struct {
	ConfigID string               `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"HAProxy Configuration ID"`
	Body     models.HAProxyConfig `json:"haproxyConfig"`
}

// UpdateHAProxyConfigResponse represents the response for updating an HAProxy config
type UpdateHAProxyConfigResponse struct {
	Body models.HAProxyConfig
}

// DeleteHAProxyConfigInput represents the input for deleting an HAProxy config
type DeleteHAProxyConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"HAProxy Configuration ID"`
}

// DeleteHAProxyConfigResponse represents the response for deleting an HAProxy config
type DeleteHAProxyConfigResponse struct {
	Body struct {
		Message string `json:"message" example:"Configuration deleted successfully"`
	}
}

// GenerateHAProxyConfigInput represents the input for generating an HAProxy config
type GenerateHAProxyConfigInput struct {
	ConfigID string `path:"configId" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"HAProxy Configuration ID"`
}

// GenerateHAProxyConfigResponse represents the response for generating an HAProxy config
type GenerateHAProxyConfigResponse struct {
	Body interface{} `json:"config" doc:"Generated HAProxy configuration"`
}

// =============================================================================
// V2 API Types - Node Management
// =============================================================================

// Node filtering and pagination types
type NodeFilters struct {
	Status string   `query:"status" example:"active" doc:"Filter by node status (active, inactive, maintenance, error)"`
	Tags   []string `query:"tags" example:"production,us-east" doc:"Filter by tags (comma-separated)"`
	Search string   `query:"search" example:"server" doc:"Search in node name, hostname, or description"`
}

// ListNodesInput represents the input for listing nodes
type ListNodesInput struct {
	Limit  int      `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset int      `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
	Status string   `query:"status" example:"active" doc:"Filter by node status"`
	Tags   []string `query:"tags" example:"production,us-east" doc:"Filter by tags"`
	Search string   `query:"search" example:"server" doc:"Search in node name, hostname, or description"`
}

// ListNodesResponse represents the response for listing nodes
type ListNodesResponse struct {
	Body struct {
		Nodes  []models.NodeV2 `json:"nodes"`
		Total  int             `json:"total" example:"25" doc:"Total number of nodes matching the filter"`
		Limit  int             `json:"limit" example:"10" doc:"Number of items per page"`
		Offset int             `json:"offset" example:"0" doc:"Current offset"`
	}
}

// CreateNodeInput represents the input for creating a node
type CreateNodeInput struct {
	Body models.NodeCreateV2 `json:"node"`
}

// CreateNodeResponse represents the response for creating a node
type CreateNodeResponse struct {
	Body models.NodeV2
}

// GetNodeInput represents the input for getting a node by ID
type GetNodeInput struct {
	NodeID string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
}

// GetNodeResponse represents the response for getting a node
type GetNodeResponse struct {
	Body models.NodeV2
}

// UpdateNodeInput represents the input for updating a node
type UpdateNodeInput struct {
	NodeID string               `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	Body   models.NodeUpdateV2 `json:"node"`
}

// UpdateNodeResponse represents the response for updating a node
type UpdateNodeResponse struct {
	Body models.NodeV2
}

// DeleteNodeInput represents the input for deleting a node
type DeleteNodeInput struct {
	NodeID string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
}

// DeleteNodeResponse represents the response for deleting a node
type DeleteNodeResponse struct {
	Body struct {
		Message string `json:"message" example:"Node deleted successfully"`
	}
}

// GetNodeServicesInput represents the input for getting services on a node
type GetNodeServicesInput struct {
	NodeID      string   `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	Limit       int      `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset      int      `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
	ServiceType string   `query:"service_type" example:"xray" doc:"Filter by service type"`
	Status      string   `query:"status" example:"running" doc:"Filter by service status"`
	Tags        []string `query:"tags" example:"proxy,production" doc:"Filter by tags"`
}

// GetNodeServicesResponse represents the response for getting services on a node
type GetNodeServicesResponse struct {
	Body struct {
		Services []models.ServiceInstanceV2 `json:"services"`
		Total    int                        `json:"total" example:"5" doc:"Total number of services on this node"`
		Limit    int                        `json:"limit" example:"10" doc:"Number of items per page"`
		Offset   int                        `json:"offset" example:"0" doc:"Current offset"`
	}
}

// =============================================================================
// V2 API Types - Service Instance Management
// =============================================================================

// ListServicesInput represents the input for listing service instances
type ListServicesInput struct {
	Limit       int      `query:"limit" minimum:"1" maximum:"100" default:"10" example:"10" doc:"Number of items to return per page"`
	Offset      int      `query:"offset" minimum:"0" default:"0" example:"0" doc:"Offset for pagination"`
	NodeID      string   `query:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Filter by node ID"`
	ServiceType string   `query:"service_type" example:"xray" doc:"Filter by service type"`
	Status      string   `query:"status" example:"running" doc:"Filter by service status"`
	Tags        []string `query:"tags" example:"proxy,production" doc:"Filter by tags"`
	Port        int      `query:"port" example:"443" doc:"Filter by port number"`
	Protocol    string   `query:"protocol" example:"tcp" doc:"Filter by protocol"`
}

// ListServicesResponse represents the response for listing service instances
type ListServicesResponse struct {
	Body struct {
		Services []models.ServiceInstanceV2 `json:"services"`
		Total    int                        `json:"total" example:"15" doc:"Total number of services matching the filter"`
		Limit    int                        `json:"limit" example:"10" doc:"Number of items per page"`
		Offset   int                        `json:"offset" example:"0" doc:"Current offset"`
	}
}

// CreateServiceInput represents the input for creating a service instance
type CreateServiceInput struct {
	NodeID string                          `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	Body   models.ServiceInstanceCreateV2 `json:"service"`
}

// CreateServiceResponse represents the response for creating a service instance
type CreateServiceResponse struct {
	Body models.ServiceInstanceV2
}

// GetServiceInput represents the input for getting a service instance by ID
type GetServiceInput struct {
	NodeID    string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
}

// GetServiceResponse represents the response for getting a service instance
type GetServiceResponse struct {
	Body models.ServiceInstanceV2
}

// UpdateServiceInput represents the input for updating a service instance
type UpdateServiceInput struct {
	NodeID    string                          `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string                          `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
	Body      models.ServiceInstanceUpdateV2 `json:"service"`
}

// UpdateServiceResponse represents the response for updating a service instance
type UpdateServiceResponse struct {
	Body models.ServiceInstanceV2
}

// DeleteServiceInput represents the input for deleting a service instance
type DeleteServiceInput struct {
	NodeID    string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
}

// DeleteServiceResponse represents the response for deleting a service instance
type DeleteServiceResponse struct {
	Body struct {
		Message string `json:"message" example:"Service deleted successfully"`
	}
}

// GetServiceConfigInput represents the input for getting service configuration
type GetServiceConfigInput struct {
	NodeID    string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
}

// GetServiceConfigResponse represents the response for getting service configuration
type GetServiceConfigResponse struct {
	Body struct {
		ServiceType string                 `json:"service_type" example:"xray" doc:"Type of service"`
		Config      map[string]interface{} `json:"config" doc:"Service configuration"`
	}
}

// UpdateServiceConfigInput represents the input for updating service configuration
type UpdateServiceConfigInput struct {
	NodeID    string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
	Body      struct {
		Config map[string]interface{} `json:"config" doc:"Updated service configuration"`
	}
}

// UpdateServiceConfigResponse represents the response for updating service configuration
type UpdateServiceConfigResponse struct {
	Body struct {
		ServiceType string                 `json:"service_type" example:"xray" doc:"Type of service"`
		Config      map[string]interface{} `json:"config" doc:"Updated service configuration"`
		UpdatedAt   time.Time              `json:"updated_at" example:"2023-01-01T13:00:00Z" doc:"Timestamp of last update"`
	}
}

// GenerateServiceConfigInput represents the input for generating service configuration
type GenerateServiceConfigInput struct {
	NodeID    string `path:"nodeId" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Node ID"`
	ServiceID string `path:"serviceId" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" doc:"Service Instance ID"`
	Body      struct {
		Parameters map[string]interface{} `json:"parameters,omitempty" doc:"Configuration parameters for generation"`
	}
}

// GenerateServiceConfigResponse represents the response for generating service configuration
type GenerateServiceConfigResponse struct {
	Body struct {
		ServiceType string                 `json:"service_type" example:"xray" doc:"Type of service"`
		Config      map[string]interface{} `json:"config" doc:"Generated service configuration"`
		GeneratedAt time.Time              `json:"generated_at" example:"2023-01-01T13:00:00Z" doc:"Timestamp of generation"`
	}
}

// =============================================================================
// V2 API Error Types
// =============================================================================

// ApiErrorV2 represents a standardized V2 API error response
type ApiErrorV2 struct {
	Body struct {
		Error   string                 `json:"error" example:"Resource not found"`
		Code    string                 `json:"code,omitempty" example:"RESOURCE_NOT_FOUND"`
		Details map[string]interface{} `json:"details,omitempty" doc:"Additional error details"`
	}
}

// ValidationErrorV2 represents a validation error response
type ValidationErrorV2 struct {
	Body struct {
		Error  string `json:"error" example:"Validation failed"`
		Code   string `json:"code" example:"VALIDATION_ERROR"`
		Fields []struct {
			Field   string `json:"field" example:"name"`
			Message string `json:"message" example:"Name is required"`
		} `json:"fields"`
	}
}