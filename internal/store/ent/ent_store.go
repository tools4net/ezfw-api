package ent

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"github.com/tools4net/ezfw/backend/ent"
	"github.com/tools4net/ezfw/backend/ent/agenttoken"
	"github.com/tools4net/ezfw/backend/ent/haproxyconfig"
	"github.com/tools4net/ezfw/backend/ent/node"
	"github.com/tools4net/ezfw/backend/ent/serviceinstance"
	"github.com/tools4net/ezfw/backend/ent/singboxconfig"
	"github.com/tools4net/ezfw/backend/ent/xrayconfig"
	"github.com/tools4net/ezfw/backend/internal/models"
	"github.com/tools4net/ezfw/backend/internal/store"
)

// EntStore implements the store.Store interface using Ent ORM.
type EntStore struct {
	client *ent.Client
}

// NewEntStore creates a new EntStore and initializes the database schema.
func NewEntStore(dataSourceName string) (*EntStore, error) {
	drv, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	client := ent.NewClient(
		ent.Driver(
			dialect.DebugWithContext(
				drv, func(ctx context.Context, i ...interface{}) {
					// Optional: Add debug logging here
				},
			),
		),
	)

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &EntStore{client: client}, nil
}

// Close closes the database connection.
func (s *EntStore) Close() error {
	return s.client.Close()
}

// Ensure EntStore implements the Store interface
var _ store.Store = (*EntStore)(nil)

// SingBox Configuration methods

func (s *EntStore) CreateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error {
	id := uuid.New().String()
	now := time.Now()

	// Convert model types to interface{} for Ent storage
	var inbounds []interface{}
	if config.Inbounds != nil {
		for _, inbound := range config.Inbounds {
			inbounds = append(inbounds, inbound)
		}
	}

	var outbounds []interface{}
	if config.Outbounds != nil {
		for _, outbound := range config.Outbounds {
			outbounds = append(outbounds, outbound)
		}
	}

	var endpoints []interface{}
	if config.Endpoints != nil {
		for _, endpoint := range config.Endpoints {
			endpoints = append(endpoints, endpoint)
		}
	}

	var experimental map[string]interface{}
	if config.Experimental != nil {
		experimental = *config.Experimental
	}

	var certificateConfig map[string]interface{}
	if config.Certificate != nil {
		// Convert certificate slice to map for storage
		certificateConfig = map[string]interface{}{
			"certificates": config.Certificate,
		}
	}

	_, err := s.client.SingBoxConfig.Create().
		SetID(id).
		SetName(config.Name).
		SetDescription(config.Description).
		SetLogConfig(config.Log).
		SetDNSConfig(config.DNS).
		SetNtpConfig(config.NTP).
		SetInbounds(inbounds).
		SetOutbounds(outbounds).
		SetRouteConfig(config.Route).
		SetExperimentalConfig(experimental).
		SetEndpoints(endpoints).
		SetCertificateConfig(certificateConfig).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Save(ctx)

	return err
}

func (s *EntStore) GetSingBoxConfig(ctx context.Context, id string) (*models.SingBoxConfig, error) {
	entConfig, err := s.client.SingBoxConfig.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("singbox config not found")
		}
		return nil, err
	}

	return s.entSingBoxConfigToModel(entConfig), nil
}

