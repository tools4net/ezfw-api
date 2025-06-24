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

// CreateSingBoxConfigHandler creates a new SingBox configuration.
// @Summary      Create SingBox Config
// @Description  Adds a new SingBox configuration to the store.
// @Tags         ConfigsSingBox
// @Accept       json
// @Produce      json
// @Param        singBoxConfig body models.SingBoxConfig true "SingBox Configuration Object"
// @Success      201 {object} models.SingBoxConfig "Successfully created SingBox configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid input or missing name"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to create configuration"
// @Security     ApiKeyAuth
// @Router       /configs/singbox [post]
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

// ListSingBoxConfigsHandler lists all SingBox configurations with pagination.
// @Summary      List SingBox Configs
// @Description  Retrieves a paginated list of SingBox configurations.
// @Tags         ConfigsSingBox
// @Produce      json
// @Param        limit query int false "Number of items to return per page" default(10) example(10)
// @Param        offset query int false "Offset for pagination" default(0) example(0)
// @Success      200 {object} map[string][]models.SingBoxConfig "A list of SingBox configurations"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to list configurations"
// @Security     ApiKeyAuth
// @Router       /configs/singbox [get]
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

// GetSingBoxConfigHandler retrieves a specific SingBox configuration by its ID.
// @Summary      Get SingBox Config by ID
// @Description  Fetches a single SingBox configuration based on its unique ID.
// @Tags         ConfigsSingBox
// @Produce      json
// @Param        configId path string true "SingBox Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} models.SingBoxConfig "Successfully retrieved SingBox configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to retrieve configuration"
// @Security     ApiKeyAuth
// @Router       /configs/singbox/{configId} [get]
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

// UpdateSingBoxConfigHandler updates an existing SingBox configuration.
// @Summary      Update SingBox Config
// @Description  Modifies an existing SingBox configuration by its ID.
// @Tags         ConfigsSingBox
// @Accept       json
// @Produce      json
// @Param        configId path string true "SingBox Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Param        singBoxConfig body models.SingBoxConfig true "SingBox Configuration Object"
// @Success      200 {object} models.SingBoxConfig "Successfully updated SingBox configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid input, missing ID, or missing name"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found for update"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to update configuration"
// @Security     ApiKeyAuth
// @Router       /configs/singbox/{configId} [put]
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

// DeleteSingBoxConfigHandler removes a SingBox configuration by its ID.
// @Summary      Delete SingBox Config
// @Description  Deletes a specific SingBox configuration by its ID.
// @Tags         ConfigsSingBox
// @Produce      json
// @Param        configId path string true "SingBox Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} map[string]string "Success message"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found for deletion"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to delete configuration"
// @Security     ApiKeyAuth
// @Router       /configs/singbox/{configId} [delete]
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

// GenerateSingBoxConfigHandler generates a sanitized SingBox configuration file by ID.
// @Summary      Generate SingBox Config File
// @Description  Retrieves a SingBox configuration and returns it in a format suitable for the sing-box binary, stripping internal metadata.
// @Tags         ConfigsSingBox
// @Produce      json
// @Param        configId path string true "SingBox Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} object "The core SingBox configuration JSON. The actual structure depends on the config (models.SingBoxConfig without metadata)."
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to retrieve configuration"
// @Security     ApiKeyAuth
// @Router       /configs/singbox/{configId}/generate [get]
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

// CreateXrayConfigHandler creates a new Xray configuration.
// @Summary      Create Xray Config
// @Description  Adds a new Xray configuration to the store. Name must be unique.
// @Tags         ConfigsXray
// @Accept       json
// @Produce      json
// @Param        xrayConfig body models.XrayConfig true "Xray Configuration Object"
// @Success      201 {object} models.XrayConfig "Successfully created Xray configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid input or missing name"
// @Failure      409 {object} models.ErrorResponse "Conflict - Configuration name already exists"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to create configuration"
// @Security     ApiKeyAuth
// @Router       /configs/xray [post]
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

// ListXrayConfigsHandler lists all Xray configurations with pagination.
// @Summary      List Xray Configs
// @Description  Retrieves a paginated list of Xray configurations.
// @Tags         ConfigsXray
// @Produce      json
// @Param        limit query int false "Number of items to return per page" default(10) example(10)
// @Param        offset query int false "Offset for pagination" default(0) example(0)
// @Success      200 {object} map[string][]models.XrayConfig "A list of Xray configurations"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to list configurations"
// @Security     ApiKeyAuth
// @Router       /configs/xray [get]
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

// GetXrayConfigHandler retrieves a specific Xray configuration by its ID.
// @Summary      Get Xray Config by ID
// @Description  Fetches a single Xray configuration based on its unique ID.
// @Tags         ConfigsXray
// @Produce      json
// @Param        configId path string true "Xray Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} models.XrayConfig "Successfully retrieved Xray configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to retrieve configuration"
// @Security     ApiKeyAuth
// @Router       /configs/xray/{configId} [get]
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

// UpdateXrayConfigHandler updates an existing Xray configuration.
// @Summary      Update Xray Config
// @Description  Modifies an existing Xray configuration by its ID. Name must be unique among Xray configs.
// @Tags         ConfigsXray
// @Accept       json
// @Produce      json
// @Param        configId path string true "Xray Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Param        xrayConfig body models.XrayConfig true "Xray Configuration Object"
// @Success      200 {object} models.XrayConfig "Successfully updated Xray configuration"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid input, missing ID, or missing name"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found for update"
// @Failure      409 {object} models.ErrorResponse "Conflict - Configuration name already exists for another configuration"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to update configuration"
// @Security     ApiKeyAuth
// @Router       /configs/xray/{configId} [put]
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

// DeleteXrayConfigHandler removes an Xray configuration by its ID.
// @Summary      Delete Xray Config
// @Description  Deletes a specific Xray configuration by its ID.
// @Tags         ConfigsXray
// @Produce      json
// @Param        configId path string true "Xray Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} map[string]string "Success message"
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found for deletion"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to delete configuration"
// @Security     ApiKeyAuth
// @Router       /configs/xray/{configId} [delete]
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

// GenerateXrayConfigHandler generates a sanitized Xray configuration file by ID.
// @Summary      Generate Xray Config File
// @Description  Retrieves an Xray configuration and returns it in a format suitable for the Xray-core binary, stripping internal metadata.
// @Tags         ConfigsXray
// @Produce      json
// @Param        configId path string true "Xray Configuration ID" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// @Success      200 {object} object "The core Xray configuration JSON. The actual structure depends on the config (models.XrayConfig without metadata)."
// @Failure      400 {object} models.ErrorResponse "Bad Request - Invalid or missing configuration ID"
// @Failure      404 {object} models.ErrorResponse "Not Found - Configuration not found"
// @Failure      500 {object} models.ErrorResponse "Internal Server Error - Failed to retrieve configuration"
// @Security     ApiKeyAuth
// @Router       /configs/xray/{configId}/generate [get]
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
		"services":          config.Services, // Added missing services field
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