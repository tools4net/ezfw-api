package store

import (
	"context"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// Store defines the interface for database operations.
type Store interface {
	// SingBox Configuration methods
	CreateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error
	GetSingBoxConfig(ctx context.Context, id string) (*models.SingBoxConfig, error)
	ListSingBoxConfigs(ctx context.Context, limit, offset int) ([]*models.SingBoxConfig, error)
	UpdateSingBoxConfig(ctx context.Context, config *models.SingBoxConfig) error
	DeleteSingBoxConfig(ctx context.Context, id string) error
	// CountSingBoxConfigs(ctx context.Context) (int, error) // Optional: for pagination metadata

	// Xray Configuration methods
	CreateXrayConfig(ctx context.Context, config *models.XrayConfig) error
	GetXrayConfig(ctx context.Context, id string) (*models.XrayConfig, error)
	ListXrayConfigs(ctx context.Context, limit, offset int) ([]*models.XrayConfig, error)
	UpdateXrayConfig(ctx context.Context, config *models.XrayConfig) error
	DeleteXrayConfig(ctx context.Context, id string) error
	// CountXrayConfigs(ctx context.Context) (int, error) // Optional: for pagination metadata
}
