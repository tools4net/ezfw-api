package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tools4net/ezfw/backend/internal/models"
)

func setupTestDB(t *testing.T) (*SQLiteStore, func()) {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "testdb_")
	require.NoError(t, err)

	dbPath := filepath.Join(tempDir, "test_proxypanel.db")
	store, err := NewSQLiteStore(dbPath)
	require.NoError(t, err, "Failed to create test SQLite store")

	cleanup := func() {
		err := store.Close()
		assert.NoError(t, err, "Failed to close test DB")
		err = os.RemoveAll(tempDir) // Remove the temp directory and its contents
		assert.NoError(t, err, "Failed to remove temp DB directory")
	}

	return store, cleanup
}

func TestCreateSingBoxConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	config := &models.SingBoxConfig{
		Name:        "Test Config 1",
		Description: "A test configuration",
		Log:         &models.SingBoxLogConfig{Level: ("info")},
		DNS:         &models.SingBoxDNSConfig{Servers: []struct{Address string "json:\"address\""}{{"8.8.8.8"}}},
		Inbounds: []models.SingBoxInbound{
			{Type: "mixed", Tag: "mixed-in", ListenPort: (1080)},
		},
		Outbounds: []models.SingBoxOutbound{
			{Type: "direct", Tag: "direct-out"},
		},
	}

	err := store.CreateSingBoxConfig(ctx, config)
	require.NoError(t, err)

	retrieved, err := store.GetSingBoxConfig(ctx, config.Tag)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, config.Name, retrieved.Name)
	assert.Equal(t, config.Description, retrieved.Description)
	assert.NotNil(t, retrieved.Log)
	assert.Equal(t, "info", *retrieved.Log.Level)
	assert.NotNil(t, retrieved.DNS)
	require.Len(t, retrieved.DNS.Servers, 1)
	assert.Equal(t, "8.8.8.8", *retrieved.DNS.Servers[0].Address)
	require.Len(t, retrieved.Inbounds, 1)
	assert.Equal(t, "mixed", retrieved.Inbounds[0].Type)
	assert.Equal(t, 1080, *retrieved.Inbounds[0].ListenPort)
}

