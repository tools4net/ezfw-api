package models

import (
	"time"
)

type XrayConfig struct {
	ID               string                      `json:"id"`
	Name             string                      `json:"name"`
	Description      string                      `json:"description"`
	Tag              string                      `json:"tag"`
	Log              *XrayLogConfig              `json:"log"`
	Inbounds         []XrayInbound               `json:"inbounds"`
	Outbounds        []XrayOutbound              `json:"outbounds"`
	API              *XrayAPIConfig              `json:"api"`
	DNS              *XrayDNSConfig              `json:"dns"`
	Reverse          *XrayReverseConfig          `json:"reverse"`
	Routing          *XrayRoutingConfig          `json:"routing"`
	Policy           *XrayPolicyConfig           `json:"policy"`
	Transport        *XrayTransportConfig        `json:"transport"`
	FakeDNS          *XrayFakeDNSConfig          `json:"fakeDNS"`
	Observatory      *XrayObservatoryConfig      `json:"observatory"`
	BurstObservatory *XrayBurstObservatoryConfig `json:"burstObservatory"`
	Stats            *XrayStatsConfig            `json:"stats"`
	Metrics          *XrayMetricsConfig          `json:"metrics"`
	CreatedAt        time.Time                   `json:"createdAt"`
	UpdatedAt        time.Time                   `json:"updatedAt"`
}
type XrayLogConfig struct {
	Loglevel string `json:"level"`
}

type XrayInbound struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	Settings struct {
		Auth string `json:"auth"`
	} `json:"settings"`
	Tag string `json:"tag"`
}
type XrayRoutingConfig struct {
	Rules []*XrayRoutingRule `json:"rules"`
}
type XrayRoutingRule struct {
	Outbound string   `json:"outbound"`
	Domain   []string `json:"domain"`
}

type XrayAPIConfig struct {
	Tag string `json:"tag"`
}
type XrayDNSConfig struct {
	Servers []*XrayDNSServer `json:"servers"`
}
type XrayDNSServer struct {
	Address string `json:"address"`
}
type XrayReverseConfig struct {
	Services []*XrayReverseService `json:"services"`
}
type XrayReverseService struct {
	Tag string `json:"tag"`
}
type XrayPolicyConfig struct {
	Level string `json:"level"`
}
type XrayTransportConfig struct {
	Level string `json:"level"`
}
type XrayFakeDNSConfig struct {
	IPPool string `json:"ipPool"`
}
type XrayObservatoryConfig struct {
	Enabled bool `json:"enabled"`
}
type XrayBurstObservatoryConfig struct {
	Enabled bool `json:"enabled"`
}
type XrayStatsConfig struct {
	Enabled bool `json:"enabled"`
}
type XrayMetricsConfig struct {
	Enabled bool `json:"enabled"`
}

type XrayOutbound struct {
	Protocol string `json:"protocol"`
	Tag      string `json:"tag"`
	Settings struct {
		Vnext []struct {
			Address string `json:"address"`
			Port    int    `json:"port"`
			Users   []struct {
				Id       string `json:"id"`
				AlterId  int    `json:"alterId"`
				Security string `json:"security"`
			} `json:"users"`
		} `json:"vnext"`
	} `json:"settings"`
}

type Certificate struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
