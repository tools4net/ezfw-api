package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AgentToken holds the schema definition for the AgentToken entity.
type AgentToken struct {
	ent.Schema
}

// Fields of the AgentToken.
func (AgentToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("node_id").NotEmpty(),
		field.String("token").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.String("status").Default("active"), // active, revoked, expired
		field.Time("expires_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("last_used").Optional().Nillable(),
	}
}

// Edges of the AgentToken.
func (AgentToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Ref("agent_tokens").
			Field("node_id").
			Unique().
			Required(),
	}
}

// Indexes of the AgentToken.
func (AgentToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token").Unique(),
		index.Fields("node_id"),
		index.Fields("status"),
		index.Fields("created_at"),
		index.Fields("expires_at"),
	}
}