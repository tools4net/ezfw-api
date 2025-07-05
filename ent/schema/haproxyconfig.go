package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// HAProxyConfig holds the schema definition for the HAProxyConfig entity.
type HAProxyConfig struct {
	ent.Schema
}

// Fields of the HAProxyConfig.
func (HAProxyConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		
		// HAProxy configuration fields as JSON
		field.JSON("global_config", &models.HAProxyGlobal{}).Optional(),
		field.JSON("defaults_config", &models.HAProxyDefaults{}).Optional(),
		field.JSON("frontends", []models.HAProxyFrontend{}).Optional(),
		field.JSON("backends", []models.HAProxyBackend{}).Optional(),
		field.JSON("listens", []models.HAProxyListen{}).Optional(),
		field.JSON("stats_config", &models.HAProxyStats{}).Optional(),
	}
}

// Edges of the HAProxyConfig.
func (HAProxyConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the HAProxyConfig.
func (HAProxyConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("created_at"),
	}
}
