package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/tools4net/ezfw/backend/internal/api/types"
	"github.com/tools4net/ezfw/backend/internal/models"
	"github.com/tools4net/ezfw/backend/internal/store"
)

// ServiceHandler handles service instance-related HTTP requests for V2 API
type ServiceHandler struct {
	store store.Store
}

// NewServiceHandler creates a new ServiceHandler instance
func NewServiceHandler(store store.Store) *ServiceHandler {
	return &ServiceHandler{
		store: store,
	}
}

// validateServiceName validates that the service name is not empty
func validateServiceName(name string) error {
	if strings.TrimSpace(name) == "" {
		return huma.Error400BadRequest("Service name cannot be empty")
	}
	return nil
}

// validateServiceType validates that the service type is supported
func validateServiceType(serviceType string) error {
	supportedTypes := map[string]bool{
		"xray":      true,
		"singbox":   true,
		"nginx":     true,
		"wireguard": true,
		"haproxy":   true,
	}

	if !supportedTypes[serviceType] {
		return huma.Error400BadRequest("Unsupported service type. Supported types: xray, singbox, nginx, wireguard, haproxy")
	}
	return nil
}

// ListServices retrieves all service instances with pagination and filtering
func (h *ServiceHandler) ListServices(ctx context.Context, input *types.ListServicesInput) (*types.ListServicesResponse, error) {
	// For global service listing, we need to implement this differently
	// Since the store interface only supports listing by nodeId, we'll need to modify the approach
	// For now, return an error indicating this needs implementation
	return nil, huma.Error501NotImplemented("Global service listing not yet implemented")
}

// CreateService creates a new service instance on a node
func (h *ServiceHandler) CreateService(ctx context.Context, input *types.CreateServiceInput) (*types.CreateServiceResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	serviceData := input.Body

	// Validate service name
	if err := validateServiceName(serviceData.Name); err != nil {
		return nil, err
	}

	// Validate service type
	if err := validateServiceType(serviceData.ServiceType); err != nil {
		return nil, err
	}

	// Check if a service with the same name already exists on this node
	existingServices, err := h.store.ListServiceInstances(ctx, input.NodeID, 100, 0) // Get up to 100 services to check
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check existing services", err)
	}

	for _, service := range existingServices {
		if service.Name == input.Body.Name {
			return nil, huma.Error409Conflict("A service with this name already exists on this node")
		}
		if service.Port == input.Body.Port && service.Protocol == input.Body.Protocol {
			return nil, huma.Error409Conflict("A service is already using this port and protocol combination on this node")
		}
	}

	// Create the service instance
	service, err := h.store.CreateServiceInstance(ctx, input.NodeID, &input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create service", err)
	}

	resp := &types.CreateServiceResponse{}
	resp.Body = *service

	return resp, nil
}

// GetService retrieves a specific service instance by ID
func (h *ServiceHandler) GetService(ctx context.Context, input *types.GetServiceInput) (*types.GetServiceResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Get the service instance
	service, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if service.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	resp := &types.GetServiceResponse{}
	resp.Body = *service

	return resp, nil
}