func (s *EntStore) ListSingBoxConfigs(ctx context.Context, limit, offset int) ([]*models.SingBoxConfig, error) {
	query := s.client.SingBoxConfig.Query().
		Order(ent.Desc(singboxconfig.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entConfigs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	configs := make([]*models.SingBoxConfig, len(entConfigs))
	for i, entConfig := range entConfigs {
		configs[i] = s.entSingBoxConfigToModel(entConfig)
	}

	return configs, nil
}

func (s *EntStore) UpdateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error {
	// Convert model types to interface{} for Ent storage
	var inbounds []interface{}
	if config.Inbounds != nil {
		for _, inbound := range config.Inbounds {
			inbounds = append(inbounds, inbound)
		}
	}

	var outbounds []interface{}
	if config.Outbounds != nil {
		for _, outbound := range config.Outbounds {
			outbounds = append(outbounds, outbound)
		}
	}

	var endpoints []interface{}
	if config.Endpoints != nil {
		for _, endpoint := range config.Endpoints {
			endpoints = append(endpoints, endpoint)
		}
	}

	var experimental map[string]interface{}
	if config.Experimental != nil {
		experimental = *config.Experimental
	}

	var certificateConfig map[string]interface{}
	if config.Certificate != nil {
		// Convert certificate slice to map for storage
		certificateConfig = map[string]interface{}{
			"certificates": config.Certificate,
		}
	}

	update := s.client.SingBoxConfig.UpdateOneID(config.ID).
		SetName(config.Name).
		SetDescription(config.Description).
		SetLogConfig(config.Log).
		SetDNSConfig(config.DNS).
		SetNtpConfig(config.NTP).
		SetInbounds(inbounds).
		SetOutbounds(outbounds).
		SetRouteConfig(config.Route).
		SetExperimentalConfig(experimental).
		SetEndpoints(endpoints).
		SetCertificateConfig(certificateConfig).
		SetUpdatedAt(time.Now())

	_, err := update.Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *EntStore) DeleteSingBoxConfig(ctx context.Context, id string) error {
	return s.client.SingBoxConfig.DeleteOneID(id).Exec(ctx)
}

// Xray Configuration methods

func (s *EntStore) CreateXrayConfig(ctx context.Context, config *models.XrayConfig) error {
	if config.ID == "" {
		config.ID = uuid.New().String()
	}
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}
	config.UpdatedAt = time.Now()

	_, err := s.client.XrayConfig.Create().
		SetID(config.ID).
		SetName(config.Name).
		SetDescription(config.Description).
		SetCreatedAt(config.CreatedAt).
		SetUpdatedAt(config.UpdatedAt).
		SetLog(config.Log).
		SetAPI(config.API).
		SetDNS(config.DNS).
		SetRouting(config.Routing).
		SetPolicy(config.Policy).
		SetInbounds(config.Inbounds).
		SetOutbounds(config.Outbounds).
		SetTransport(config.Transport).
		SetStats(config.Stats).
		SetReverse(config.Reverse).
		SetFakedns(config.FakeDNS).
		SetMetrics(config.Metrics).
		SetObservatory(config.Observatory).
		SetBurstObservatory(config.BurstObservatory).
		Save(ctx)

	return err
}

func (s *EntStore) GetXrayConfig(ctx context.Context, id string) (*models.XrayConfig, error) {
	entConfig, err := s.client.XrayConfig.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("xray config not found")
		}
		return nil, err
	}

	return s.entXrayConfigToModel(entConfig), nil
}