func TestGetSingBoxConfig_NotFound(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	_, err := store.GetSingBoxConfig(ctx, uuid.NewString())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdateSingBoxConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	originalConfig := &models.SingBoxConfig{Name: "Original Name"}
	err := store.CreateSingBoxConfig(ctx, originalConfig)
	require.NoError(t, err)

	originalConfig.Name = "Updated Name"
	originalConfig.Description = "Updated Description"
	originalConfig.Log = &models.SingBoxLogConfig{Level: ("debug")}
	originalUpdatedAt := originalConfig.UpdatedAt

	// Brief pause to ensure UpdatedAt changes
	time.Sleep(10 * time.Millisecond)

	err = store.UpdateSingBoxConfig(ctx, originalConfig)
	require.NoError(t, err)

	updatedConfig, err := store.GetSingBoxConfig(ctx, originalConfig.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedConfig.Name)
	assert.Equal(t, "Updated Description", updatedConfig.Description)
	require.NotNil(t, updatedConfig.Log)
	assert.Equal(t, "debug", *updatedConfig.Log.Level)
	assert.True(t, updatedConfig.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be more recent")
}

func TestUpdateSingBoxConfig_NotFound(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	nonExistentConfig := &models.SingBoxConfig{ID: uuid.NewString(), Name: "Non Existent"}
	err := store.UpdateSingBoxConfig(ctx, nonExistentConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found for update")
}

func TestDeleteSingBoxConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	config := &models.SingBoxConfig{Name: "To Be Deleted"}
	err := store.CreateSingBoxConfig(ctx, config)
	require.NoError(t, err)

	err = store.DeleteSingBoxConfig(ctx, config.ID)
	require.NoError(t, err)

	_, err = store.GetSingBoxConfig(ctx, config.ID)
	require.Error(t, err, "Expected error when getting deleted config")
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteSingBoxConfig_NotFound(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	err := store.DeleteSingBoxConfig(ctx, uuid.NewString())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found for deletion")
}

func TestListSingBoxConfigs(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a few configs
	config1 := &models.SingBoxConfig{Name: "Config A", Log: &models.SingBoxLogConfig{Level: ("info")}}
	time.Sleep(5 * time.Millisecond) // Ensure different UpdatedAt
	config2 := &models.SingBoxConfig{Name: "Config B", Log: &models.SingBoxLogConfig{Level: ("debug")}}
	time.Sleep(5 * time.Millisecond)
	config3 := &models.SingBoxConfig{Name: "Config C", Log: &models.SingBoxLogConfig{Level: ("warn")}}

	require.NoError(t, store.CreateSingBoxConfig(ctx, config1))
	require.NoError(t, store.CreateSingBoxConfig(ctx, config2))
	require.NoError(t, store.CreateSingBoxConfig(ctx, config3))

	// Test listing all (limit > count)
	configs, err := store.ListSingBoxConfigs(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, configs, 3)
	// Check order (default is by UpdatedAt DESC)
	assert.Equal(t, config3.ID, configs[0].ID)
	assert.Equal(t, config2.ID, configs[1].ID)
	assert.Equal(t, config1.ID, configs[2].ID)
	assert.NotNil(t, configs[0].Log) // Check if JSON unmarshalling worked
	assert.Equal(t, "warn", *configs[0].Log.Level)

	// Test limit
	configs, err = store.ListSingBoxConfigs(ctx, 1, 0)
	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, config3.ID, configs[0].ID)

	// Test offset
	configs, err = store.ListSingBoxConfigs(ctx, 1, 1)
	require.NoError(t, err)
	require.Len(t, configs, 1)
	assert.Equal(t, config2.ID, configs[0].ID)

	// Test empty list
	storeNoConf, cleanupNoConf := setupTestDB(t)
	defer cleanupNoConf()
	emptyConfigs, err := storeNoConf.ListSingBoxConfigs(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, emptyConfigs, 0)
}

// Helper function to create a pointer to a string.
// () etc. can be used directly in tests.
// This is just for local use if models package isn't directly modifiable for test helpers.
func (s string) *string { return &s }
func (i int) *int          { return &i }

// Add a test for marshalling/unmarshalling of all field types
func TestSingBoxConfig_JSONMarshalling(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	fullConfig := &models.SingBoxConfig{
		Name:        "Full JSON Test",
		Description: "Testing all fields",
		Log:         &models.SingBoxLogConfig{Level: ("debug"), Timestamp: (true)},
		DNS: &models.SingBoxDNSConfig{
			Servers:  []*models.SingBoxDNSServer{{Address: ("1.1.1.1"), Strategy: ("ipv4_only")}},
			Strategy: ("prefer_ipv6"),
			FakeIP:   &models.SingBoxFakeIPConfig{Enabled: (true), Inet4Range: ("198.18.0.0/15")},
		},
		NTP: &models.SingBoxNTPConfig{Enabled: (true), Server: ("time.apple.com")},
		Inbounds: []models.SingBoxInbound{
			{Type: "socks", Tag: "socks-in", ListenPort: (1080), Sniff: (true)},
		},
		Outbounds: []models.SingBoxOutbound{
			{Type: "direct", Tag: "direct-out"},
			{Type: "vmess", Tag: "vmess-out", Settings: map[string]interface{}{"address": "server.com", "port": 12345}},
		},
		Route: &models.SingBoxRouteConfig{
			Final: ("direct-out"),
			Rules: []*models.SingBoxRouteRule{{Outbound: ("vmess-out"), Domain: []string{"test.com"}}},
		},
		Experimental: &map[string]interface{}{"cache_file": "/path/to/cache"},
		Services:     []*map[string]interface{}{{"type": "some_service", "enabled": true}},
		Endpoints:    []*map[string]interface{}{{"type": "wg", "interface_name": "wg0"}},
		Certificate:  []*map[string]interface{}{{"ca_file": "/path/to/ca.pem"}},
	}

	err := store.CreateSingBoxConfig(ctx, fullConfig)
	require.NoError(t, err)
	require.NotEmpty(t, fullConfig.ID)

	retrieved, err := store.GetSingBoxConfig(ctx, fullConfig.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	// Assertions for all fields
	assert.Equal(t, fullConfig.Name, retrieved.Name)
	assert.Equal(t, fullConfig.Description, retrieved.Description)

	require.NotNil(t, retrieved.Log)
	assert.Equal(t, *fullConfig.Log.Level, *retrieved.Log.Level)
	assert.Equal(t, *fullConfig.Log.Timestamp, *retrieved.Log.Timestamp)

	require.NotNil(t, retrieved.DNS)
	assert.Equal(t, *fullConfig.DNS.Strategy, *retrieved.DNS.Strategy)
	require.Len(t, retrieved.DNS.Servers, 1)
	assert.Equal(t, *fullConfig.DNS.Servers[0].Address, *retrieved.DNS.Servers[0].Address)
	require.NotNil(t, retrieved.DNS.FakeIP)
	assert.Equal(t, *fullConfig.DNS.FakeIP.Enabled, *retrieved.DNS.FakeIP.Enabled)
	assert.Equal(t, *fullConfig.DNS.FakeIP.Inet4Range, *retrieved.DNS.FakeIP.Inet4Range)

	require.NotNil(t, retrieved.NTP)
	assert.Equal(t, *fullConfig.NTP.Enabled, *retrieved.NTP.Enabled)
	assert.Equal(t, *full_config.NTP.Server, *retrieved.NTP.Server)

	require.Len(t, retrieved.Inbounds, 1)
	assert.Equal(t, fullConfig.Inbounds[0].Type, retrieved.Inbounds[0].Type)
	assert.Equal(t, *fullConfig.Inbounds[0].ListenPort, *retrieved.Inbounds[0].ListenPort)
	assert.Equal(t, *fullConfig.Inbounds[0].Sniff, *retrieved.Inbounds[0].Sniff)

	require.Len(t, retrieved.Outbounds, 2)
	assert.Equal(t, fullConfig.Outbounds[1].Type, retrieved.Outbounds[1].Type)
	assert.Equal(t, fullConfig.Outbounds[1].Settings["address"], retrieved.Outbounds[1].Settings["address"])

	require.NotNil(t, retrieved.Route)
	assert.Equal(t, *fullConfig.Route.Final, *retrieved.Route.Final)
	require.Len(t, retrieved.Route.Rules, 1)
	assert.Equal(t, *fullConfig.Route.Rules[0].Outbound, *retrieved.Route.Rules[0].Outbound)

	require.NotNil(t, retrieved.Experimental)
	assert.Equal(t, (*fullConfig.Experimental)["cache_file"], (*retrieved.Experimental)["cache_file"])

	require.Len(t, retrieved.Services, 1)
	assert.Equal(t, (*fullConfig.Services[0])["type"], (*retrieved.Services[0])["type"])

	require.Len(t, retrieved.Endpoints, 1)
	assert.Equal(t, (*fullConfig.Endpoints[0])["interface_name"], (*retrieved.Endpoints[0])["interface_name"])

	require.Len(t, retrieved.Certificate, 1)
	assert.Equal(t, (*fullConfig.Certificate[0])["ca_file"], (*retrieved.Certificate[0])["ca_file"])
}

// --- Xray Tests ---

func TestCreateXrayConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	config := &models.XrayConfig{
		Name:        "Test Xray Config",
		Description: "An Xray test configuration",
		Log:         &models.LogObject{Loglevel: "warning"},
		DNS:         &models.DNSObject{ClientIP: "1.2.3.4"}, // Using a different field for variety
		Inbounds: []models.InboundObject{ // Note: Inbounds is []InboundObject, not []*InboundObject in the model
			{Protocol: "socks", Port: 1088, Tag: "socks-xray-in"},
		},
	}

	err := store.CreateXrayConfig(ctx, config)
	require.NoError(t, err)
	require.NotEmpty(t, config.ID)

	retrieved, err := store.GetXrayConfig(ctx, config.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, config.Name, retrieved.Name)
	require.NotNil(t, retrieved.Log)
	assert.Equal(t, "warning", retrieved.Log.Loglevel) // No pointer for Loglevel in new LogObject
	require.NotNil(t, retrieved.DNS)
	assert.Equal(t, "1.2.3.4", retrieved.DNS.ClientIP) // No pointer for ClientIP
	require.Len(t, retrieved.Inbounds, 1)
	assert.Equal(t, "socks", retrieved.Inbounds[0].Protocol)

	// Handle Port being interface{}
	// When unmarshalling from JSON, numbers are typically float64 unless specified otherwise
	portVal, ok := retrieved.Inbounds[0].Port.(float64)
	require.True(t, ok, "Port should be unmarshalled as float64 from JSON number")
	assert.Equal(t, float64(1088), portVal)
}

func TestCreateXrayConfig_NameConflict(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	config1 := &models.XrayConfig{Name: "Conflict Xray"}
	err := store.CreateXrayConfig(ctx, config1)
	require.NoError(t, err)

	config2 := &models.XrayConfig{Name: "Conflict Xray"}
	err = store.CreateXrayConfig(ctx, config2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed: xray_configs.name")
}

func TestGetXrayConfig_NotFound(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	_, err := store.GetXrayConfig(ctx, uuid.NewString())
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found"))
}

func TestUpdateXrayConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	original := &models.XrayConfig{Name: "Xray Original", Description: "Original Desc"}
	err := store.CreateXrayConfig(ctx, original)
	require.NoError(t, err)
	require.NotEmpty(t, original.ID)

	toUpdate := &models.XrayConfig{
		ID:          original.ID, // Must provide ID for update
		Name:        "Xray Updated",
		Description: "Updated Desc",
		API:         &models.APIObject{Tag: "api-tag"},
		// CreatedAt will be ignored, UpdatedAt will be set by the store
	}
	originalUpdatedAt := original.UpdatedAt // Store this before it's potentially modified by GetXrayConfig

	// Retrieve and store UpdatedAt from DB to ensure accurate comparison
	// because CreateXrayConfig sets it, and we need the DB version.
	createdConfig, err := store.GetXrayConfig(ctx, original.ID)
	require.NoError(t, err)
	originalUpdatedAtFromDB := createdConfig.UpdatedAt

	time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt will change

	err = store.UpdateXrayConfig(ctx, toUpdate)
	require.NoError(t, err)

	updated, err := store.GetXrayConfig(ctx, original.ID)
	require.NoError(t, err)
	assert.Equal(t, "Xray Updated", updated.Name)
	assert.Equal(t, "Updated Desc", updated.Description)
	require.NotNil(t, updated.API)
	assert.Equal(t, "api-tag", updated.API.Tag) // No pointer for Tag in new APIObject
	assert.True(t, updated.UpdatedAt.After(originalUpdatedAtFromDB), "UpdatedAt should be more recent")
}

func TestUpdateXrayConfig_NotFound(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	nonExistentConfig := &models.XrayConfig{ID: uuid.NewString(), Name: "Non Existent"}
	err := store.UpdateXrayConfig(ctx, nonExistentConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found for update")
}

func TestUpdateXrayConfig_NameConflict(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	config1 := &models.XrayConfig{Name: "Xray Config One"}
	err := store.CreateXrayConfig(ctx, config1)
	require.NoError(t, err)

	config2 := &models.XrayConfig{Name: "Xray Config Two"}
	err = store.CreateXrayConfig(ctx, config2)
	require.NoError(t, err)

	// Try to update config2 to have config1's name
	config2.Name = "Xray Config One"
	err = store.UpdateXrayConfig(ctx, config2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed: xray_configs.name")
}


func TestDeleteXrayConfig(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	config := &models.XrayConfig{Name: "Xray To Delete"}
	err := store.CreateXrayConfig(ctx, config)
	require.NoError(t, err)

	err = store.DeleteXrayConfig(ctx, config.ID)
	require.NoError(t, err)

	_, err = store.GetXrayConfig(ctx, config.ID)
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found"))
}

func TestListXrayConfigs(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	x1 := &models.XrayConfig{Name: "Xray Cfg 1", Log: &models.LogObject{Loglevel: "info"}}
	time.Sleep(5 * time.Millisecond)
	x2 := &models.XrayConfig{Name: "Xray Cfg 2", Log: &models.LogObject{Loglevel: "debug"}}

	require.NoError(t, store.CreateXrayConfig(ctx, x1))
	require.NoError(t, store.CreateXrayConfig(ctx, x2))

	configs, err := store.ListXrayConfigs(ctx, 5, 0)
	require.NoError(t, err)
	require.Len(t, configs, 2)
	assert.Equal(t, x2.ID, configs[0].ID) // Ordered by UpdatedAt DESC
	assert.Equal(t, x1.ID, configs[1].ID)
	require.NotNil(t, configs[0].Log)
	assert.Equal(t, "debug", configs[0].Log.Loglevel)

	// Test empty list
	storeNoConf, cleanupNoConf := setupTestDB(t)
	defer cleanupNoConf()
	emptyConfigs, err := storeNoConf.ListXrayConfigs(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, emptyConfigs, 0)

}

func TestXrayConfig_JSONMarshalling_Partial(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	partialXrayConfig := &models.XrayConfig{
		Name: "Partial Xray JSON Test",
		Log:  &models.LogObject{Loglevel: "error", Access: "/var/log/xray/access.log"},
		API:  &models.APIObject{Tag: "proxy-api", Services: []string{"StatsService"}},
		Inbounds: []models.InboundObject{ // Value type for slice elements
			{
				Protocol: "vless",
				Port:     "443", // Port as string
				Tag:      "vless-in",
				Settings: map[string]interface{}{"decryption": "none"},
			},
		},
		Outbounds: []models.OutboundObject{ // Value type for slice elements
			{Protocol: "freedom", Tag: "direct", Settings: map[string]interface{}{"domainStrategy": "UseIP"}},
		},
		Stats:   &models.StatsObject{}, // Empty struct to enable stats
		FakeDNS: &models.FakeDNSObject{IPPool: "198.18.0.0/15"},
	}

	err := store.CreateXrayConfig(ctx, partialXrayConfig)
	require.NoError(t, err)
	require.NotEmpty(t, partialXrayConfig.ID)

	retrieved, err := store.GetXrayConfig(ctx, partialXrayConfig.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	assert.Equal(t, partialXrayConfig.Name, retrieved.Name)

	require.NotNil(t, retrieved.Log)
	assert.Equal(t, partialXrayConfig.Log.Loglevel, retrieved.Log.Loglevel)
	assert.Equal(t, partialXrayConfig.Log.Access, retrieved.Log.Access)


	require.NotNil(t, retrieved.API)
	assert.Equal(t, partialXrayConfig.API.Tag, retrieved.API.Tag)
	assert.Equal(t, partialXrayConfig.API.Services, retrieved.API.Services)


	require.Len(t, retrieved.Inbounds, 1)
	assert.Equal(t, "vless", retrieved.Inbounds[0].Protocol)
	assert.Equal(t, "443", retrieved.Inbounds[0].Port) // Port is string, direct comparison
	require.NotNil(t, retrieved.Inbounds[0].Settings)
	assert.Equal(t, "none", retrieved.Inbounds[0].Settings["decryption"])


	require.Len(t, retrieved.Outbounds, 1)
	assert.Equal(t, "freedom", retrieved.Outbounds[0].Protocol)
	require.NotNil(t, retrieved.Outbounds[0].Settings)
	assert.Equal(t, "UseIP", retrieved.Outbounds[0].Settings["domainStrategy"])


	require.NotNil(t, retrieved.Stats) // Just checking it's not nil, as it's an empty struct

	require.NotNil(t, retrieved.FakeDNS)
	assert.Equal(t, partialXrayConfig.FakeDNS.IPPool, retrieved.FakeDNS.IPPool)

	// Check that fields not set are nil or empty
	assert.Nil(t, retrieved.DNS)
	assert.Nil(t, retrieved.Routing)
	assert.Nil(t, retrieved.Policy)
	assert.Nil(t, retrieved.Transport)
	assert.Nil(t, retrieved.Reverse)
	assert.Nil(t, retrieved.Metrics)
	assert.Nil(t, retrieved.Observatory)
	assert.Nil(t, retrieved.BurstObservatory)
}

func TestGetXrayConfigByName(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	config := &models.XrayConfig{
		Name: "Xray Find By Name",
		Log:  &models.LogObject{Loglevel: "debug"},
	}
	err := store.CreateXrayConfig(ctx, config)
	require.NoError(t, err)
	require.NotEmpty(t, config.ID)

	// Test found
	retrieved, err := store.GetXrayConfigByName(ctx, "Xray Find By Name")
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, config.ID, retrieved.ID)
	assert.Equal(t, "Xray Find By Name", retrieved.Name)
	require.NotNil(t, retrieved.Log)
	assert.Equal(t, "debug", retrieved.Log.Loglevel)

	// Test not found
	_, err = store.GetXrayConfigByName(ctx, "NonExistentName")
	require.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found"))
}


// It's good practice to also add pointer helpers to the models package itself if possible,
// or a shared utility package.
// For example, in models/utils.go:
// func (s string) *string { return &s }
// func (i int) *int { return &i }
// func BoolPtr(b bool) *bool { return &b }
// Then they can be used as ("value")