// UpdateService updates an existing service instance
func (h *ServiceHandler) UpdateService(ctx context.Context, input *types.UpdateServiceInput) (*types.UpdateServiceResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Check if the service exists
	existingService, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if existingService.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	// Validate service name if provided
	if input.Body.Name != nil {
		if err := validateServiceName(*input.Body.Name); err != nil {
			return nil, err
		}

		// Check if another service with the same name exists on this node (excluding current service)
		if *input.Body.Name != existingService.Name {
			existingServices, err := h.store.ListServiceInstances(ctx, input.NodeID, 100, 0) // Get up to 100 services to check
			if err != nil {
				return nil, huma.Error500InternalServerError("Failed to check existing services", err)
			}

			for _, service := range existingServices {
				if service.ID != input.ServiceID && service.Name == *input.Body.Name {
					return nil, huma.Error409Conflict("A service with this name already exists on this node")
				}
			}
		}
	}

	// Check for port/protocol conflicts if they are being updated
	if input.Body.Port != nil || input.Body.Protocol != nil {
		newPort := existingService.Port
		newProtocol := existingService.Protocol

		if input.Body.Port != nil {
			newPort = *input.Body.Port
		}
		if input.Body.Protocol != nil {
			newProtocol = *input.Body.Protocol
		}

		// Only check for conflicts if port or protocol actually changed
		if newPort != existingService.Port || newProtocol != existingService.Protocol {
			existingServices, err := h.store.ListServiceInstances(ctx, input.NodeID, 100, 0) // Get up to 100 services to check
			if err != nil {
				return nil, huma.Error500InternalServerError("Failed to check existing services", err)
			}

			for _, service := range existingServices {
				if service.ID != input.ServiceID && service.Port == newPort && service.Protocol == newProtocol {
					return nil, huma.Error409Conflict("A service is already using this port and protocol combination on this node")
				}
			}
		}
	}

	// Update the service instance
	updatedService, err := h.store.UpdateServiceInstance(ctx, input.NodeID, input.ServiceID, &input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update service", err)
	}

	resp := &types.UpdateServiceResponse{}
	resp.Body = *updatedService

	return resp, nil
}

// DeleteService deletes a service instance
func (h *ServiceHandler) DeleteService(ctx context.Context, input *types.DeleteServiceInput) (*types.DeleteServiceResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Check if the service exists
	existingService, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if existingService.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	// Delete the service instance
	err = h.store.DeleteServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete service", err)
	}

	resp := &types.DeleteServiceResponse{}
	resp.Body.Message = "Service deleted successfully"

	return resp, nil
}

// GetServiceConfig retrieves the configuration for a specific service instance
func (h *ServiceHandler) GetServiceConfig(ctx context.Context, input *types.GetServiceConfigInput) (*types.GetServiceConfigResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Get the service instance
	service, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if service.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	resp := &types.GetServiceConfigResponse{}
	resp.Body.ServiceType = service.ServiceType
	resp.Body.Config = service.Config

	return resp, nil
}

// UpdateServiceConfig updates the configuration for a specific service instance
func (h *ServiceHandler) UpdateServiceConfig(ctx context.Context, input *types.UpdateServiceConfigInput) (*types.UpdateServiceConfigResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Check if the service exists
	existingService, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if existingService.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	// Validate configuration based on service type
	if err := h.validateServiceConfig(existingService.ServiceType, input.Body.Config); err != nil {
		return nil, err
	}

	// Update the service configuration
	updateData := &models.ServiceInstanceUpdateV2{
		Config: input.Body.Config,
	}

	updatedService, err := h.store.UpdateServiceInstance(ctx, input.NodeID, input.ServiceID, updateData)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update service configuration", err)
	}

	resp := &types.UpdateServiceConfigResponse{}
	resp.Body.ServiceType = updatedService.ServiceType
	resp.Body.Config = updatedService.Config
	resp.Body.UpdatedAt = updatedService.UpdatedAt

	return resp, nil
}

// validateServiceConfig validates configuration based on service type
func (h *ServiceHandler) validateServiceConfig(serviceType string, config interface{}) error {
	if config == nil {
		return huma.Error400BadRequest("Configuration cannot be empty")
	}

	// Convert to map[string]interface{}
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return huma.Error400BadRequest("Configuration must be a valid JSON object")
	}

	// Type-specific validation can be added here
	switch serviceType {
	case "xray":
		return validateXrayConfig(configMap)
	case "singbox":
		return validateSingBoxConfig(configMap)
	case "nginx":
		return validateNginxConfig(configMap)
	case "haproxy":
		return validateHAProxyConfig(configMap)
	case "wireguard":
		return validateWireGuardConfig(configMap)
	default:
		return huma.Error400BadRequest("Unsupported service type for configuration validation")
	}
}

