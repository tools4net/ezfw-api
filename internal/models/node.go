package models

import "time"

// NodeV2 represents a physical or virtual server in the V2 architecture
type NodeV2 struct {
	ID          string    `json:"id" gorm:"primaryKey" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Name        string    `json:"name" gorm:"uniqueIndex" example:"Production Server 1"`
	Description string    `json:"description,omitempty" example:"Main production server in US-East"`
	Hostname    string    `json:"hostname" example:"server1.example.com"`
	IPAddress   string    `json:"ip_address" example:"192.168.1.100"`
	Port        int       `json:"port" example:"22"`
	Status      string    `json:"status" example:"active"` // active, inactive, maintenance, error
	OSInfo      *OSInfo   `json:"os_info,omitempty"`
	AgentInfo   *AgentInfo `json:"agent_info,omitempty"`
	Tags        []string  `json:"tags,omitempty" example:"production,us-east,high-performance"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T13:00:00Z"`
	LastSeen    *time.Time `json:"last_seen,omitempty" example:"2023-01-01T13:00:00Z"`
}

// NodeCreateV2 represents the data needed to create a new node
type NodeCreateV2 struct {
	Name        string   `json:"name" validate:"required" example:"Production Server 1"`
	Description string   `json:"description,omitempty" example:"Main production server in US-East"`
	Hostname    string   `json:"hostname" validate:"required" example:"server1.example.com"`
	IPAddress   string   `json:"ip_address" validate:"required,ip" example:"192.168.1.100"`
	Port        int      `json:"port" validate:"required,min=1,max=65535" example:"22"`
	Tags        []string `json:"tags,omitempty" example:"production,us-east,high-performance"`
}

// NodeUpdateV2 represents the data that can be updated for a node
type NodeUpdateV2 struct {
	Name        *string  `json:"name,omitempty" example:"Updated Server Name"`
	Description *string  `json:"description,omitempty" example:"Updated description"`
	Hostname    *string  `json:"hostname,omitempty" example:"updated-server.example.com"`
	IPAddress   *string  `json:"ip_address,omitempty" validate:"omitempty,ip" example:"192.168.1.101"`
	Port        *int     `json:"port,omitempty" validate:"omitempty,min=1,max=65535" example:"2222"`
	Status      *string  `json:"status,omitempty" example:"maintenance"`
	Tags        []string `json:"tags,omitempty" example:"production,us-east,updated"`
}

// OSInfo contains operating system information for a node
type OSInfo struct {
	Name         string `json:"name" example:"Ubuntu"`
	Version      string `json:"version" example:"22.04 LTS"`
	Architecture string `json:"architecture" example:"x86_64"`
	Kernel       string `json:"kernel" example:"5.15.0-56-generic"`
}

// AgentInfo contains information about the agent running on the node
type AgentInfo struct {
	Version     string     `json:"version" example:"v1.2.3"`
	Status      string     `json:"status" example:"connected"` // connected, disconnected, error
	LastContact *time.Time `json:"last_contact,omitempty" example:"2023-01-01T13:00:00Z"`
	TokenID     string     `json:"token_id,omitempty" example:"token-xxxxxxxx"`
}

// NodeFilters represents filters for listing nodes
type NodeFilters struct {
	Status   string   `json:"status,omitempty" example:"active"`
	Tags     []string `json:"tags,omitempty" example:"production,us-east"`
	Hostname string   `json:"hostname,omitempty" example:"server1.example.com"`
}