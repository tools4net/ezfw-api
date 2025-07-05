package store

import (
	"context"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// Store defines the interface for database operations.
type Store interface {
	// SingBox Configuration methods
	CreateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error
	GetSingBoxConfig(ctx context.Context, id string) (*models.SingBoxConfig, error)
	ListSingBoxConfigs(ctx context.Context, limit, offset int) ([]*models.SingBoxConfig, error)
	UpdateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error
	DeleteSingBoxConfig(ctx context.Context, id string) error
	// CountSingBoxConfigs(ctx context.Context) (int, error) // Optional: for pagination metadata

	// Xray Configuration methods
	CreateXrayConfig(ctx context.Context, config *models.XrayConfig) error
	GetXrayConfig(ctx context.Context, id string) (*models.XrayConfig, error)
	ListXrayConfigs(ctx context.Context, limit, offset int) ([]*models.XrayConfig, error)
	UpdateXrayConfig(ctx context.Context, config *models.XrayConfig) error
	DeleteXrayConfig(ctx context.Context, id string) error
	// CountXrayConfigs(ctx context.Context) (int, error) // Optional: for pagination metadata

	// HAProxy Configuration methods
	CreateHAProxyConfig(ctx context.Context, config *models.HAProxyConfig) error
	GetHAProxyConfig(ctx context.Context, id string) (*models.HAProxyConfig, error)
	ListHAProxyConfigs(ctx context.Context, limit, offset int) ([]*models.HAProxyConfig, error)
	UpdateHAProxyConfig(ctx context.Context, config *models.HAProxyConfig) error
	DeleteHAProxyConfig(ctx context.Context, id string) error
	// CountHAProxyConfigs(ctx context.Context) (int, error) // Optional: for pagination metadata

	// V2 Node management methods
	CreateNode(ctx context.Context, node *models.NodeCreateV2) (*models.NodeV2, error)
	GetNode(ctx context.Context, id string) (*models.NodeV2, error)
	ListNodes(ctx context.Context, filters models.NodeFilters, limit, offset int) ([]*models.NodeV2, error)
	UpdateNode(ctx context.Context, id string, updates *models.NodeUpdateV2) (*models.NodeV2, error)
	DeleteNode(ctx context.Context, id string) error

	// V2 Service instance management methods
	CreateServiceInstance(ctx context.Context, nodeId string, service *models.ServiceInstanceCreateV2) (*models.ServiceInstanceV2, error)
	GetServiceInstance(ctx context.Context, nodeId, serviceId string) (*models.ServiceInstanceV2, error)
	ListServiceInstances(ctx context.Context, nodeId string, limit, offset int) ([]*models.ServiceInstanceV2, error)
	UpdateServiceInstance(ctx context.Context, nodeId, serviceId string, updates *models.ServiceInstanceUpdateV2) (*models.ServiceInstanceV2, error)
	DeleteServiceInstance(ctx context.Context, nodeId, serviceId string) error

	// Agent token management methods
	CreateAgentToken(ctx context.Context, token *models.AgentTokenCreate) (*models.AgentToken, error)
	GetAgentToken(ctx context.Context, id string) (*models.AgentToken, error)
	GetAgentTokenByToken(ctx context.Context, token string) (*models.AgentToken, error)
	ListAgentTokens(ctx context.Context, filters models.AgentTokenFilters, limit, offset int) ([]*models.AgentToken, error)
	UpdateAgentToken(ctx context.Context, id string, updates *models.AgentTokenUpdate) (*models.AgentToken, error)
	DeleteAgentToken(ctx context.Context, id string) error
	RevokeAgentToken(ctx context.Context, id string) error
}