// validateXrayConfig validates Xray-specific configuration
func validateXrayConfig(config map[string]interface{}) error {
	// Basic Xray configuration validation
	if config == nil {
		return errors.New("xray configuration cannot be empty")
	}

	// Check for required fields
	if _, ok := config["inbounds"]; !ok {
		return errors.New("xray configuration must include 'inbounds'")
	}

	// Additional Xray-specific validation can be added here
	return nil
}

// validateSingBoxConfig validates SingBox-specific configuration
func validateSingBoxConfig(config map[string]interface{}) error {
	// Basic SingBox configuration validation
	if config == nil {
		return errors.New("singbox configuration cannot be empty")
	}

	// Check for required fields
	if _, ok := config["inbounds"]; !ok {
		return errors.New("singbox configuration must include 'inbounds'")
	}
	if _, ok := config["outbounds"]; !ok {
		return errors.New("singbox configuration must include 'outbounds'")
	}

	// Additional SingBox-specific validation can be added here
	return nil
}

// validateNginxConfig validates Nginx-specific configuration
func validateNginxConfig(config map[string]interface{}) error {
	// Basic Nginx configuration validation
	if config == nil {
		return errors.New("nginx configuration cannot be empty")
	}

	// Check for required fields
	if _, ok := config["server"]; !ok {
		return errors.New("nginx configuration must include 'server' block")
	}

	// Additional Nginx-specific validation can be added here
	return nil
}

// validateHAProxyConfig validates HAProxy-specific configuration
func validateHAProxyConfig(config map[string]interface{}) error {
	// Basic HAProxy configuration validation
	if config == nil {
		return errors.New("haproxy configuration cannot be empty")
	}

	// Check for required fields
	if _, ok := config["frontend"]; !ok {
		return errors.New("haproxy configuration must include 'frontend'")
	}
	if _, ok := config["backend"]; !ok {
		return errors.New("haproxy configuration must include 'backend'")
	}

	// Additional HAProxy-specific validation can be added here
	return nil
}

// validateWireGuardConfig validates WireGuard-specific configuration
func validateWireGuardConfig(config map[string]interface{}) error {
	// Basic WireGuard configuration validation
	if config == nil {
		return errors.New("wireguard configuration cannot be empty")
	}

	// Check for required fields
	if _, ok := config["interface"]; !ok {
		return errors.New("wireguard configuration must include 'interface'")
	}

	// Additional WireGuard-specific validation can be added here
	return nil
}

// GenerateServiceConfig generates a service configuration based on service type and parameters
func (h *ServiceHandler) GenerateServiceConfig(ctx context.Context, input *types.GenerateServiceConfigInput) (*types.GenerateServiceConfigResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Check if the service exists
	service, err := h.store.GetServiceInstance(ctx, input.NodeID, input.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Service not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve service", err)
	}

	// Verify the service belongs to the specified node
	if service.NodeID != input.NodeID {
		return nil, huma.Error404NotFound("Service not found on this node")
	}

	// Generate configuration based on service type
	generatedConfig, err := generateConfigByType(service.ServiceType, input.Body.Parameters)
	if err != nil {
		return nil, err
	}

	resp := &types.GenerateServiceConfigResponse{}
	resp.Body.ServiceType = service.ServiceType
	resp.Body.Config = generatedConfig
	resp.Body.GeneratedAt = time.Now()

	return resp, nil
}

// GenerateServiceConfig generates a new configuration for a service instance (helper function)
func GenerateServiceConfig(ctx context.Context, serviceType string, parameters map[string]interface{}) (map[string]interface{}, error) {
	return generateConfigByType(serviceType, parameters)
}

