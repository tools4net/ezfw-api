package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// Node holds the schema definition for the Node entity.
type Node struct {
	ent.Schema
}

// Fields of the Node.
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.String("hostname").NotEmpty(),
		field.String("ip_address").NotEmpty(),
		field.Int("port").Positive(),
		field.String("status").Default("inactive"), // active, inactive, maintenance, error
		field.JSON("os_info", &models.OSInfo{}).Optional(),
		field.JSON("agent_info", &models.AgentInfo{}).Optional(),
		field.JSON("tags", []string{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("last_seen").Optional().Nillable(),
	}
}

// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("service_instances", ServiceInstance.Type),
		edge.To("agent_tokens", AgentToken.Type),
	}
}

// Indexes of the Node.
func (Node) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("hostname"),
		index.Fields("status"),
		index.Fields("created_at"),
		index.Fields("last_seen"),
	}
}