func (s *EntStore) ListXrayConfigs(ctx context.Context, limit, offset int) ([]*models.XrayConfig, error) {
	query := s.client.XrayConfig.Query().
		Order(ent.Desc(xrayconfig.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entConfigs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	configs := make([]*models.XrayConfig, len(entConfigs))
	for i, entConfig := range entConfigs {
		configs[i] = s.entXrayConfigToModel(entConfig)
	}

	return configs, nil
}

func (s *EntStore) UpdateXrayConfig(ctx context.Context, config *models.XrayConfig) error {
	config.UpdatedAt = time.Now()

	_, err := s.client.XrayConfig.UpdateOneID(config.ID).
		SetName(config.Name).
		SetDescription(config.Description).
		SetUpdatedAt(config.UpdatedAt).
		SetLog(config.Log).
		SetAPI(config.API).
		SetDNS(config.DNS).
		SetRouting(config.Routing).
		SetPolicy(config.Policy).
		SetInbounds(config.Inbounds).
		SetOutbounds(config.Outbounds).
		SetTransport(config.Transport).
		SetStats(config.Stats).
		SetReverse(config.Reverse).
		SetFakedns(config.FakeDNS).
		SetMetrics(config.Metrics).
		SetObservatory(config.Observatory).
		SetBurstObservatory(config.BurstObservatory).
		Save(ctx)

	return err
}

func (s *EntStore) DeleteXrayConfig(ctx context.Context, id string) error {
	return s.client.XrayConfig.DeleteOneID(id).Exec(ctx)
}

// Helper function to convert Ent SingBoxConfig to model
func (s *EntStore) entSingBoxConfigToModel(entConfig *ent.SingBoxConfig) *models.SingBoxConfig {
	// Convert interface{} slices back to typed slices
	var inbounds []*models.SingBoxInbound
	if entConfig.Inbounds != nil {
		for _, inbound := range entConfig.Inbounds {
			if typedInbound, ok := inbound.(*models.SingBoxInbound); ok {
				inbounds = append(inbounds, typedInbound)
			}
		}
	}

	var outbounds []*models.SingBoxOutbound
	if entConfig.Outbounds != nil {
		for _, outbound := range entConfig.Outbounds {
			if typedOutbound, ok := outbound.(*models.SingBoxOutbound); ok {
				outbounds = append(outbounds, typedOutbound)
			}
		}
	}

	var endpoints []map[string]interface{}
	if entConfig.Endpoints != nil {
		for _, endpoint := range entConfig.Endpoints {
			if typedEndpoint, ok := endpoint.(map[string]interface{}); ok {
				endpoints = append(endpoints, typedEndpoint)
			}
		}
	}

	var experimental *map[string]interface{}
	if entConfig.ExperimentalConfig != nil {
		experimental = &entConfig.ExperimentalConfig
	}

	var certificates []*models.SingBoxCertificate
	if entConfig.CertificateConfig != nil {
		if certData, ok := entConfig.CertificateConfig["certificates"]; ok {
			if certSlice, ok := certData.([]*models.SingBoxCertificate); ok {
				certificates = certSlice
			}
		}
	}

	return &models.SingBoxConfig{
		ID:           entConfig.ID,
		Name:         entConfig.Name,
		Description:  entConfig.Description,
		CreatedAt:    entConfig.CreatedAt,
		UpdatedAt:    entConfig.UpdatedAt,
		Log:          entConfig.LogConfig,
		DNS:          entConfig.DNSConfig,
		NTP:          entConfig.NtpConfig,
		Inbounds:     inbounds,
		Outbounds:    outbounds,
		Route:        entConfig.RouteConfig,
		Experimental: experimental,
		Endpoints:    endpoints,
		Certificate:  certificates,
	}
}

// HAProxy Configuration methods

func (s *EntStore) CreateHAProxyConfig(ctx context.Context, config *models.HAProxyConfig) error {
	id := uuid.New().String()
	now := time.Now()

	create := s.client.HAProxyConfig.Create().
		SetID(id).
		SetName(config.Name).
		SetDescription(config.Description).
		SetCreatedAt(now).
		SetUpdatedAt(now)

	if config.Global != nil {
		create = create.SetGlobalConfig(config.Global)
	}
	if config.Defaults != nil {
		create = create.SetDefaultsConfig(config.Defaults)
	}
	if config.Frontends != nil {
		create = create.SetFrontends(config.Frontends)
	}
	if config.Backends != nil {
		create = create.SetBackends(config.Backends)
	}
	if config.Listens != nil {
		create = create.SetListens(config.Listens)
	}
	if config.Stats != nil {
		create = create.SetStatsConfig(config.Stats)
	}

	_, err := create.Save(ctx)
	return err
}

func (s *EntStore) GetHAProxyConfig(ctx context.Context, id string) (*models.HAProxyConfig, error) {
	entConfig, err := s.client.HAProxyConfig.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("haproxy config not found")
		}
		return nil, err
	}

	return s.entHAProxyConfigToModel(entConfig), nil
}