// generateConfigByType generates configuration based on service type
func generateConfigByType(serviceType string, parameters map[string]interface{}) (map[string]interface{}, error) {
	switch serviceType {
	case "xray":
		return generateXrayConfig(parameters)
	case "singbox":
		return generateSingBoxConfig(parameters)
	case "nginx":
		return generateNginxConfig(parameters)
	case "haproxy":
		return generateHAProxyConfig(parameters)
	case "wireguard":
		return generateWireGuardConfig(parameters)
	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}
}

// generateXrayConfig generates Xray configuration
func generateXrayConfig(parameters map[string]interface{}) (map[string]interface{}, error) {
	// Start with default Xray configuration
	config := map[string]interface{}{
		"log": map[string]interface{}{
			"loglevel": "warning",
		},
		"inbounds": []map[string]interface{}{
			{
				"port":     10808,
				"protocol": "vmess",
				"settings": map[string]interface{}{
					"clients": []map[string]interface{}{},
				},
			},
		},
		"outbounds": []map[string]interface{}{
			{
				"protocol": "freedom",
				"settings": map[string]interface{}{},
			},
		},
	}

	// Apply custom parameters
	for key, value := range parameters {
		config[key] = value
	}

	return config, nil
}

// generateSingBoxConfig generates SingBox configuration
func generateSingBoxConfig(parameters map[string]interface{}) (map[string]interface{}, error) {
	// Start with default SingBox configuration
	config := map[string]interface{}{
		"log": map[string]interface{}{
			"level": "warn",
		},
		"inbounds": []map[string]interface{}{
			{
				"type":        "mixed",
				"listen":      "0.0.0.0",
				"listen_port": 10808,
			},
		},
		"outbounds": []map[string]interface{}{
			{
				"type": "direct",
			},
		},
	}

	// Apply custom parameters
	for key, value := range parameters {
		config[key] = value
	}

	return config, nil
}

// generateNginxConfig generates Nginx configuration
func generateNginxConfig(parameters map[string]interface{}) (map[string]interface{}, error) {
	// Start with default Nginx configuration
	config := map[string]interface{}{
		"server": map[string]interface{}{
			"listen":      80,
			"server_name": "_",
			"location": map[string]interface{}{
				"/": map[string]interface{}{
					"proxy_pass": "http://backend",
					"proxy_set_header": map[string]interface{}{
						"Host":            "$host",
						"X-Real-IP":       "$remote_addr",
						"X-Forwarded-For": "$proxy_add_x_forwarded_for",
					},
				},
			},
		},
	}

	// Apply custom parameters
	for key, value := range parameters {
		config[key] = value
	}

	return config, nil
}

// generateHAProxyConfig generates HAProxy configuration
func generateHAProxyConfig(parameters map[string]interface{}) (map[string]interface{}, error) {
	config := map[string]interface{}{
		"global": map[string]interface{}{
			"daemon": true,
			"maxconn": 4096,
		},
		"defaults": map[string]interface{}{
			"mode": "http",
			"timeout": map[string]interface{}{
				"connect": "5000ms",
				"client": "50000ms",
				"server": "50000ms",
			},
		},
		"frontend": map[string]interface{}{
			"main": map[string]interface{}{
				"bind": "*:80",
				"default_backend": "servers",
			},
		},
		"backend": map[string]interface{}{
			"servers": map[string]interface{}{
				"balance": "roundrobin",
				"server": []string{"web1 127.0.0.1:8080 check"},
			},
		},
	}

	// Apply parameters
	for key, value := range parameters {
		config[key] = value
	}

	return config, nil
}

// generateWireGuardConfig generates WireGuard configuration
func generateWireGuardConfig(parameters map[string]interface{}) (map[string]interface{}, error) {
	config := map[string]interface{}{
		"interface": map[string]interface{}{
			"private_key": "PRIVATE_KEY_PLACEHOLDER",
			"address": "10.0.0.1/24",
			"listen_port": 51820,
		},
		"peers": []map[string]interface{}{},
	}

	// Apply parameters
	for key, value := range parameters {
		config[key] = value
	}

	return config, nil
}