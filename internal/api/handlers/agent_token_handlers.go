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

// AgentTokenHandlers contains handlers for agent token management
type AgentTokenHandlers struct {
	store  store.Store
	logger *zap.Logger
}

// NewAgentTokenHandlers creates a new instance of agent token handlers
func NewAgentTokenHandlers(store store.Store, logger *zap.Logger) *AgentTokenHandlers {
	return &AgentTokenHandlers{
		store:  store,
		logger: logger,
	}
}

// Agent Token Request/Response Models
type CreateAgentTokenRequest struct {
	Body struct {
		NodeID    string     `json:"node_id" validate:"required" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		Name      string     `json:"name" validate:"required" example:"Production Agent Token"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	} `json:"body"`
}

type CreateAgentTokenResponse struct {
	Body struct {
		ID        string     `json:"id" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		NodeID    string     `json:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		Token     string     `json:"token" example:"agt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
		Name      string     `json:"name" example:"Production Agent Token"`
		Status    string     `json:"status" example:"active"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
		CreatedAt time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	} `json:"body"`
}

type GetAgentTokenRequest struct {
	TokenID string `path:"token_id" validate:"required" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
}

type GetAgentTokenResponse struct {
	Body struct {
		ID        string     `json:"id" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		NodeID    string     `json:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		Name      string     `json:"name" example:"Production Agent Token"`
		Status    string     `json:"status" example:"active"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
		CreatedAt time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
		UpdatedAt time.Time  `json:"updated_at" example:"2023-01-01T13:00:00Z"`
		LastUsed  *time.Time `json:"last_used,omitempty" example:"2023-01-01T13:00:00Z"`
		Node      *models.NodeV2 `json:"node,omitempty"`
	} `json:"body"`
}

type ListAgentTokensRequest struct {
	NodeID string `query:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Status string `query:"status" example:"active"`
	Limit  int    `query:"limit" minimum:"1" maximum:"100" default:"20" example:"20"`
	Offset int    `query:"offset" minimum:"0" default:"0" example:"0"`
}

type ListAgentTokensResponse struct {
	Body struct {
		Tokens []AgentTokenSummary `json:"tokens"`
		Total  int                 `json:"total" example:"50"`
		Limit  int                 `json:"limit" example:"20"`
		Offset int                 `json:"offset" example:"0"`
	} `json:"body"`
}

type UpdateAgentTokenRequest struct {
	TokenID string `path:"token_id" validate:"required" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Body    struct {
		Name      *string    `json:"name,omitempty" example:"Updated Token Name"`
		Status    *string    `json:"status,omitempty" example:"revoked"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	} `json:"body"`
}

type UpdateAgentTokenResponse struct {
	Body struct {
		ID        string     `json:"id" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		NodeID    string     `json:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
		Name      string     `json:"name" example:"Updated Token Name"`
		Status    string     `json:"status" example:"revoked"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
		UpdatedAt time.Time  `json:"updated_at" example:"2023-01-01T13:00:00Z"`
	} `json:"body"`
}

type DeleteAgentTokenRequest struct {
	TokenID string `path:"token_id" validate:"required" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
}

type DeleteAgentTokenResponse struct {
	Body struct {
		Message string `json:"message" example:"Agent token deleted successfully"`
	} `json:"body"`
}

type RevokeAgentTokenRequest struct {
	TokenID string `path:"token_id" validate:"required" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
}

type RevokeAgentTokenResponse struct {
	Body struct {
		Message string `json:"message" example:"Agent token revoked successfully"`
	} `json:"body"`
}