func (s *EntStore) ListHAProxyConfigs(ctx context.Context, limit, offset int) ([]*models.HAProxyConfig, error) {
	query := s.client.HAProxyConfig.Query().
		Order(ent.Desc(haproxyconfig.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entConfigs, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	configs := make([]*models.HAProxyConfig, len(entConfigs))
	for i, entConfig := range entConfigs {
		configs[i] = s.entHAProxyConfigToModel(entConfig)
	}

	return configs, nil
}

func (s *EntStore) UpdateHAProxyConfig(ctx context.Context, config *models.HAProxyConfig) error {
	update := s.client.HAProxyConfig.UpdateOneID(config.ID).
		SetName(config.Name).
		SetDescription(config.Description).
		SetUpdatedAt(time.Now())

	if config.Global != nil {
		update = update.SetGlobalConfig(config.Global)
	}
	if config.Defaults != nil {
		update = update.SetDefaultsConfig(config.Defaults)
	}
	if config.Frontends != nil {
		update = update.SetFrontends(config.Frontends)
	}
	if config.Backends != nil {
		update = update.SetBackends(config.Backends)
	}
	if config.Listens != nil {
		update = update.SetListens(config.Listens)
	}
	if config.Stats != nil {
		update = update.SetStatsConfig(config.Stats)
	}

	_, err := update.Save(ctx)
	return err
}

func (s *EntStore) DeleteHAProxyConfig(ctx context.Context, id string) error {
	return s.client.HAProxyConfig.DeleteOneID(id).Exec(ctx)
}

// Helper function to convert Ent XrayConfig to model
func (s *EntStore) entXrayConfigToModel(entConfig *ent.XrayConfig) *models.XrayConfig {
	return &models.XrayConfig{
		ID:               entConfig.ID,
		Name:             entConfig.Name,
		Description:      entConfig.Description,
		CreatedAt:        entConfig.CreatedAt,
		UpdatedAt:        entConfig.UpdatedAt,
		Log:              entConfig.Log,
		API:              entConfig.API,
		DNS:              entConfig.DNS,
		Routing:          entConfig.Routing,
		Policy:           entConfig.Policy,
		Inbounds:         entConfig.Inbounds,
		Outbounds:        entConfig.Outbounds,
		Transport:        entConfig.Transport,
		Stats:            entConfig.Stats,
		Reverse:          entConfig.Reverse,
		FakeDNS:          entConfig.Fakedns,
		Metrics:          entConfig.Metrics,
		Observatory:      entConfig.Observatory,
		BurstObservatory: entConfig.BurstObservatory,
	}
}

// V2 Node management methods

func (s *EntStore) CreateNode(ctx context.Context, nodeCreate *models.NodeCreateV2) (*models.NodeV2, error) {
	id := uuid.New().String()
	now := time.Now()

	create := s.client.Node.Create().
		SetID(id).
		SetName(nodeCreate.Name).
		SetDescription(nodeCreate.Description).
		SetHostname(nodeCreate.Hostname).
		SetIPAddress(nodeCreate.IPAddress).
		SetPort(nodeCreate.Port).
		SetStatus("inactive").
		SetCreatedAt(now).
		SetUpdatedAt(now)

	if nodeCreate.Tags != nil {
		create = create.SetTags(nodeCreate.Tags)
	}

	entNode, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.entNodeToModel(entNode), nil
}

func (s *EntStore) GetNode(ctx context.Context, id string) (*models.NodeV2, error) {
	entNode, err := s.client.Node.Query().
		Where(node.ID(id)).
		WithServiceInstances().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("node not found")
		}
		return nil, err
	}

	return s.entNodeToModel(entNode), nil
}

