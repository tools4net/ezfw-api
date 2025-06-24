package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

// CreateSingBoxConfigHandler handles POST requests to create a new SingBox configuration
func (h *ConfigHandler) CreateSingBoxConfigHandler(c *gin.Context) {
	var config models.SingBoxConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that the configuration name is not empty
	if !validateConfigName(c, config.Name, "SingBox Configuration") {
		return
	}

	if err := h.store.CreateSingBoxConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// ListSingBoxConfigsHandler handles GET requests to list SingBox configurations
func (h *ConfigHandler) ListSingBoxConfigsHandler(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	configs, err := h.store.ListSingBoxConfigs(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetSingBoxConfigHandler handles GET requests to retrieve a specific SingBox configuration
func (h *ConfigHandler) GetSingBoxConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	config, err := h.store.GetSingBoxConfig(c.Request.Context(), configID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateSingBoxConfigHandler handles PUT requests to update a SingBox configuration
func (h *ConfigHandler) UpdateSingBoxConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	var config models.SingBoxConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID matches the URL parameter
	config.ID = configID

	// Validate that the configuration name is not empty if provided
	if !validateConfigName(c, config.Name, "SingBox Configuration") {
		return
	}

	if err := h.store.UpdateSingBoxConfig(c.Request.Context(), &config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found for update"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteSingBoxConfigHandler handles DELETE requests to remove a SingBox configuration
func (h *ConfigHandler) DeleteSingBoxConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	if err := h.store.DeleteSingBoxConfig(c.Request.Context(), configID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found for deletion"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration deleted successfully"})
}

// GenerateSingBoxConfigHandler handles GET requests to generate a SingBox configuration file
func (h *ConfigHandler) GenerateSingBoxConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	config, err := h.store.GetSingBoxConfig(c.Request.Context(), configID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration: " + err.Error()})
		}
		return
	}

	// Return the configuration in a format suitable for SingBox, stripping internal metadata
	singboxCoreConfig := gin.H{
		"log":          config.Log,
		"dns":          config.DNS,
		"ntp":          config.NTP,
		"inbounds":     config.Inbounds,
		"outbounds":    config.Outbounds,
		"route":        config.Route,
		"experimental": config.Experimental,
		"services":     config.Services,
		"endpoints":    config.Endpoints,
		"certificate":  config.Certificate,
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=singbox-config-"+config.Name+".json")
	c.JSON(http.StatusOK, singboxCoreConfig)
}

// CreateXrayConfigHandler handles POST requests to create a new Xray configuration
func (h *ConfigHandler) CreateXrayConfigHandler(c *gin.Context) {
	var config models.XrayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Basic validation: Name is required and cannot be empty
	if !validateConfigName(c, config.Name, "Xray Configuration") {
		return
	}

	// ID will be generated by the store if not provided, or use provided if it is.
	// For Create, we typically let store handle ID generation.
	// If ID is provided by client during create, it might be overwritten or cause unique constraint issues if not careful.
	// Current store CreateXrayConfig generates a new UUID if config.ID is empty.

	if err := h.store.CreateXrayConfig(c.Request.Context(), &config); err != nil {
		// Check for unique constraint error on name (specific to SQLite error message)
		if strings.Contains(err.Error(), "UNIQUE constraint failed: xray_configs.name") {
			c.JSON(http.StatusConflict, gin.H{"error": "Configuration name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Xray configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

// ListXrayConfigsHandler handles GET requests to list Xray configurations
func (h *ConfigHandler) ListXrayConfigsHandler(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	configs, err := h.store.ListXrayConfigs(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// GetXrayConfigHandler handles GET requests to retrieve a specific Xray configuration
func (h *ConfigHandler) GetXrayConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	config, err := h.store.GetXrayConfig(c.Request.Context(), configID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateXrayConfigHandler handles PUT requests to update an Xray configuration
func (h *ConfigHandler) UpdateXrayConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	var config models.XrayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Name is required in the payload for an update and cannot be empty.
	if !validateConfigName(c, config.Name, "Xray Configuration") {
		return
	}


	// Ensure the ID matches the URL parameter
	config.ID = configID

	if err := h.store.UpdateXrayConfig(c.Request.Context(), &config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found for update"})
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: xray_configs.name") { // Check for unique constraint error on name
			c.JSON(http.StatusConflict, gin.H{"error": "Configuration name already exists for another configuration"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Xray configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteXrayConfigHandler handles DELETE requests to remove an Xray configuration
func (h *ConfigHandler) DeleteXrayConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	if err := h.store.DeleteXrayConfig(c.Request.Context(), configID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found for deletion"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Xray configuration: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration deleted successfully"})
}

// GenerateXrayConfigHandler handles GET requests to generate an Xray configuration file
func (h *ConfigHandler) GenerateXrayConfigHandler(c *gin.Context) {
	configID, ok := validateConfigID(c)
	if !ok {
		return
	}

	config, err := h.store.GetXrayConfig(c.Request.Context(), configID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configuration: " + err.Error()})
		}
		return
	}

	// Return the configuration in a format suitable for Xray
	// We only want to marshal the core Xray fields, not our internal DB/API metadata
	xrayCoreConfig := gin.H{
		"log":               config.Log,
		"api":               config.API,
		"dns":               config.DNS,
		"routing":           config.Routing,
		"policy":            config.Policy,
		"inbounds":          config.Inbounds,
		"outbounds":         config.Outbounds,
		"transport":         config.Transport,
		"stats":             config.Stats,
		"reverse":           config.Reverse,
		"fakedns":           config.FakeDNS,
		"metrics":           config.Metrics,
		"observatory":       config.Observatory,
		"burstObservatory": config.BurstObservatory,
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=xray-config-"+config.Name+".json") // Add config name to filename
	c.JSON(http.StatusOK, xrayCoreConfig)
}

// --- Helper Functions ---

// validateConfigID checks if the configId path parameter is present.
// If not, it writes a 400 error to the context and returns (id="", ok=false).
// Otherwise, it returns (id=configID, ok=true).
func validateConfigID(c *gin.Context) (string, bool) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config ID is required in path"})
		return "", false
	}
	return configID, true
}

// validateConfigName checks if the Name field of a model is empty (after trimming spaces).
// If empty, it writes a 400 error to the context and returns false.
// entityType is a string like "SingBox Configuration" or "Xray Configuration" for the error message.
func validateConfigName(c *gin.Context, name string, entityType string) bool {
	if strings.TrimSpace(name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": entityType + " 'name' is required and cannot be empty"})
		return false
	}
	return true
}