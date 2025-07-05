package handlers

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tools4net/ezfw/backend/internal/api/types"
	"github.com/tools4net/ezfw/backend/internal/models"
	"github.com/tools4net/ezfw/backend/internal/store"
)

// ConfigHandler handles configuration-related HTTP requests
type ConfigHandler struct {
	store store.Store
}

// NewConfigHandler creates a new ConfigHandler instance
func NewConfigHandler(store store.Store) *ConfigHandler {
	return &ConfigHandler{
		store: store,
	}
}

// validateConfigName validates that the configuration name is not empty
func validateConfigName(name, configType string) error {
	if strings.TrimSpace(name) == "" {
		return huma.Error400BadRequest(configType + " name cannot be empty")
	}
	return nil
}

// Health check handler
func (h *ConfigHandler) HealthCheck(ctx context.Context, input *struct{}) (*types.HealthResponse, error) {
	resp := &types.HealthResponse{}
	resp.Body.Status = "API v2 is healthy"
	return resp, nil
}

// SingBox Config Handlers

// CreateSingBoxConfig creates a new SingBox configuration
func (h *ConfigHandler) CreateSingBoxConfig(ctx context.Context, input *types.CreateSingBoxConfigInput) (*types.CreateSingBoxConfigResponse, error) {
	config := input.Body

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "SingBox Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.CreateSingBoxConfig(ctx, &config); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create configuration: " + err.Error())
	}

	resp := &types.CreateSingBoxConfigResponse{}
	resp.Body = config
	return resp, nil
}

// ListSingBoxConfigs retrieves a paginated list of SingBox configurations
func (h *ConfigHandler) ListSingBoxConfigs(ctx context.Context, input *types.ListSingBoxConfigsInput) (*types.ListSingBoxConfigsResponse, error) {
	configs, err := h.store.ListSingBoxConfigs(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve SingBox configurations", err)
	}

	// Convert []*models.SingBoxConfig to []models.SingBoxConfig
	configList := make([]models.SingBoxConfig, len(configs))
	for i, config := range configs {
		configList[i] = *config
	}

	resp := &types.ListSingBoxConfigsResponse{}
	resp.Body.Configs = configList
	return resp, nil
}

// GetSingBoxConfig retrieves a specific SingBox configuration by its ID
func (h *ConfigHandler) GetSingBoxConfig(ctx context.Context, input *types.GetSingBoxConfigInput) (*types.GetSingBoxConfigResponse, error) {
	config, err := h.store.GetSingBoxConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("SingBox configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	resp := &types.GetSingBoxConfigResponse{}
	resp.Body = *config
	return resp, nil
}

// UpdateSingBoxConfig updates an existing SingBox configuration
func (h *ConfigHandler) UpdateSingBoxConfig(ctx context.Context, input *types.UpdateSingBoxConfigInput) (*types.UpdateSingBoxConfigResponse, error) {
	config := input.Body
	config.ID = input.ConfigID

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "SingBox Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.UpdateSingBoxConfig(ctx, &config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("SingBox configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to update configuration: " + err.Error())
	}

	resp := &types.UpdateSingBoxConfigResponse{}
	resp.Body = config
	return resp, nil
}

// DeleteSingBoxConfig deletes a SingBox configuration
func (h *ConfigHandler) DeleteSingBoxConfig(ctx context.Context, input *types.DeleteSingBoxConfigInput) (*types.DeleteSingBoxConfigResponse, error) {
	if err := h.store.DeleteSingBoxConfig(ctx, input.ConfigID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("SingBox configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete configuration: " + err.Error())
	}

	resp := &types.DeleteSingBoxConfigResponse{}
	resp.Body.Message = "Configuration deleted successfully"
	return resp, nil
}

// GenerateSingBoxConfig generates a SingBox configuration
func (h *ConfigHandler) GenerateSingBoxConfig(ctx context.Context, input *types.GenerateSingBoxConfigInput) (*types.GenerateSingBoxConfigResponse, error) {
	config, err := h.store.GetSingBoxConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("SingBox configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	// Generate the actual configuration (this would be implemented based on your business logic)
	generatedConfig := map[string]interface{}{
		"log":       config.Log,
		"dns":       config.DNS,
		"inbounds":  config.Inbounds,
		"outbounds": config.Outbounds,
		"route":     config.Route,
	}

	resp := &types.GenerateSingBoxConfigResponse{}
	resp.Body = generatedConfig
	return resp, nil
}

// Xray Config Handlers

// CreateXrayConfig creates a new Xray configuration
func (h *ConfigHandler) CreateXrayConfig(ctx context.Context, input *types.CreateXrayConfigInput) (*types.CreateXrayConfigResponse, error) {
	config := input.Body

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "Xray Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.CreateXrayConfig(ctx, &config); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create configuration: " + err.Error())
	}

	resp := &types.CreateXrayConfigResponse{}
	resp.Body = config
	return resp, nil
}

// ListXrayConfigs retrieves a paginated list of Xray configurations
func (h *ConfigHandler) ListXrayConfigs(ctx context.Context, input *types.ListXrayConfigsInput) (*types.ListXrayConfigsResponse, error) {
	configs, err := h.store.ListXrayConfigs(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve Xray configurations", err)
	}

	// Convert []*models.XrayConfig to []models.XrayConfig
	configList := make([]models.XrayConfig, len(configs))
	for i, config := range configs {
		configList[i] = *config
	}

	resp := &types.ListXrayConfigsResponse{}
	resp.Body.Configs = configList
	return resp, nil
}

