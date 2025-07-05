package models

import (
	"time"
)

// AgentToken represents an authentication token for node agents
type AgentToken struct {
	ID        string    `json:"id" gorm:"primaryKey" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	NodeID    string    `json:"node_id" gorm:"index" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Token     string    `json:"token" gorm:"uniqueIndex" example:"agt_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`
	Name      string    `json:"name" example:"Production Agent Token"`
	Status    string    `json:"status" example:"active"` // active, revoked, expired
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T13:00:00Z"`
	LastUsed  *time.Time `json:"last_used,omitempty" example:"2023-01-01T13:00:00Z"`

	// Relationships
	Node *NodeV2 `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// AgentTokenCreate represents the data needed to create a new agent token
type AgentTokenCreate struct {
	NodeID    string     `json:"node_id" validate:"required" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Name      string     `json:"name" validate:"required" example:"Production Agent Token"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// AgentTokenUpdate represents the data that can be updated for an agent token
type AgentTokenUpdate struct {
	Name      *string    `json:"name,omitempty" example:"Updated Token Name"`
	Status    *string    `json:"status,omitempty" example:"revoked"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" example:"2024-12-31T23:59:59Z"`
}

// AgentTokenFilters represents filters for listing agent tokens
type AgentTokenFilters struct {
	NodeID string `json:"node_id,omitempty" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Status string `json:"status,omitempty" example:"active"`
}

// UserContext represents the authenticated user context from Clerk JWT
type UserContext struct {
	UserID       string   `json:"user_id" example:"user_xxxxxxxxxxxxxxxxxxxxxxxx"`
	Email        string   `json:"email" example:"user@example.com"`
	FirstName    string   `json:"first_name,omitempty" example:"John"`
	LastName     string   `json:"last_name,omitempty" example:"Doe"`
	Roles        []string `json:"roles,omitempty" example:"admin,user"`
	Organization string   `json:"organization,omitempty" example:"org_xxxxxxxxxxxxxxxxxxxxxxxx"`
}

// AgentContext represents the authenticated agent context
type AgentContext struct {
	TokenID string  `json:"token_id" example:"token-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	NodeID  string  `json:"node_id" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Node    *NodeV2 `json:"node,omitempty"`
}

// AuthContext represents the combined authentication context
type AuthContext struct {
	User  *UserContext  `json:"user,omitempty"`
	Agent *AgentContext `json:"agent,omitempty"`
	Type  string        `json:"type" example:"user"` // user, agent
}