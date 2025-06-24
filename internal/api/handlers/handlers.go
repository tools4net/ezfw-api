package handlers

import (
	"net/http"
	"strconv"

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
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	config, err := h.store.GetSingBoxConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateSingBoxConfigHandler handles PUT requests to update a SingBox configuration
func (h *ConfigHandler) UpdateSingBoxConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	var config models.SingBoxConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID matches the URL parameter
	config.ID = configID

	if err := h.store.UpdateSingBoxConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteSingBoxConfigHandler handles DELETE requests to remove a SingBox configuration
func (h *ConfigHandler) DeleteSingBoxConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	if err := h.store.DeleteSingBoxConfig(c.Request.Context(), configID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration deleted successfully"})
}

// GenerateSingBoxConfigHandler handles GET requests to generate a SingBox configuration file
func (h *ConfigHandler) GenerateSingBoxConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	config, err := h.store.GetSingBoxConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	// Return the configuration in a format suitable for SingBox
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=singbox-config.json")
	c.JSON(http.StatusOK, config)
}

// CreateXrayConfigHandler handles POST requests to create a new Xray configuration
func (h *ConfigHandler) CreateXrayConfigHandler(c *gin.Context) {
	var config models.XrayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.store.CreateXrayConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	config, err := h.store.GetXrayConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateXrayConfigHandler handles PUT requests to update an Xray configuration
func (h *ConfigHandler) UpdateXrayConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	var config models.XrayConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the ID matches the URL parameter
	config.ID = configID

	if err := h.store.UpdateXrayConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// DeleteXrayConfigHandler handles DELETE requests to remove an Xray configuration
func (h *ConfigHandler) DeleteXrayConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	if err := h.store.DeleteXrayConfig(c.Request.Context(), configID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration deleted successfully"})
}

// GenerateXrayConfigHandler handles GET requests to generate an Xray configuration file
func (h *ConfigHandler) GenerateXrayConfigHandler(c *gin.Context) {
	configID := c.Param("configId")
	if configID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config ID is required"})
		return
	}

	config, err := h.store.GetXrayConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	// Return the configuration in a format suitable for Xray
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=xray-config.json")
	c.JSON(http.StatusOK, config)
}