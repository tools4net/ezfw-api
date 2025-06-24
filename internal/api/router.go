package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tools4net/ezfw/backend/internal/api/handlers"
	"github.com/tools4net/ezfw/backend/internal/store"
)

// SetupRouter initializes the Gin router and sets up API routes.
func SetupRouter(router *gin.Engine, dbStore store.Store) {
	// Create handlers
	configHandler := handlers.NewConfigHandler(dbStore)

	// Group API routes under /api/v1
	v1 := router.Group("/api/v1")
	{
		// Health check for API v1
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "API v1 is healthy"})
		})

		// SingBox Config Routes
		singboxGroup := v1.Group("/configs/singbox")
		{
			singboxGroup.POST("", configHandler.CreateSingBoxConfigHandler)
			singboxGroup.GET("", configHandler.ListSingBoxConfigsHandler)
			singboxGroup.GET("/:configId", configHandler.GetSingBoxConfigHandler)
			singboxGroup.PUT("/:configId", configHandler.UpdateSingBoxConfigHandler)
			singboxGroup.DELETE("/:configId", configHandler.DeleteSingBoxConfigHandler)
			singboxGroup.GET("/:configId/generate", configHandler.GenerateSingBoxConfigHandler)
		}

		// Xray Config Routes
		xrayGroup := v1.Group("/configs/xray")
		{
			xrayGroup.POST("", configHandler.CreateXrayConfigHandler)
			xrayGroup.GET("", configHandler.ListXrayConfigsHandler)
			xrayGroup.GET("/:configId", configHandler.GetXrayConfigHandler)
			xrayGroup.PUT("/:configId", configHandler.UpdateXrayConfigHandler)
			xrayGroup.DELETE("/:configId", configHandler.DeleteXrayConfigHandler)
			xrayGroup.GET("/:configId/generate", configHandler.GenerateXrayConfigHandler)
		}
	}

	// It's generally better to let the frontend handle 404s for non-API routes if this Go app is purely an API.
	// If it serves frontend assets too, then this NoRoute might be more relevant.
	// For now, commenting out the NoRoute as API should be specific.
	// router.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Endpoint not found"}})
	// })
}
