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

// NodeHandler handles node-related HTTP requests for V2 API
type NodeHandler struct {
	store store.Store
}

// NewNodeHandler creates a new NodeHandler instance
func NewNodeHandler(store store.Store) *NodeHandler {
	return &NodeHandler{
		store: store,
	}
}

// validateNodeName validates that the node name is not empty
func validateNodeName(name string) error {
	if strings.TrimSpace(name) == "" {
		return huma.Error400BadRequest("Node name cannot be empty")
	}
	return nil
}

// ListNodes retrieves all nodes with pagination and filtering
func (h *NodeHandler) ListNodes(ctx context.Context, input *types.ListNodesInput) (*types.ListNodesResponse, error) {
	// Apply filters and pagination
	filters := models.NodeFilters{
		Status:   input.Status,
		Tags:     input.Tags,
		Hostname: input.Search,
	}

	nodes, err := h.store.ListNodes(ctx, filters, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve nodes", err)
	}

	// Convert []*models.NodeV2 to []models.NodeV2
	nodeValues := make([]models.NodeV2, len(nodes))
	for i, node := range nodes {
		nodeValues[i] = *node
	}

	resp := &types.ListNodesResponse{}
	resp.Body.Nodes = nodeValues
	resp.Body.Total = len(nodes)
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset

	return resp, nil
}

// CreateNode creates a new node
func (h *NodeHandler) CreateNode(ctx context.Context, input *types.CreateNodeInput) (*types.CreateNodeResponse, error) {
	nodeData := input.Body

	// Validate that the node name is not empty
	if err := validateNodeName(nodeData.Name); err != nil {
		return nil, err
	}

	// Check if node with same name exists by listing with filters
	nameFilters := models.NodeFilters{Hostname: nodeData.Name}
	existingNodes, err := h.store.ListNodes(ctx, nameFilters, 1, 0)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check existing node", err)
	}
	if len(existingNodes) > 0 {
		return nil, huma.Error409Conflict("A node with this name already exists")
	}

	// Create the node
	node, err := h.store.CreateNode(ctx, &nodeData)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create node", err)
	}

	resp := &types.CreateNodeResponse{}
	resp.Body = *node

	return resp, nil
}

// GetNode retrieves a specific node by ID
func (h *NodeHandler) GetNode(ctx context.Context, input *types.GetNodeInput) (*types.GetNodeResponse, error) {
	node, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	resp := &types.GetNodeResponse{}
	resp.Body = *node

	return resp, nil
}

// UpdateNode updates an existing node
func (h *NodeHandler) UpdateNode(ctx context.Context, input *types.UpdateNodeInput) (*types.UpdateNodeResponse, error) {
	// Check if the node exists
	existingNode, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Validate node name if provided
	if input.Body.Name != nil {
		if err := validateNodeName(*input.Body.Name); err != nil {
			return nil, err
		}

		// Check if another node with the same name exists (excluding current node)
		if *input.Body.Name != existingNode.Name {
			nameFilters := models.NodeFilters{Hostname: *input.Body.Name}
			existingNodes, err := h.store.ListNodes(ctx, nameFilters, 1, 0)
			if err != nil {
				return nil, huma.Error500InternalServerError("Failed to check existing node", err)
			}
			if len(existingNodes) > 0 && existingNodes[0].ID != input.NodeID {
				return nil, huma.Error409Conflict("A node with this name already exists")
			}
		}
	}

	// Update the node
	updatedNode, err := h.store.UpdateNode(ctx, input.NodeID, &input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update node", err)
	}

	resp := &types.UpdateNodeResponse{}
	resp.Body = *updatedNode

	return resp, nil
}

// DeleteNode deletes a node
func (h *NodeHandler) DeleteNode(ctx context.Context, input *types.DeleteNodeInput) (*types.DeleteNodeResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Check if there are any service instances running on this node
	services, err := h.store.ListServiceInstances(ctx, input.NodeID, 1, 0) // Just check if any services exist
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check node services", err)
	}
	if len(services) > 0 {
		return nil, huma.Error409Conflict("Cannot delete node with active service instances. Please remove all services first.")
	}

	// Delete the node
	err = h.store.DeleteNode(ctx, input.NodeID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to delete node", err)
	}

	resp := &types.DeleteNodeResponse{}
	resp.Body.Message = "Node deleted successfully"

	return resp, nil
}

// GetNodeServices retrieves all service instances for a specific node
func (h *NodeHandler) GetNodeServices(ctx context.Context, input *types.GetNodeServicesInput) (*types.GetNodeServicesResponse, error) {
	// Check if the node exists
	_, err := h.store.GetNode(ctx, input.NodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("Node not found")
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve node", err)
	}

	// Get services for the node
	services, err := h.store.ListServiceInstances(ctx, input.NodeID, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve node services", err)
	}

	// Convert []*models.ServiceInstanceV2 to []models.ServiceInstanceV2
	serviceValues := make([]models.ServiceInstanceV2, len(services))
	for i, service := range services {
		serviceValues[i] = *service
	}

	resp := &types.GetNodeServicesResponse{}
	resp.Body.Services = serviceValues
	resp.Body.Total = len(services)
	resp.Body.Limit = input.Limit
	resp.Body.Offset = input.Offset

	return resp, nil
}