func (s *EntStore) ListNodes(ctx context.Context, filters models.NodeFilters, limit, offset int) ([]*models.NodeV2, error) {
	query := s.client.Node.Query().
		WithServiceInstances().
		Order(ent.Desc(node.FieldCreatedAt))

	// Apply filters
	if filters.Status != "" {
		query = query.Where(node.Status(filters.Status))
	}
	if filters.Hostname != "" {
		query = query.Where(node.Hostname(filters.Hostname))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entNodes, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]*models.NodeV2, len(entNodes))
	for i, entNode := range entNodes {
		nodes[i] = s.entNodeToModel(entNode)
	}

	return nodes, nil
}

func (s *EntStore) UpdateNode(ctx context.Context, id string, updates *models.NodeUpdateV2) (*models.NodeV2, error) {
	update := s.client.Node.UpdateOneID(id).
		SetUpdatedAt(time.Now())

	if updates.Name != nil {
		update = update.SetName(*updates.Name)
	}
	if updates.Description != nil {
		update = update.SetDescription(*updates.Description)
	}
	if updates.Hostname != nil {
		update = update.SetHostname(*updates.Hostname)
	}
	if updates.IPAddress != nil {
		update = update.SetIPAddress(*updates.IPAddress)
	}
	if updates.Port != nil {
		update = update.SetPort(*updates.Port)
	}
	if updates.Status != nil {
		update = update.SetStatus(*updates.Status)
	}
	if updates.Tags != nil {
		update = update.SetTags(updates.Tags)
	}

	entNode, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch with relations
	entNode, err = s.client.Node.Query().
		Where(node.ID(id)).
		WithServiceInstances().
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return s.entNodeToModel(entNode), nil
}

func (s *EntStore) DeleteNode(ctx context.Context, id string) error {
	return s.client.Node.DeleteOneID(id).Exec(ctx)
}

// Helper function to convert Ent HAProxyConfig to model
func (s *EntStore) entHAProxyConfigToModel(entConfig *ent.HAProxyConfig) *models.HAProxyConfig {
	return &models.HAProxyConfig{
		ID:          entConfig.ID,
		Name:        entConfig.Name,
		Description: entConfig.Description,
		CreatedAt:   entConfig.CreatedAt,
		UpdatedAt:   entConfig.UpdatedAt,
		Global:      entConfig.GlobalConfig,
		Defaults:    entConfig.DefaultsConfig,
		Frontends:   entConfig.Frontends,
		Backends:    entConfig.Backends,
		Listens:     entConfig.Listens,
		Stats:       entConfig.StatsConfig,
	}
}

// V2 ServiceInstance management methods

func (s *EntStore) CreateServiceInstance(ctx context.Context, nodeId string, serviceCreate *models.ServiceInstanceCreateV2) (*models.ServiceInstanceV2, error) {
	id := uuid.New().String()
	now := time.Now()

	entService, err := s.client.ServiceInstance.Create().
		SetID(id).
		SetNodeID(nodeId).
		SetName(serviceCreate.Name).
		SetDescription(serviceCreate.Description).
		SetServiceType(serviceCreate.ServiceType).
		SetStatus("stopped").
		SetPort(serviceCreate.Port).
		SetProtocol(serviceCreate.Protocol).
		SetConfig(serviceCreate.Config).
		SetTags(serviceCreate.Tags).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return s.entServiceInstanceToModel(entService), nil
}

func (s *EntStore) GetServiceInstance(ctx context.Context, nodeId, serviceId string) (*models.ServiceInstanceV2, error) {
	entService, err := s.client.ServiceInstance.Query().
		Where(
			serviceinstance.ID(serviceId),
			serviceinstance.NodeID(nodeId),
		).
		WithNode().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("service instance not found")
		}
		return nil, err
	}

	return s.entServiceInstanceToModel(entService), nil
}

