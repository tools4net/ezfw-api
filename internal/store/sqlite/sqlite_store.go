package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/tools4net/ezfw/backend/internal/models"
)

// SQLiteStore implements the store.Store interface using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLiteStore and initializes the database schema.
func NewSQLiteStore(dataSourceName string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close() // Close the DB if ping fails
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	store := &SQLiteStore{db: db}
	if err := store.initSchema(); err != nil {
		db.Close() // Close the DB if schema init fails
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the necessary tables if they don't exist.
func (s *SQLiteStore) initSchema() error {
	createSingBoxTableSQL := `
    CREATE TABLE IF NOT EXISTS singbox_configs (
        id TEXT PRIMARY KEY,
        name TEXT,
        description TEXT,
        created_at DATETIME,
        updated_at DATETIME,
        log_config TEXT,
        dns_config TEXT,
        ntp_config TEXT,
        inbounds TEXT,
        outbounds TEXT,
        route_config TEXT,
        experimental_config TEXT,
        services_config TEXT,
        endpoints_config TEXT,
        certificate_config TEXT
    );`
	if _, err := s.db.Exec(createSingBoxTableSQL); err != nil {
		return fmt.Errorf("failed to create singbox_configs table: %w", err)
	}

	createXrayTableSQL := `
	CREATE TABLE IF NOT EXISTS xray_configs (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE,
		description TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		log_config TEXT,
		api_config TEXT,
		dns_config TEXT,
		routing_config TEXT,
		policy_config TEXT,
		inbounds TEXT,
		outbounds TEXT,
		transport_config TEXT,
		stats_config TEXT,
		reverse_config TEXT,
		fakedns_config TEXT,
		metrics_config TEXT,
		observatory_config TEXT,
		burst_observatory_config TEXT
	);`
	if _, err := s.db.Exec(createXrayTableSQL); err != nil {
		return fmt.Errorf("failed to create xray_configs table: %w", err)
	}
	return nil
}

// marshalToJSON marshals a Go struct into JSON or stores nil if the struct is nil.
// It includes specific checks for known model types to correctly handle nil pointers.
func marshalToJSON(v interface{}) (sql.NullString, error) {
	if v == nil {
		return sql.NullString{}, nil
	}

	// Type-specific nil checks for pointer fields
	switch val := v.(type) {
	// SingBox types
	case *models.SingBoxLogConfig:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.SingBoxDNSConfig:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.SingBoxNTPConfig:
		if val == nil {
			return sql.NullString{}, nil
		}
	case []models.SingBoxInbound: // Changed to value type as per singbox.go
		if val == nil { // Check if slice itself is nil
			return sql.NullString{}, nil
		}
		if len(val) == 0 {
			return sql.NullString{String: "[]", Valid: true}, nil
		} // Store empty array for empty slice
	case []models.SingBoxOutbound: // Changed to value type
		if val == nil {
			return sql.NullString{}, nil
		}
		if len(val) == 0 {
			return sql.NullString{String: "[]", Valid: true}, nil
		}
	case *models.SingBoxRouteConfig:
		if val == nil {
			return sql.NullString{}, nil
		}
	// Xray types (using new definitions from models/xray.go)
	case *models.LogObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.APIObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.DNSObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.RoutingObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.PolicyObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case []models.InboundObject: // Changed to value type as per xray.go
		if val == nil {
			return sql.NullString{}, nil
		}
		if len(val) == 0 {
			return sql.NullString{String: "[]", Valid: true}, nil
		}
	case []models.OutboundObject: // Changed to value type
		if val == nil {
			return sql.NullString{}, nil
		}
		if len(val) == 0 {
			return sql.NullString{String: "[]", Valid: true}, nil
		}
	case *models.TransportObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.StatsObject: // Note: StatsObject is an empty struct for Xray
		if val == nil { // This check might be redundant if it's always a non-nil empty struct
			return sql.NullString{}, nil
		}
		// For empty structs like StatsObject, marshal will produce "{}"
		// which is fine to store.
	case *models.ReverseObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.FakeDNSObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.MetricsObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.ObservatoryObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	case *models.BurstObservatoryObject:
		if val == nil {
			return sql.NullString{}, nil
		}
	// Generic map/slice types (often used for placeholders)
	case *map[string]interface{}: // For SingBox Experimental, etc.
		if val == nil {
			return sql.NullString{}, nil
		}
	case []map[string]interface{}: // For SingBox Services, Endpoints, Certificate
		if val == nil {
			return sql.NullString{}, nil
		}
		if len(val) == 0 {
			return sql.NullString{String: "[]", Valid: true}, nil
		}
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		return sql.NullString{}, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	if string(jsonData) == "null" { // If json.Marshal results in "null" for a non-nil but empty-like struct
		return sql.NullString{}, nil
	}
	return sql.NullString{String: string(jsonData), Valid: true}, nil
}

// unmarshalFromJSON unmarshals JSON data from sql.NullString into a target struct.
// Ptr is a pointer to the field that will hold the unmarshalled data, e.g., &config.Log.
func unmarshalFromJSON(ns sql.NullString, ptr interface{}) error {
	if !ns.Valid || ns.String == "" || ns.String == "null" {
		// Value is NULL in DB or empty string or literal "null", treat as nil/empty struct
		// Setting ptr to nil is tricky as ptr is interface{}, its underlying type needs to be a pointer.
		// The caller should ensure ptr is a pointer to a nillable type (pointer, slice, map).
		return nil
	}
	return json.Unmarshal([]byte(ns.String), ptr)
}

// --- SingBox Methods ---

// CreateSingBoxConfig creates a new SingBox configuration.
func (s *SQLiteStore) CreateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error {
	if config.ID == "" {
		config.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	config.CreatedAt = now
	config.UpdatedAt = now

	logJSON, err := marshalToJSON(config.Log)
	if err != nil {
		return fmt.Errorf("marshal Log: %w", err)
	}
	dnsJSON, err := marshalToJSON(config.DNS)
	if err != nil {
		return fmt.Errorf("marshal DNS: %w", err)
	}
	ntpJSON, err := marshalToJSON(config.NTP)
	if err != nil {
		return fmt.Errorf("marshal NTP: %w", err)
	}
	inboundsJSON, err := marshalToJSON(config.Inbounds)
	if err != nil {
		return fmt.Errorf("marshal Inbounds: %w", err)
	}
	outboundsJSON, err := marshalToJSON(config.Outbounds)
	if err != nil {
		return fmt.Errorf("marshal Outbounds: %w", err)
	}
	routeJSON, err := marshalToJSON(config.Route)
	if err != nil {
		return fmt.Errorf("marshal Route: %w", err)
	}
	experimentalJSON, err := marshalToJSON(config.Experimental)
	if err != nil {
		return fmt.Errorf("marshal Experimental: %w", err)
	}
	servicesJSON, err := marshalToJSON(config.Services)
	if err != nil {
		return fmt.Errorf("marshal Services: %w", err)
	}
	endpointsJSON, err := marshalToJSON(config.Endpoints)
	if err != nil {
		return fmt.Errorf("marshal Endpoints: %w", err)
	}
	certificateJSON, err := marshalToJSON(config.Certificate)
	if err != nil {
		return fmt.Errorf("marshal Certificate: %w", err)
	}

	stmt := `
    INSERT INTO singbox_configs (
        id, name, description, created_at, updated_at,
        log_config, dns_config, ntp_config, inbounds, outbounds, route_config,
        experimental_config, services_config, endpoints_config, certificate_config
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.ExecContext(
		ctx, stmt,
		config.ID, config.Name, config.Description, config.CreatedAt, config.UpdatedAt,
		logJSON, dnsJSON, ntpJSON, inboundsJSON, outboundsJSON, routeJSON,
		experimentalJSON, servicesJSON, endpointsJSON, certificateJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert singbox config: %w", err)
	}
	return nil
}

// GetSingBoxConfig retrieves a SingBox configuration by its ID.
func (s *SQLiteStore) GetSingBoxConfig(ctx context.Context, id string) (*models.SingBoxConfig, error) {
	stmt := `
    SELECT id, name, description, created_at, updated_at,
           log_config, dns_config, ntp_config, inbounds, outbounds, route_config,
           experimental_config, services_config, endpoints_config, certificate_config
    FROM singbox_configs WHERE id = ?`

	row := s.db.QueryRowContext(ctx, stmt, id)
	config := &models.SingBoxConfig{}

	var logJSON, dnsJSON, ntpJSON, inboundsJSON, outboundsJSON, routeJSON sql.NullString
	var experimentalJSON, servicesJSON, endpointsJSON, certificateJSON sql.NullString

	err := row.Scan(
		&config.ID, &config.Name, &config.Description, &config.CreatedAt, &config.UpdatedAt,
		&logJSON, &dnsJSON, &ntpJSON, &inboundsJSON, &outboundsJSON, &routeJSON,
		&experimentalJSON, &servicesJSON, &endpointsJSON, &certificateJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("singbox config with id %s not found: %w", id, sql.ErrNoRows) // Wrap ErrNoRows
		}
		return nil, fmt.Errorf("failed to scan singbox config: %w", err)
	}

	if err := unmarshalFromJSON(logJSON, &config.Log); err != nil {
		return nil, fmt.Errorf("unmarshal Log: %w", err)
	}
	if err := unmarshalFromJSON(dnsJSON, &config.DNS); err != nil {
		return nil, fmt.Errorf("unmarshal DNS: %w", err)
	}
	if err := unmarshalFromJSON(ntpJSON, &config.NTP); err != nil {
		return nil, fmt.Errorf("unmarshal NTP: %w", err)
	}
	if err := unmarshalFromJSON(inboundsJSON, &config.Inbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Inbounds: %w", err)
	}
	if err := unmarshalFromJSON(outboundsJSON, &config.Outbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Outbounds: %w", err)
	}
	if err := unmarshalFromJSON(routeJSON, &config.Route); err != nil {
		return nil, fmt.Errorf("unmarshal Route: %w", err)
	}
	if err := unmarshalFromJSON(experimentalJSON, &config.Experimental); err != nil {
		return nil, fmt.Errorf("unmarshal Experimental: %w", err)
	}
	if err := unmarshalFromJSON(servicesJSON, &config.Services); err != nil {
		return nil, fmt.Errorf("unmarshal Services: %w", err)
	}
	if err := unmarshalFromJSON(endpointsJSON, &config.Endpoints); err != nil {
		return nil, fmt.Errorf("unmarshal Endpoints: %w", err)
	}
	if err := unmarshalFromJSON(certificateJSON, &config.Certificate); err != nil {
		return nil, fmt.Errorf("unmarshal Certificate: %w", err)
	}

	return config, nil
}

// GetXrayConfigByName retrieves an Xray configuration by its name.
func (s *SQLiteStore) GetXrayConfigByName(ctx context.Context, name string) (*models.XrayConfig, error) {
	stmt := `
    SELECT id, name, description, created_at, updated_at,
           log_config, api_config, dns_config, routing_config, policy_config,
           inbounds, outbounds, transport_config, stats_config, reverse_config,
           fakedns_config, metrics_config, observatory_config, burst_observatory_config
    FROM xray_configs WHERE name = ?`

	row := s.db.QueryRowContext(ctx, stmt, name)
	config := &models.XrayConfig{}

	var logJ, apiJ, dnsJ, routingJ, policyJ, inboundsJ, outboundsJ, transportJ, statsJ, reverseJ, fakednsJ, metricsJ, obsJ, burstObsJ sql.NullString

	err := row.Scan(
		&config.ID, &config.Name, &config.Description, &config.CreatedAt, &config.UpdatedAt,
		&logJ, &apiJ, &dnsJ, &routingJ, &policyJ, &inboundsJ, &outboundsJ, &transportJ,
		&statsJ, &reverseJ, &fakednsJ, &metricsJ, &obsJ, &burstObsJ,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("xray config with name %s not found: %w", name, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to scan xray config by name: %w", err)
	}

	// Unmarshal JSON blobs
	if err := unmarshalFromJSON(logJ, &config.Log); err != nil {
		return nil, fmt.Errorf("unmarshal Log: %w", err)
	}
	if err := unmarshalFromJSON(apiJ, &config.API); err != nil {
		return nil, fmt.Errorf("unmarshal API: %w", err)
	}
	if err := unmarshalFromJSON(dnsJ, &config.DNS); err != nil {
		return nil, fmt.Errorf("unmarshal DNS: %w", err)
	}
	if err := unmarshalFromJSON(routingJ, &config.Routing); err != nil {
		return nil, fmt.Errorf("unmarshal Routing: %w", err)
	}
	if err := unmarshalFromJSON(policyJ, &config.Policy); err != nil {
		return nil, fmt.Errorf("unmarshal Policy: %w", err)
	}
	if err := unmarshalFromJSON(inboundsJ, &config.Inbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Inbounds: %w", err)
	}
	if err := unmarshalFromJSON(outboundsJ, &config.Outbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Outbounds: %w", err)
	}
	if err := unmarshalFromJSON(transportJ, &config.Transport); err != nil {
		return nil, fmt.Errorf("unmarshal Transport: %w", err)
	}
	if err := unmarshalFromJSON(statsJ, &config.Stats); err != nil {
		return nil, fmt.Errorf("unmarshal Stats: %w", err)
	}
	if err := unmarshalFromJSON(reverseJ, &config.Reverse); err != nil {
		return nil, fmt.Errorf("unmarshal Reverse: %w", err)
	}
	if err := unmarshalFromJSON(fakednsJ, &config.FakeDNS); err != nil {
		return nil, fmt.Errorf("unmarshal FakeDNS: %w", err)
	}
	if err := unmarshalFromJSON(metricsJ, &config.Metrics); err != nil {
		return nil, fmt.Errorf("unmarshal Metrics: %w", err)
	}
	if err := unmarshalFromJSON(obsJ, &config.Observatory); err != nil {
		return nil, fmt.Errorf("unmarshal Observatory: %w", err)
	}
	if err := unmarshalFromJSON(burstObsJ, &config.BurstObservatory); err != nil {
		return nil, fmt.Errorf("unmarshal BurstObservatory: %w", err)
	}

	return config, nil
}

// ListSingBoxConfigs retrieves a list of SingBox configurations with pagination.
func (s *SQLiteStore) ListSingBoxConfigs(ctx context.Context, limit, offset int) ([]*models.SingBoxConfig, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	stmt := `
    SELECT id, name, description, created_at, updated_at,
           log_config, dns_config, ntp_config, inbounds, outbounds, route_config,
           experimental_config, services_config, endpoints_config, certificate_config
    FROM singbox_configs ORDER BY updated_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query singbox configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.SingBoxConfig
	for rows.Next() {
		config := &models.SingBoxConfig{}
		var logJSON, dnsJSON, ntpJSON, inboundsJSON, outboundsJSON, routeJSON sql.NullString
		var experimentalJSON, servicesJSON, endpointsJSON, certificateJSON sql.NullString

		err := rows.Scan(
			&config.ID, &config.Name, &config.Description, &config.CreatedAt, &config.UpdatedAt,
			&logJSON, &dnsJSON, &ntpJSON, &inboundsJSON, &outboundsJSON, &routeJSON,
			&experimentalJSON, &servicesJSON, &endpointsJSON, &certificateJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan singbox config row: %w", err)
		}

		if err := unmarshalFromJSON(logJSON, &config.Log); err != nil {
			return nil, fmt.Errorf("unmarshal Log for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(dnsJSON, &config.DNS); err != nil {
			return nil, fmt.Errorf("unmarshal DNS for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(ntpJSON, &config.NTP); err != nil {
			return nil, fmt.Errorf("unmarshal NTP for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(inboundsJSON, &config.Inbounds); err != nil {
			return nil, fmt.Errorf("unmarshal Inbounds for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(outboundsJSON, &config.Outbounds); err != nil {
			return nil, fmt.Errorf("unmarshal Outbounds for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(routeJSON, &config.Route); err != nil {
			return nil, fmt.Errorf("unmarshal Route for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(experimentalJSON, &config.Experimental); err != nil {
			return nil, fmt.Errorf("unmarshal Experimental for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(servicesJSON, &config.Services); err != nil {
			return nil, fmt.Errorf("unmarshal Services for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(endpointsJSON, &config.Endpoints); err != nil {
			return nil, fmt.Errorf("unmarshal Endpoints for %s: %w", config.ID, err)
		}
		if err := unmarshalFromJSON(certificateJSON, &config.Certificate); err != nil {
			return nil, fmt.Errorf("unmarshal Certificate for %s: %w", config.ID, err)
		}
		configs = append(configs, config)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating singbox config rows: %w", err)
	}
	return configs, nil
}

// UpdateSingBoxConfig updates an existing SingBox configuration.
func (s *SQLiteStore) UpdateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error {
	if config.ID == "" {
		return fmt.Errorf("cannot update singbox config: ID is missing")
	}
	config.UpdatedAt = time.Now().UTC()

	logJSON, err := marshalToJSON(config.Log)
	if err != nil {
		return fmt.Errorf("marshal Log: %w", err)
	}
	dnsJSON, err := marshalToJSON(config.DNS)
	if err != nil {
		return fmt.Errorf("marshal DNS: %w", err)
	}
	ntpJSON, err := marshalToJSON(config.NTP)
	if err != nil {
		return fmt.Errorf("marshal NTP: %w", err)
	}
	inboundsJSON, err := marshalToJSON(config.Inbounds)
	if err != nil {
		return fmt.Errorf("marshal Inbounds: %w", err)
	}
	outboundsJSON, err := marshalToJSON(config.Outbounds)
	if err != nil {
		return fmt.Errorf("marshal Outbounds: %w", err)
	}
	routeJSON, err := marshalToJSON(config.Route)
	if err != nil {
		return fmt.Errorf("marshal Route: %w", err)
	}
	experimentalJSON, err := marshalToJSON(config.Experimental)
	if err != nil {
		return fmt.Errorf("marshal Experimental: %w", err)
	}
	servicesJSON, err := marshalToJSON(config.Services)
	if err != nil {
		return fmt.Errorf("marshal Services: %w", err)
	}
	endpointsJSON, err := marshalToJSON(config.Endpoints)
	if err != nil {
		return fmt.Errorf("marshal Endpoints: %w", err)
	}
	certificateJSON, err := marshalToJSON(config.Certificate)
	if err != nil {
		return fmt.Errorf("marshal Certificate: %w", err)
	}

	stmt := `
    UPDATE singbox_configs SET
        name = ?, description = ?, updated_at = ?,
        log_config = ?, dns_config = ?, ntp_config = ?, inbounds = ?, outbounds = ?, route_config = ?,
        experimental_config = ?, services_config = ?, endpoints_config = ?, certificate_config = ?
    WHERE id = ?`

	result, err := s.db.ExecContext(
		ctx, stmt,
		config.Name, config.Description, config.UpdatedAt,
		logJSON, dnsJSON, ntpJSON, inboundsJSON, outboundsJSON, routeJSON,
		experimentalJSON, servicesJSON, endpointsJSON, certificateJSON,
		config.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update singbox config: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for singbox update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("singbox config with id %s not found for update: %w", config.ID, sql.ErrNoRows)
	}
	return nil
}

// DeleteSingBoxConfig deletes a SingBox configuration by its ID.
func (s *SQLiteStore) DeleteSingBoxConfig(ctx context.Context, id string) error {
	stmt := `DELETE FROM singbox_configs WHERE id = ?`
	result, err := s.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete singbox config: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for singbox delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("singbox config with id %s not found for deletion: %w", id, sql.ErrNoRows)
	}
	return nil
}

// --- Xray Methods ---

// CreateXrayConfig creates a new Xray configuration.
func (s *SQLiteStore) CreateXrayConfig(ctx context.Context, config *models.XrayConfig) error {

	if config.ID == "" {
		config.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	config.CreatedAt = now
	config.UpdatedAt = now

	var err error
	logJSON, err := marshalToJSON(config.Log)
	if err != nil {
		return fmt.Errorf("marshal Log: %w", err)
	}
	apiJSON, err := marshalToJSON(config.API)
	if err != nil {
		return fmt.Errorf("marshal API: %w", err)
	}
	dnsJSON, err := marshalToJSON(config.DNS)
	if err != nil {
		return fmt.Errorf("marshal DNS: %w", err)
	}
	routingJSON, err := marshalToJSON(config.Routing)
	if err != nil {
		return fmt.Errorf("marshal Routing: %w", err)
	}
	policyJSON, err := marshalToJSON(config.Policy)
	if err != nil {
		return fmt.Errorf("marshal Policy: %w", err)
	}
	inboundsJSON, err := marshalToJSON(config.Inbounds)
	if err != nil {
		return fmt.Errorf("marshal Inbounds: %w", err)
	}
	outboundsJSON, err := marshalToJSON(config.Outbounds)
	if err != nil {
		return fmt.Errorf("marshal Outbounds: %w", err)
	}
	transportJSON, err := marshalToJSON(config.Transport)
	if err != nil {
		return fmt.Errorf("marshal Transport: %w", err)
	}
	statsJSON, err := marshalToJSON(config.Stats)
	if err != nil {
		return fmt.Errorf("marshal Stats: %w", err)
	}
	reverseJSON, err := marshalToJSON(config.Reverse)
	if err != nil {
		return fmt.Errorf("marshal Reverse: %w", err)
	}
	fakednsJSON, err := marshalToJSON(config.FakeDNS)
	if err != nil {
		return fmt.Errorf("marshal FakeDNS: %w", err)
	}
	metricsJSON, err := marshalToJSON(config.Metrics)
	if err != nil {
		return fmt.Errorf("marshal Metrics: %w", err)
	}
	observatoryJSON, err := marshalToJSON(config.Observatory)
	if err != nil {
		return fmt.Errorf("marshal Observatory: %w", err)
	}
	burstObservatoryJSON, err := marshalToJSON(config.BurstObservatory)
	if err != nil {
		return fmt.Errorf("marshal BurstObservatory: %w", err)
	}

	stmt := `
    INSERT INTO xray_configs (
        id, name, description, created_at, updated_at,
        log_config, api_config, dns_config, routing_config, policy_config,
        inbounds, outbounds, transport_config, stats_config, reverse_config,
        fakedns_config, metrics_config, observatory_config, burst_observatory_config
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.db.ExecContext(
		ctx, stmt,
		config.ID, config.Name, config.Description, config.CreatedAt, config.UpdatedAt,
		logJSON, apiJSON, dnsJSON, routingJSON, policyJSON,
		inboundsJSON, outboundsJSON, transportJSON, statsJSON, reverseJSON,
		fakednsJSON, metricsJSON, observatoryJSON, burstObservatoryJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert xray config: %w", err)
	}
	return nil
}

// GetXrayConfig retrieves an Xray configuration by its ID.
func (s *SQLiteStore) GetXrayConfig(ctx context.Context, id string) (*models.XrayConfig, error) {
	stmt := `
    SELECT id, name, description, created_at, updated_at,
           log_config, api_config, dns_config, routing_config, policy_config,
           inbounds, outbounds, transport_config, stats_config, reverse_config,
           fakedns_config, metrics_config, observatory_config, burst_observatory_config
    FROM xray_configs WHERE id = ?`

	row := s.db.QueryRowContext(ctx, stmt, id)
	config := &models.XrayConfig{}

	var logJ, apiJ, dnsJ, routingJ, policyJ, inboundsJ, outboundsJ, transportJ, statsJ, reverseJ, fakednsJ, metricsJ, obsJ, burstObsJ sql.NullString

	err := row.Scan(
		&config.ID, &config.Name, &config.Description, &config.CreatedAt, &config.UpdatedAt,
		&logJ, &apiJ, &dnsJ, &routingJ, &policyJ, &inboundsJ, &outboundsJ, &transportJ,
		&statsJ, &reverseJ, &fakednsJ, &metricsJ, &obsJ, &burstObsJ,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("xray config with id %s not found: %w", id, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to scan xray config: %w", err)
	}

	// Unmarshal JSON blobs
	if err := unmarshalFromJSON(logJ, &config.Log); err != nil {
		return nil, fmt.Errorf("unmarshal Log: %w", err)
	}
	if err := unmarshalFromJSON(apiJ, &config.API); err != nil {
		return nil, fmt.Errorf("unmarshal API: %w", err)
	}
	if err := unmarshalFromJSON(dnsJ, &config.DNS); err != nil {
		return nil, fmt.Errorf("unmarshal DNS: %w", err)
	}
	if err := unmarshalFromJSON(routingJ, &config.Routing); err != nil {
		return nil, fmt.Errorf("unmarshal Routing: %w", err)
	}
	if err := unmarshalFromJSON(policyJ, &config.Policy); err != nil {
		return nil, fmt.Errorf("unmarshal Policy: %w", err)
	}
	if err := unmarshalFromJSON(inboundsJ, &config.Inbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Inbounds: %w", err)
	}
	if err := unmarshalFromJSON(outboundsJ, &config.Outbounds); err != nil {
		return nil, fmt.Errorf("unmarshal Outbounds: %w", err)
	}
	if err := unmarshalFromJSON(transportJ, &config.Transport); err != nil {
		return nil, fmt.Errorf("unmarshal Transport: %w", err)
	}
	if err := unmarshalFromJSON(statsJ, &config.Stats); err != nil {
		return nil, fmt.Errorf("unmarshal Stats: %w", err)
	}
	if err := unmarshalFromJSON(reverseJ, &config.Reverse); err != nil {
		return nil, fmt.Errorf("unmarshal Reverse: %w", err)
	}
	if err := unmarshalFromJSON(fakednsJ, &config.FakeDNS); err != nil {
		return nil, fmt.Errorf("unmarshal FakeDNS: %w", err)
	}
	if err := unmarshalFromJSON(metricsJ, &config.Metrics); err != nil {
		return nil, fmt.Errorf("unmarshal Metrics: %w", err)
	}
	if err := unmarshalFromJSON(obsJ, &config.Observatory); err != nil {
		return nil, fmt.Errorf("unmarshal Observatory: %w", err)
	}
	if err := unmarshalFromJSON(burstObsJ, &config.BurstObservatory); err != nil {
		return nil, fmt.Errorf("unmarshal BurstObservatory: %w", err)
	}

	return config, nil
}

// ListXrayConfigs retrieves a list of Xray configurations with pagination.
func (s *SQLiteStore) ListXrayConfigs(ctx context.Context, limit, offset int) ([]*models.XrayConfig, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	stmt := `
    SELECT id, name, description, created_at, updated_at,
           log_config, api_config, dns_config, routing_config, policy_config,
           inbounds, outbounds, transport_config, stats_config, reverse_config,
           fakedns_config, metrics_config, observatory_config, burst_observatory_config
    FROM xray_configs ORDER BY updated_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query xray configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.XrayConfig
	for rows.Next() {
		config := &models.XrayConfig{}
		var logJ, apiJ, dnsJ, routingJ, policyJ, inboundsJ, outboundsJ, transportJ, statsJ, reverseJ, fakednsJ, metricsJ, obsJ, burstObsJ sql.NullString
		err := rows.Scan(
			&config.ID, &config.Name, &config.Description, &config.CreatedAt, &config.UpdatedAt,
			&logJ, &apiJ, &dnsJ, &routingJ, &policyJ, &inboundsJ, &outboundsJ, &transportJ,
			&statsJ, &reverseJ, &fakednsJ, &metricsJ, &obsJ, &burstObsJ,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan xray config row: %w", err)
		}

		if errU := unmarshalFromJSON(logJ, &config.Log); errU != nil {
			return nil, fmt.Errorf("unmarshal Log for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(apiJ, &config.API); errU != nil {
			return nil, fmt.Errorf("unmarshal API for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(dnsJ, &config.DNS); errU != nil {
			return nil, fmt.Errorf("unmarshal DNS for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(routingJ, &config.Routing); errU != nil {
			return nil, fmt.Errorf("unmarshal Routing for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(policyJ, &config.Policy); errU != nil {
			return nil, fmt.Errorf("unmarshal Policy for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(inboundsJ, &config.Inbounds); errU != nil {
			return nil, fmt.Errorf("unmarshal Inbounds for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(outboundsJ, &config.Outbounds); errU != nil {
			return nil, fmt.Errorf("unmarshal Outbounds for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(transportJ, &config.Transport); errU != nil {
			return nil, fmt.Errorf("unmarshal Transport for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(statsJ, &config.Stats); errU != nil {
			return nil, fmt.Errorf("unmarshal Stats for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(reverseJ, &config.Reverse); errU != nil {
			return nil, fmt.Errorf("unmarshal Reverse for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(fakednsJ, &config.FakeDNS); errU != nil {
			return nil, fmt.Errorf("unmarshal FakeDNS for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(metricsJ, &config.Metrics); errU != nil {
			return nil, fmt.Errorf("unmarshal Metrics for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(obsJ, &config.Observatory); errU != nil {
			return nil, fmt.Errorf("unmarshal Observatory for %s: %w", config.ID, errU)
		}
		if errU := unmarshalFromJSON(burstObsJ, &config.BurstObservatory); errU != nil {
			return nil, fmt.Errorf("unmarshal BurstObservatory for %s: %w", config.ID, errU)
		}
		configs = append(configs, config)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating xray config rows: %w", err)
	}
	return configs, nil
}

// UpdateXrayConfig updates an existing Xray configuration.
func (s *SQLiteStore) UpdateXrayConfig(ctx context.Context, config *models.XrayConfig) error {
	if config.ID == "" {
		return fmt.Errorf("cannot update xray config: ID is missing")
	}
	config.UpdatedAt = time.Now().UTC()

	var err error
	logJSON, err := marshalToJSON(config.Log)
	if err != nil {
		return fmt.Errorf("marshal Log: %w", err)
	}
	apiJSON, err := marshalToJSON(config.API)
	if err != nil {
		return fmt.Errorf("marshal API: %w", err)
	}
	dnsJSON, err := marshalToJSON(config.DNS)
	if err != nil {
		return fmt.Errorf("marshal DNS: %w", err)
	}
	routingJSON, err := marshalToJSON(config.Routing)
	if err != nil {
		return fmt.Errorf("marshal Routing: %w", err)
	}
	policyJSON, err := marshalToJSON(config.Policy)
	if err != nil {
		return fmt.Errorf("marshal Policy: %w", err)
	}
	inboundsJSON, err := marshalToJSON(config.Inbounds)
	if err != nil {
		return fmt.Errorf("marshal Inbounds: %w", err)
	}
	outboundsJSON, err := marshalToJSON(config.Outbounds)
	if err != nil {
		return fmt.Errorf("marshal Outbounds: %w", err)
	}
	transportJSON, err := marshalToJSON(config.Transport)
	if err != nil {
		return fmt.Errorf("marshal Transport: %w", err)
	}
	statsJSON, err := marshalToJSON(config.Stats)
	if err != nil {
		return fmt.Errorf("marshal Stats: %w", err)
	}
	reverseJSON, err := marshalToJSON(config.Reverse)
	if err != nil {
		return fmt.Errorf("marshal Reverse: %w", err)
	}
	fakednsJSON, err := marshalToJSON(config.FakeDNS)
	if err != nil {
		return fmt.Errorf("marshal FakeDNS: %w", err)
	}
	metricsJSON, err := marshalToJSON(config.Metrics)
	if err != nil {
		return fmt.Errorf("marshal Metrics: %w", err)
	}
	observatoryJSON, err := marshalToJSON(config.Observatory)
	if err != nil {
		return fmt.Errorf("marshal Observatory: %w", err)
	}
	burstObservatoryJSON, err := marshalToJSON(config.BurstObservatory)
	if err != nil {
		return fmt.Errorf("marshal BurstObservatory: %w", err)
	}

	stmt := `
    UPDATE xray_configs SET
        name = ?, description = ?, updated_at = ?,
        log_config = ?, api_config = ?, dns_config = ?, routing_config = ?, policy_config = ?,
        inbounds = ?, outbounds = ?, transport_config = ?, stats_config = ?, reverse_config = ?,
        fakedns_config = ?, metrics_config = ?, observatory_config = ?, burst_observatory_config = ?
    WHERE id = ?`

	result, err := s.db.ExecContext(
		ctx, stmt,
		config.Name, config.Description, config.UpdatedAt,
		logJSON, apiJSON, dnsJSON, routingJSON, policyJSON,
		inboundsJSON, outboundsJSON, transportJSON, statsJSON, reverseJSON,
		fakednsJSON, metricsJSON, observatoryJSON, burstObservatoryJSON,
		config.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update xray config: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for xray update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("xray config with id %s not found for update: %w", config.ID, sql.ErrNoRows)
	}
	return nil
}

// DeleteXrayConfig deletes an Xray configuration by its ID.
func (s *SQLiteStore) DeleteXrayConfig(ctx context.Context, id string) error {
	stmt := `DELETE FROM xray_configs WHERE id = ?`
	result, err := s.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete xray config: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for xray delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("xray config with id %s not found for deletion: %w", id, sql.ErrNoRows)
	}
	return nil
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
