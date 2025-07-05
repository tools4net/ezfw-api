package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tools4net/ezfw/backend/internal/models"
)

// XrayConfig holds the schema definition for the XrayConfig entity.
type XrayConfig struct {
	ent.Schema
}

// Fields of the XrayConfig.
func (XrayConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		
		// Xray configuration fields as JSON
		field.JSON("log", &models.LogObject{}).Optional(),
		field.JSON("api", &models.APIObject{}).Optional(),
		field.JSON("dns", &models.DNSObject{}).Optional(),
		field.JSON("routing", &models.RoutingObject{}).Optional(),
		field.JSON("policy", &models.PolicyObject{}).Optional(),
		field.JSON("inbounds", []models.InboundObject{}).Optional(),
		field.JSON("outbounds", []models.OutboundObject{}).Optional(),
		field.JSON("transport", &models.TransportObject{}).Optional(),
		field.JSON("stats", &models.StatsObject{}).Optional(),
		field.JSON("reverse", &models.ReverseObject{}).Optional(),
		field.JSON("fakedns", &models.FakeDNSObject{}).Optional(),
		field.JSON("metrics", &models.MetricsObject{}).Optional(),
		field.JSON("observatory", &models.ObservatoryObject{}).Optional(),
		field.JSON("burst_observatory", &models.BurstObservatoryObject{}).Optional(),
		field.JSON("services", map[string]interface{}{}).Optional(),
	}
}

// Edges of the XrayConfig.
func (XrayConfig) Edges() []ent.Edge {
	return nil
}

// Indexes of the XrayConfig.
func (XrayConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("created_at"),
	}
}