func (s *EntStore) ListServiceInstances(ctx context.Context, nodeId string, limit, offset int) ([]*models.ServiceInstanceV2, error) {
	query := s.client.ServiceInstance.Query().
		Where(serviceinstance.NodeID(nodeId)).
		WithNode().
		Order(ent.Desc(serviceinstance.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entServices, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	services := make([]*models.ServiceInstanceV2, len(entServices))
	for i, entService := range entServices {
		services[i] = s.entServiceInstanceToModel(entService)
	}

	return services, nil
}

func (s *EntStore) UpdateServiceInstance(ctx context.Context, nodeId, serviceId string, updates *models.ServiceInstanceUpdateV2) (*models.ServiceInstanceV2, error) {
	update := s.client.ServiceInstance.UpdateOneID(serviceId).
		SetUpdatedAt(time.Now())

	if updates.Name != nil {
		update = update.SetName(*updates.Name)
	}
	if updates.Description != nil {
		update = update.SetDescription(*updates.Description)
	}
	if updates.Status != nil {
		update = update.SetStatus(*updates.Status)
	}
	if updates.Port != nil {
		update = update.SetPort(*updates.Port)
	}
	if updates.Protocol != nil {
		update = update.SetProtocol(*updates.Protocol)
	}
	if updates.Config != nil {
		update = update.SetConfig(updates.Config)
	}
	if updates.Tags != nil {
		update = update.SetTags(updates.Tags)
	}

	entService, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch with relations
	entService, err = s.client.ServiceInstance.Query().
		Where(
			serviceinstance.ID(serviceId),
			serviceinstance.NodeID(nodeId),
		).
		WithNode().
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return s.entServiceInstanceToModel(entService), nil
}

func (s *EntStore) DeleteServiceInstance(ctx context.Context, nodeId, serviceId string) error {
	// First verify the service instance belongs to the specified node
	_, err := s.client.ServiceInstance.Query().
		Where(
			serviceinstance.ID(serviceId),
			serviceinstance.NodeID(nodeId),
		).
		Only(ctx)
	if err != nil {
		return err
	}

	return s.client.ServiceInstance.DeleteOneID(serviceId).Exec(ctx)
}

// Helper function to convert Ent Node to model
func (s *EntStore) entNodeToModel(entNode *ent.Node) *models.NodeV2 {
	node := &models.NodeV2{
		ID:          entNode.ID,
		Name:        entNode.Name,
		Description: entNode.Description,
		Hostname:    entNode.Hostname,
		IPAddress:   entNode.IPAddress,
		Port:        entNode.Port,
		Status:      entNode.Status,
		OSInfo:      entNode.OsInfo,
		AgentInfo:   entNode.AgentInfo,
		Tags:        entNode.Tags,
		CreatedAt:   entNode.CreatedAt,
		UpdatedAt:   entNode.UpdatedAt,
		LastSeen:    entNode.LastSeen,
	}



	return node
}

// AgentToken management methods

func (s *EntStore) CreateAgentToken(ctx context.Context, tokenCreate *models.AgentTokenCreate) (*models.AgentToken, error) {
	id := uuid.New().String()
	token := s.generateAgentToken()
	now := time.Now()

	entToken, err := s.client.AgentToken.Create().
		SetID(id).
		SetNodeID(tokenCreate.NodeID).
		SetToken(token).
		SetName(tokenCreate.Name).
		SetStatus("active").
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetNillableExpiresAt(tokenCreate.ExpiresAt).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return s.entAgentTokenToModel(entToken), nil
}

func (s *EntStore) GetAgentToken(ctx context.Context, id string) (*models.AgentToken, error) {
	entToken, err := s.client.AgentToken.Query().
		Where(agenttoken.ID(id)).
		WithNode().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("agent token not found")
		}
		return nil, err
	}

	return s.entAgentTokenToModel(entToken), nil
}

func (s *EntStore) GetAgentTokenByToken(ctx context.Context, token string) (*models.AgentToken, error) {
	entToken, err := s.client.AgentToken.Query().
		Where(agenttoken.Token(token)).
		WithNode().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("agent token not found")
		}
		return nil, err
	}

	// Update last used timestamp
	_, err = s.client.AgentToken.UpdateOneID(entToken.ID).
		SetLastUsed(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.entAgentTokenToModel(entToken), nil
}