// GetXrayConfig retrieves a specific Xray configuration by its ID
func (h *ConfigHandler) GetXrayConfig(ctx context.Context, input *types.GetXrayConfigInput) (*types.GetXrayConfigResponse, error) {
	config, err := h.store.GetXrayConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Xray configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	resp := &types.GetXrayConfigResponse{}
	resp.Body = *config
	return resp, nil
}

// UpdateXrayConfig updates an existing Xray configuration
func (h *ConfigHandler) UpdateXrayConfig(ctx context.Context, input *types.UpdateXrayConfigInput) (*types.UpdateXrayConfigResponse, error) {
	config := input.Body
	config.ID = input.ConfigID

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "Xray Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.UpdateXrayConfig(ctx, &config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Xray configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to update configuration: " + err.Error())
	}

	resp := &types.UpdateXrayConfigResponse{}
	resp.Body = config
	return resp, nil
}

// DeleteXrayConfig deletes an Xray configuration
func (h *ConfigHandler) DeleteXrayConfig(ctx context.Context, input *types.DeleteXrayConfigInput) (*types.DeleteXrayConfigResponse, error) {
	if err := h.store.DeleteXrayConfig(ctx, input.ConfigID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Xray configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete configuration: " + err.Error())
	}

	resp := &types.DeleteXrayConfigResponse{}
	resp.Body.Message = "Configuration deleted successfully"
	return resp, nil
}

// GenerateXrayConfig generates an Xray configuration
func (h *ConfigHandler) GenerateXrayConfig(ctx context.Context, input *types.GenerateXrayConfigInput) (*types.GenerateXrayConfigResponse, error) {
	config, err := h.store.GetXrayConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Xray configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	// Generate the actual configuration (this would be implemented based on your business logic)
	generatedConfig := map[string]interface{}{
		"log":       config.Log,
		"dns":       config.DNS,
		"inbounds":  config.Inbounds,
		"outbounds": config.Outbounds,
		"routing":   config.Routing,
	}

	resp := &types.GenerateXrayConfigResponse{}
	resp.Body = generatedConfig
	return resp, nil
}

// HAProxy Config Handlers

// CreateHAProxyConfig creates a new HAProxy configuration
func (h *ConfigHandler) CreateHAProxyConfig(ctx context.Context, input *types.CreateHAProxyConfigInput) (*types.CreateHAProxyConfigResponse, error) {
	config := input.Body

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "HAProxy Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.CreateHAProxyConfig(ctx, &config); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create configuration: " + err.Error())
	}

	resp := &types.CreateHAProxyConfigResponse{}
	resp.Body = config
	return resp, nil
}

// ListHAProxyConfigs retrieves a paginated list of HAProxy configurations
func (h *ConfigHandler) ListHAProxyConfigs(ctx context.Context, input *types.ListHAProxyConfigsInput) (*types.ListHAProxyConfigsResponse, error) {
	configs, err := h.store.ListHAProxyConfigs(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve HAProxy configurations", err)
	}

	// Convert []*models.HAProxyConfig to []models.HAProxyConfig
	configList := make([]models.HAProxyConfig, len(configs))
	for i, config := range configs {
		configList[i] = *config
	}

	resp := &types.ListHAProxyConfigsResponse{}
	resp.Body.Configs = configList
	return resp, nil
}

// GetHAProxyConfig retrieves a specific HAProxy configuration by its ID
func (h *ConfigHandler) GetHAProxyConfig(ctx context.Context, input *types.GetHAProxyConfigInput) (*types.GetHAProxyConfigResponse, error) {
	config, err := h.store.GetHAProxyConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("HAProxy configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	resp := &types.GetHAProxyConfigResponse{}
	resp.Body = *config
	return resp, nil
}

// UpdateHAProxyConfig updates an existing HAProxy configuration
func (h *ConfigHandler) UpdateHAProxyConfig(ctx context.Context, input *types.UpdateHAProxyConfigInput) (*types.UpdateHAProxyConfigResponse, error) {
	config := input.Body
	config.ID = input.ConfigID

	// Validate that the configuration name is not empty
	if err := validateConfigName(config.Name, "HAProxy Configuration"); err != nil {
		return nil, err
	}

	if err := h.store.UpdateHAProxyConfig(ctx, &config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("HAProxy configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to update configuration: " + err.Error())
	}

	resp := &types.UpdateHAProxyConfigResponse{}
	resp.Body = config
	return resp, nil
}

// DeleteHAProxyConfig deletes an HAProxy configuration
func (h *ConfigHandler) DeleteHAProxyConfig(ctx context.Context, input *types.DeleteHAProxyConfigInput) (*types.DeleteHAProxyConfigResponse, error) {
	if err := h.store.DeleteHAProxyConfig(ctx, input.ConfigID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("HAProxy configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete configuration: " + err.Error())
	}

	resp := &types.DeleteHAProxyConfigResponse{}
	resp.Body.Message = "Configuration deleted successfully"
	return resp, nil
}

// GenerateHAProxyConfig generates an HAProxy configuration
func (h *ConfigHandler) GenerateHAProxyConfig(ctx context.Context, input *types.GenerateHAProxyConfigInput) (*types.GenerateHAProxyConfigResponse, error) {
	config, err := h.store.GetHAProxyConfig(ctx, input.ConfigID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("HAProxy configuration not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get configuration: " + err.Error())
	}

	// Generate the actual configuration (this would be implemented based on your business logic)
	generatedConfig := map[string]interface{}{
		"global":    config.Global,
		"defaults":  config.Defaults,
		"frontends": config.Frontends,
		"backends":  config.Backends,
		"listens":   config.Listens,
		"stats":     config.Stats,
	}

	resp := &types.GenerateHAProxyConfigResponse{}
	resp.Body = generatedConfig
	return resp, nil
}