// Supporting Models
type AgentTokenSummary struct {
	ID        string     `json:"id" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	NodeID    string     `json:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Name      string     `json:"name" example:"Production Agent Token"`
	Status    string     `json:"status" example:"active"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	LastUsed  *time.Time `json:"last_used,omitempty" example:"2023-01-01T13:00:00Z"`
	Node      *models.NodeV2 `json:"node,omitempty"`
}

// RegisterAgentTokenRoutes registers all agent token management routes
func (h *AgentTokenHandlers) RegisterAgentTokenRoutes(api huma.API, authMiddleware *auth.AuthMiddleware) {
	// Create agent token
	huma.Register(api, huma.Operation{
		OperationID: "create-agent-token",
		Method:      http.MethodPost,
		Path:        "/api/v2/agent-tokens",
		Summary:     "Create Agent Token",
		Description: "Create a new authentication token for a node agent",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.CreateAgentToken)

	// Get agent token
	huma.Register(api, huma.Operation{
		OperationID: "get-agent-token",
		Method:      http.MethodGet,
		Path:        "/api/v2/agent-tokens/{token_id}",
		Summary:     "Get Agent Token",
		Description: "Get details of a specific agent token",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.GetAgentToken)

	// List agent tokens
	huma.Register(api, huma.Operation{
		OperationID: "list-agent-tokens",
		Method:      http.MethodGet,
		Path:        "/api/v2/agent-tokens",
		Summary:     "List Agent Tokens",
		Description: "List agent tokens with optional filtering",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.ListAgentTokens)

	// Update agent token
	huma.Register(api, huma.Operation{
		OperationID: "update-agent-token",
		Method:      http.MethodPut,
		Path:        "/api/v2/agent-tokens/{token_id}",
		Summary:     "Update Agent Token",
		Description: "Update an existing agent token",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.UpdateAgentToken)

	// Revoke agent token
	huma.Register(api, huma.Operation{
		OperationID: "revoke-agent-token",
		Method:      http.MethodPost,
		Path:        "/api/v2/agent-tokens/{token_id}/revoke",
		Summary:     "Revoke Agent Token",
		Description: "Revoke an agent token (sets status to revoked)",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.RevokeAgentToken)

	// Delete agent token
	huma.Register(api, huma.Operation{
		OperationID: "delete-agent-token",
		Method:      http.MethodDelete,
		Path:        "/api/v2/agent-tokens/{token_id}",
		Summary:     "Delete Agent Token",
		Description: "Permanently delete an agent token",
		Tags:        []string{"Agent Tokens"},
		Security: []map[string][]string{
			{"BearerAuth": {}},
		},
		Middlewares: huma.Middlewares{authMiddleware.RequireUserAuth},
	}, h.DeleteAgentToken)
}

// CreateAgentToken creates a new agent token for a node
func (h *AgentTokenHandlers) CreateAgentToken(ctx context.Context, req *CreateAgentTokenRequest) (*CreateAgentTokenResponse, error) {
	// Get user context from authentication middleware
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		h.logger.Error("User context not found in create agent token request")
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	// Verify that the node exists and belongs to the user
	_, err := h.store.GetNode(ctx, req.Body.NodeID)
	if err != nil {
		h.logger.Error("Node not found for agent token creation",
			zap.String("node_id", req.Body.NodeID),
			zap.String("user_id", userCtx.UserID),
			zap.Error(err),
		)
		return nil, huma.Error404NotFound("Node not found")
	}

	// TODO: Add user ownership validation for the node
	// This would require adding user_id field to NodeV2 model

	// Generate a secure token
	token, err := auth.GenerateAgentToken()
	if err != nil {
		h.logger.Error("Failed to generate agent token", zap.Error(err))
		return nil, huma.Error500InternalServerError("Failed to generate token")
	}

	// Create the agent token
	tokenCreate := &models.AgentTokenCreate{
		NodeID:    req.Body.NodeID,
		Name:      req.Body.Name,
		ExpiresAt: req.Body.ExpiresAt,
	}

	createdToken, err := h.store.CreateAgentToken(ctx, tokenCreate)
	if err != nil {
		h.logger.Error("Failed to create agent token",
			zap.String("node_id", req.Body.NodeID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to create agent token")
	}

	// Update the created token with the generated token value
	// Note: This is a security consideration - we store the token but only return it once
	createdToken.Token = token

	h.logger.Info("Agent token created successfully",
		zap.String("token_id", createdToken.ID),
		zap.String("node_id", req.Body.NodeID),
		zap.String("user_id", userCtx.UserID),
	)

	response := &CreateAgentTokenResponse{}
	response.Body.ID = createdToken.ID
	response.Body.NodeID = createdToken.NodeID
	response.Body.Token = token // Only returned on creation
	response.Body.Name = createdToken.Name
	response.Body.Status = createdToken.Status
	response.Body.ExpiresAt = createdToken.ExpiresAt
	response.Body.CreatedAt = createdToken.CreatedAt

	return response, nil
}

// GetAgentToken retrieves a specific agent token
func (h *AgentTokenHandlers) GetAgentToken(ctx context.Context, req *GetAgentTokenRequest) (*GetAgentTokenResponse, error) {
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	token, err := h.store.GetAgentToken(ctx, req.TokenID)
	if err != nil {
		h.logger.Error("Agent token not found",
			zap.String("token_id", req.TokenID),
			zap.Error(err),
		)
		return nil, huma.Error404NotFound("Agent token not found")
	}

	// TODO: Verify user ownership of the token's node

	response := &GetAgentTokenResponse{}
	response.Body.ID = token.ID
	response.Body.NodeID = token.NodeID
	response.Body.Name = token.Name
	response.Body.Status = token.Status
	response.Body.ExpiresAt = token.ExpiresAt
	response.Body.CreatedAt = token.CreatedAt
	response.Body.UpdatedAt = token.UpdatedAt
	response.Body.LastUsed = token.LastUsed
	response.Body.Node = token.Node

	return response, nil
}

// ListAgentTokens lists agent tokens with optional filtering
func (h *AgentTokenHandlers) ListAgentTokens(ctx context.Context, req *ListAgentTokensRequest) (*ListAgentTokensResponse, error) {
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	// Set default values
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	filters := models.AgentTokenFilters{
		NodeID: req.NodeID,
		Status: req.Status,
	}

	// TODO: Add user ownership filtering

	tokens, err := h.store.ListAgentTokens(ctx, filters, req.Limit, req.Offset)
	if err != nil {
		h.logger.Error("Failed to list agent tokens", zap.Error(err))
		return nil, huma.Error500InternalServerError("Failed to retrieve agent tokens")
	}

	// Convert to summary format (without sensitive token data)
	tokenSummaries := make([]AgentTokenSummary, 0, len(tokens))
	for _, token := range tokens {
		summary := AgentTokenSummary{
			ID:        token.ID,
			NodeID:    token.NodeID,
			Name:      token.Name,
			Status:    token.Status,
			ExpiresAt: token.ExpiresAt,
			CreatedAt: token.CreatedAt,
			LastUsed:  token.LastUsed,
			Node:      token.Node,
		}
		tokenSummaries = append(tokenSummaries, summary)
	}

	response := &ListAgentTokensResponse{}
	response.Body.Tokens = tokenSummaries
	response.Body.Total = len(tokenSummaries) // TODO: Implement proper count
	response.Body.Limit = req.Limit
	response.Body.Offset = req.Offset

	return response, nil
}

// UpdateAgentToken updates an existing agent token
func (h *AgentTokenHandlers) UpdateAgentToken(ctx context.Context, req *UpdateAgentTokenRequest) (*UpdateAgentTokenResponse, error) {
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	// Verify token exists and user has access
	_, err := h.store.GetAgentToken(ctx, req.TokenID)
	if err != nil {
		return nil, huma.Error404NotFound("Agent token not found")
	}

	// TODO: Verify user ownership

	update := &models.AgentTokenUpdate{
		Name:      req.Body.Name,
		Status:    req.Body.Status,
		ExpiresAt: req.Body.ExpiresAt,
	}

	updatedToken, err := h.store.UpdateAgentToken(ctx, req.TokenID, update)
	if err != nil {
		h.logger.Error("Failed to update agent token",
			zap.String("token_id", req.TokenID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to update agent token")
	}

	response := &UpdateAgentTokenResponse{}
	response.Body.ID = updatedToken.ID
	response.Body.NodeID = updatedToken.NodeID
	response.Body.Name = updatedToken.Name
	response.Body.Status = updatedToken.Status
	response.Body.ExpiresAt = updatedToken.ExpiresAt
	response.Body.UpdatedAt = updatedToken.UpdatedAt

	return response, nil
}

// RevokeAgentToken revokes an agent token
func (h *AgentTokenHandlers) RevokeAgentToken(ctx context.Context, req *RevokeAgentTokenRequest) (*RevokeAgentTokenResponse, error) {
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	// Verify token exists and user has access
	_, err := h.store.GetAgentToken(ctx, req.TokenID)
	if err != nil {
		return nil, huma.Error404NotFound("Agent token not found")
	}

	// TODO: Verify user ownership

	err = h.store.RevokeAgentToken(ctx, req.TokenID)
	if err != nil {
		h.logger.Error("Failed to revoke agent token",
			zap.String("token_id", req.TokenID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to revoke agent token")
	}

	h.logger.Info("Agent token revoked",
		zap.String("token_id", req.TokenID),
		zap.String("user_id", userCtx.UserID),
	)

	response := &RevokeAgentTokenResponse{}
	response.Body.Message = "Agent token revoked successfully"

	return response, nil
}

// DeleteAgentToken permanently deletes an agent token
func (h *AgentTokenHandlers) DeleteAgentToken(ctx context.Context, req *DeleteAgentTokenRequest) (*DeleteAgentTokenResponse, error) {
	userCtx := auth.GetUserContext(ctx)
	if userCtx == nil {
		return nil, huma.Error400BadRequest("Invalid user context")
	}

	// Verify token exists and user has access
	_, err := h.store.GetAgentToken(ctx, req.TokenID)
	if err != nil {
		return nil, huma.Error404NotFound("Agent token not found")
	}

	// TODO: Verify user ownership

	err = h.store.DeleteAgentToken(ctx, req.TokenID)
	if err != nil {
		h.logger.Error("Failed to delete agent token",
			zap.String("token_id", req.TokenID),
			zap.Error(err),
		)
		return nil, huma.Error500InternalServerError("Failed to delete agent token")
	}

	h.logger.Info("Agent token deleted",
		zap.String("token_id", req.TokenID),
		zap.String("user_id", userCtx.UserID),
	)

	response := &DeleteAgentTokenResponse{}
	response.Body.Message = "Agent token deleted successfully"

	return response, nil
}