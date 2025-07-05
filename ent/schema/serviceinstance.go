package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ServiceInstance holds the schema definition for the ServiceInstance entity.
type ServiceInstance struct {
	ent.Schema
}

// Fields of the ServiceInstance.
func (ServiceInstance) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("node_id").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.String("service_type").NotEmpty(), // xray, singbox, nginx, wireguard, haproxy
		field.String("status").Default("stopped"), // running, stopped, error, starting, stopping
		field.Int("port").Positive(),
		field.String("protocol").NotEmpty(), // tcp, udp, both
		field.JSON("config", map[string]interface{}{}).Optional(),
		field.JSON("tags", []string{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("last_seen").Optional().Nillable(),
	}
}

// Edges of the ServiceInstance.
func (ServiceInstance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node", Node.Type).
			Ref("service_instances").
			Field("node_id").
			Unique().
			Required(),
	}
}

// Indexes of the ServiceInstance.
func (ServiceInstance) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("node_id"),
		index.Fields("service_type"),
		index.Fields("status"),
		index.Fields("port"),
		index.Fields("created_at"),
		index.Fields("last_seen"),
		index.Fields("node_id", "name").Unique(),
	}
}
