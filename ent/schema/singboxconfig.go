package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// SingBoxConfig holds the schema definition for the SingBoxConfig entity.
type SingBoxConfig struct {
	ent.Schema
}

// Fields of the SingBoxConfig.
func (SingBoxConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		// JSON fields for SingBox configuration components
		field.JSON("log_config", &models.SingBoxLogConfig{}).Optional(),
		field.JSON("dns_config", &models.SingBoxDNSConfig{}).Optional(),
		field.JSON("ntp_config", &models.SingBoxNTPConfig{}).Optional(),
		field.JSON("inbounds", []interface{}{}).Optional(),
		field.JSON("outbounds", []interface{}{}).Optional(),
		field.JSON("route_config", &models.SingBoxRouteConfig{}).Optional(),
		field.JSON("experimental_config", map[string]interface{}{}).Optional(),
		field.JSON("endpoints", []interface{}{}).Optional(),
		field.JSON("certificate_config", map[string]interface{}{}).Optional(),
	}
}

// Edges of the SingBoxConfig.
func (SingBoxConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the SingBoxConfig.
func (SingBoxConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("created_at"),
	}
}
