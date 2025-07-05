package models

import (
	"encoding/json"
	"time"
)

// ServiceInstanceV2 represents a service running on a node
type ServiceInstanceV2 struct {
	ID          string                 `json:"id" gorm:"primaryKey" example:"service-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	NodeID      string                 `json:"node_id" gorm:"index" example:"node-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Name        string                 `json:"name" example:"Xray Proxy Service"`
	Description string                 `json:"description,omitempty" example:"Main Xray proxy service for US-East"`
	ServiceType string                 `json:"service_type" example:"xray"` // xray, singbox, nginx, wireguard, haproxy
	Status      string                 `json:"status" example:"running"`     // running, stopped, error, starting, stopping
	Port        int                    `json:"port" example:"443"`
	Protocol    string                 `json:"protocol" example:"tcp"`       // tcp, udp, both
	Config      map[string]interface{} `json:"config,omitempty"`             // Polymorphic configuration based on service type
	Tags        []string               `json:"tags,omitempty" example:"proxy,production,high-availability"`
	CreatedAt   time.Time              `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   time.Time              `json:"updated_at" example:"2023-01-01T13:00:00Z"`
	LastSeen    *time.Time             `json:"last_seen,omitempty" example:"2023-01-01T13:00:00Z"`

	// Relationships
	Node *NodeV2 `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// ServiceInstanceCreateV2 represents the data needed to create a new service instance
type ServiceInstanceCreateV2 struct {
	Name        string                 `json:"name" validate:"required" example:"Xray Proxy Service"`
	Description string                 `json:"description,omitempty" example:"Main Xray proxy service for US-East"`
	ServiceType string                 `json:"service_type" validate:"required,oneof=xray singbox nginx wireguard haproxy" example:"xray"`
	Port        int                    `json:"port" validate:"required,min=1,max=65535" example:"443"`
	Protocol    string                 `json:"protocol" validate:"required,oneof=tcp udp both" example:"tcp"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string               `json:"tags,omitempty" example:"proxy,production,high-availability"`
}

// ServiceInstanceUpdateV2 represents the data that can be updated for a service instance
type ServiceInstanceUpdateV2 struct {
	Name        *string                `json:"name,omitempty" example:"Updated Service Name"`
	Description *string                `json:"description,omitempty" example:"Updated description"`
	Port        *int                   `json:"port,omitempty" validate:"omitempty,min=1,max=65535" example:"8443"`
	Protocol    *string                `json:"protocol,omitempty" validate:"omitempty,oneof=tcp udp both" example:"udp"`
	Status      *string                `json:"status,omitempty" example:"stopped"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string               `json:"tags,omitempty" example:"proxy,production,updated"`
}

// ServiceInstanceFilters represents filters for listing service instances
type ServiceInstanceFilters struct {
	ServiceType string   `json:"service_type,omitempty" example:"xray"`
	Status      string   `json:"status,omitempty" example:"running"`
	Tags        []string `json:"tags,omitempty" example:"production,proxy"`
	Port        int      `json:"port,omitempty" example:"443"`
	Protocol    string   `json:"protocol,omitempty" example:"tcp"`
}

// GetXrayConfig extracts Xray configuration from the polymorphic config field
func (s *ServiceInstanceV2) GetXrayConfig() (*XrayConfig, error) {
	if s.ServiceType != "xray" {
		return nil, ErrInvalidServiceType
	}

	configBytes, err := json.Marshal(s.Config)
	if err != nil {
		return nil, err
	}

	var xrayConfig XrayConfig
	err = json.Unmarshal(configBytes, &xrayConfig)
	if err != nil {
		return nil, err
	}

	return &xrayConfig, nil
}

// GetSingBoxConfig extracts SingBox configuration from the polymorphic config field
func (s *ServiceInstanceV2) GetSingBoxConfig() (*SingBoxConfig, error) {
	if s.ServiceType != "singbox" {
		return nil, ErrInvalidServiceType
	}

	configBytes, err := json.Marshal(s.Config)
	if err != nil {
		return nil, err
	}

	var singboxConfig SingBoxConfig
	err = json.Unmarshal(configBytes, &singboxConfig)
	if err != nil {
		return nil, err
	}

	return &singboxConfig, nil
}

// SetXrayConfig sets Xray configuration in the polymorphic config field
func (s *ServiceInstanceV2) SetXrayConfig(config *XrayConfig) error {
	if s.ServiceType != "xray" {
		return ErrInvalidServiceType
	}

	configMap := make(map[string]interface{})
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configBytes, &configMap)
	if err != nil {
		return err
	}

	s.Config = configMap
	return nil
}

// SetSingBoxConfig sets SingBox configuration in the polymorphic config field
func (s *ServiceInstanceV2) SetSingBoxConfig(config *SingBoxConfig) error {
	if s.ServiceType != "singbox" {
		return ErrInvalidServiceType
	}

	configMap := make(map[string]interface{})
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configBytes, &configMap)
	if err != nil {
		return err
	}

	s.Config = configMap
	return nil
}