func (s *EntStore) ListAgentTokens(ctx context.Context, filters models.AgentTokenFilters, limit, offset int) ([]*models.AgentToken, error) {
	query := s.client.AgentToken.Query().
		WithNode().
		Order(ent.Desc(agenttoken.FieldCreatedAt))

	// Apply filters
	if filters.NodeID != "" {
		query = query.Where(agenttoken.NodeID(filters.NodeID))
	}
	if filters.Status != "" {
		query = query.Where(agenttoken.Status(filters.Status))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	entTokens, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	tokens := make([]*models.AgentToken, len(entTokens))
	for i, entToken := range entTokens {
		tokens[i] = s.entAgentTokenToModel(entToken)
	}

	return tokens, nil
}

func (s *EntStore) UpdateAgentToken(ctx context.Context, id string, updates *models.AgentTokenUpdate) (*models.AgentToken, error) {
	update := s.client.AgentToken.UpdateOneID(id).
		SetUpdatedAt(time.Now())

	if updates.Name != nil {
		update = update.SetName(*updates.Name)
	}
	if updates.Status != nil {
		update = update.SetStatus(*updates.Status)
	}
	if updates.ExpiresAt != nil {
		update = update.SetNillableExpiresAt(updates.ExpiresAt)
	}

	entToken, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Fetch with relations
	entToken, err = s.client.AgentToken.Query().
		Where(agenttoken.ID(id)).
		WithNode().
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return s.entAgentTokenToModel(entToken), nil
}

func (s *EntStore) DeleteAgentToken(ctx context.Context, id string) error {
	return s.client.AgentToken.DeleteOneID(id).Exec(ctx)
}

func (s *EntStore) RevokeAgentToken(ctx context.Context, id string) error {
	_, err := s.client.AgentToken.UpdateOneID(id).
		SetStatus("revoked").
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

// Helper function to convert Ent ServiceInstance to model
func (s *EntStore) entServiceInstanceToModel(entService *ent.ServiceInstance) *models.ServiceInstanceV2 {
	service := &models.ServiceInstanceV2{
		ID:          entService.ID,
		NodeID:      entService.NodeID,
		Name:        entService.Name,
		Description: entService.Description,
		ServiceType: entService.ServiceType,
		Status:      entService.Status,
		Port:        entService.Port,
		Protocol:    entService.Protocol,
		Config:      entService.Config,
		Tags:        entService.Tags,
		CreatedAt:   entService.CreatedAt,
		UpdatedAt:   entService.UpdatedAt,
		LastSeen:    entService.LastSeen,
	}

	// Convert node if loaded
	if entService.Edges.Node != nil {
		service.Node = s.entNodeToModel(entService.Edges.Node)
	}

	return service
}

// Helper function to convert Ent AgentToken to model
func (s *EntStore) entAgentTokenToModel(entToken *ent.AgentToken) *models.AgentToken {
	token := &models.AgentToken{
		ID:        entToken.ID,
		NodeID:    entToken.NodeID,
		Token:     entToken.Token,
		Name:      entToken.Name,
		Status:    entToken.Status,
		ExpiresAt: entToken.ExpiresAt,
		CreatedAt: entToken.CreatedAt,
		UpdatedAt: entToken.UpdatedAt,
		LastUsed:  entToken.LastUsed,
	}

	// Convert node if loaded
	if entToken.Edges.Node != nil {
		token.Node = s.entNodeToModel(entToken.Edges.Node)
	}

	return token
}

// generateAgentToken generates a secure random token for agent authentication
func (s *EntStore) generateAgentToken() string {
	// Generate a 32-byte random token and encode as base64
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to timestamp-based token if random generation fails
		return fmt.Sprintf("token_%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
