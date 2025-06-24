package models

import "time"

type SingBoxLogConfig struct {
Disabled  *bool   `json:"disabled,omitempty"`
Level     *string `json:"level,omitempty"` // e.g., "trace", "debug", "info", "warn", "error", "fatal", "panic"
Output    *string `json:"output,omitempty"` // file path
Timestamp *bool   `json:"timestamp,omitempty"`
// TODO: Add other log options as per SingBox documentation if needed.
}

// SingBoxDNSServer defines a single DNS server configuration.
type SingBoxDNSServer struct {
	Tag             *string  `json:"tag,omitempty"`
	Address         *string  `json:"address,omitempty"`
	AddressResolver *string  `json:"address_resolver,omitempty"`
	AddressStrategy *string  `json:"address_strategy,omitempty"` // prefer_ipv4, prefer_ipv6, ipv4_only, ipv6_only
	Strategy        *string  `json:"strategy,omitempty"`         // #dns_strategy
	Detour          *string  `json:"detour,omitempty"`
	ClientSubnet    *string  `json:"client_subnet,omitempty"`
}

// SingBoxFakeIPConfig defines FakeIP settings within DNS.
type SingBoxFakeIPConfig struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	Inet4Range *string `json:"inet4_range,omitempty"` // CIDR
	Inet6Range *string `json:"inet6_range,omitempty"` // CIDR
	// ... other FakeIP options
}

// SingBoxDNSRule defines a simplified DNS rule.
type SingBoxDNSRule struct {
	Type        *string  `json:"type,omitempty"`
	Outbound    []string `json:"outbound,omitempty"`
	Server      *string  `json:"server,omitempty"`
	Domains     []string `json:"domains,omitempty"`
	IPVersion   *int     `json:"ip_version,omitempty"`
	Mode        *string  `json:"mode,omitempty"` // "and", "or" for logical rules
	Rules       []*SingBoxDNSRule `json:"rules,omitempty"` // Nested rules for logical type
	// Other matching fields
	DomainSuffix    []string `json:"domain_suffix,omitempty"`
	DomainKeyword   []string `json:"domain_keyword,omitempty"`
	DomainRegex     []string `json:"domain_regex,omitempty"`
	Geosite         []string `json:"geosite,omitempty"`
	SourceGeoIP     []string `json:"source_geoip,omitempty"`
	GeoIP           []string `json:"geoip,omitempty"`
	IPCidr          []string `json:"ip_cidr,omitempty"`
	SourceIPCidr    []string `json:"source_ip_cidr,omitempty"`
	Port            []string `json:"port,omitempty"` // Port or port range string
	ProcessName     []string `json:"process_name,omitempty"`
	PackageName     []string `json:"package_name,omitempty"`
	Protocol        []string `json:"protocol,omitempty"`
	QueryType       []string `json:"query_type,omitempty"` // A, AAAA, CNAME, etc.
	Network         *string  `json:"network,omitempty"` // "tcp", "udp"
	ClashMode       *string  `json:"clash_mode,omitempty"`
	Invert          *bool    `json:"invert,omitempty"`
	DisableCache    *bool    `json:"disable_cache,omitempty"`
	RewriteTTL      *uint32  `json:"rewrite_ttl,omitempty"`
}


// SingBoxDNSConfig defines the structure for DNS configurations in SingBox.
type SingBoxDNSConfig struct {
	Servers          []*SingBoxDNSServer     `json:"servers,omitempty"`
	Rules            []*SingBoxDNSRule       `json:"rules,omitempty"`
	Final            *string                 `json:"final,omitempty"`   // tag of final DNS server
	Strategy         *string                 `json:"strategy,omitempty"` // #dns_strategy
	DisableCache     *bool                   `json:"disable_cache,omitempty"`
	DisableExpire    *bool                   `json:"disable_expire,omitempty"`
	IndependentCache *bool                   `json:"independent_cache,omitempty"`
	ReverseMapping   *bool                   `json:"reverse_mapping,omitempty"`
	FakeIP           *SingBoxFakeIPConfig    `json:"fakeip,omitempty"`
	Hosts            map[string]interface{}  `json:"hosts,omitempty"` // domain: IP or IP array
}

// SingBoxNTPConfig defines the structure for NTP configurations in SingBox.
type SingBoxNTPConfig struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	Server     *string `json:"server,omitempty"`
	ServerPort *int    `json:"server_port,omitempty"`
	Interval   *string `json:"interval,omitempty"` // duration string
	Detour     *string `json:"detour,omitempty"`   // outbound tag
}

// SingBoxInbound defines the structure for inbound connections in SingBox.
type SingBoxInbound struct {
	Type                     string                 `json:"type"`
	Tag                      string                 `json:"tag"`
	Listen                   *string                `json:"listen,omitempty"`
	ListenPort               *int                   `json:"listen_port,omitempty"`
	TCPFastOpen              *bool                  `json:"tcp_fast_open,omitempty"`
	UDPFragment              *bool                  `json:"udp_fragment,omitempty"` // Deprecated, use udp_mtu
	UDPMTU                   *int                   `json:"udp_mtu,omitempty"`
	Sniff                    *bool                  `json:"sniff,omitempty"`
	SniffOverrideDestination *bool                  `json:"sniff_override_destination,omitempty"`
	SniffTimeout             *string                `json:"sniff_timeout,omitempty"` // duration string
	DomainStrategy           *string                `json:"domain_strategy,omitempty"`
	UDPTimeout               *string                `json:"udp_timeout,omitempty"` // duration string
	ProxyProtocol            *int                   `json:"proxy_protocol,omitempty"` // 0, 1, 2
	Settings                 map[string]interface{} `json:"settings,omitempty"` // Protocol-specific settings
	TLS                      map[string]interface{} `json:"tls,omitempty"`      // TLS settings object
	Transport                map[string]interface{} `json:"transport,omitempty"`// Transport settings object
}

// SingBoxOutbound defines the structure for outbound connections in SingBox.
type SingBoxOutbound struct {
	Type     string                 `json:"type"`
	Tag      string                 `json:"tag"`
	Settings map[string]interface{} `json:"settings,omitempty"`   // Protocol-specific settings
	TLS      map[string]interface{} `json:"tls,omitempty"`        // TLS settings object
	Transport map[string]interface{}`json:"transport,omitempty"` // Transport settings object
	Multiplex map[string]interface{} `json:"multiplex,omitempty"`// Multiplex settings object
}

// SingBoxRouteRule defines a rule within the routing configuration.
type SingBoxRouteRule struct {
	// Matching fields
	Inbound           interface{} `json:"inbound,omitempty"`           // string or []string
	IPVersion         *int        `json:"ip_version,omitempty"`        // 0, 4, 6
	Network           *string     `json:"network,omitempty"`           // "tcp", "udp"
	Domain            []string    `json:"domain,omitempty"`
	DomainSuffix      []string    `json:"domain_suffix,omitempty"`
	DomainKeyword     []string    `json:"domain_keyword,omitempty"`
	DomainRegex       []string    `json:"domain_regex,omitempty"`
	Geosite           []string    `json:"geosite,omitempty"`
	SourceGeoIP       []string    `json:"source_geoip,omitempty"`
	GeoIP             []string    `json:"geoip,omitempty"`
	SourceIPCidr      []string    `json:"source_ip_cidr,omitempty"`    // Renamed from source_ip
	IPCidr            []string    `json:"ip_cidr,omitempty"`           // Renamed from ip
	SourcePort        []string    `json:"source_port,omitempty"`       // Port or port range string
	Port              []string    `json:"port,omitempty"`              // Port or port range string
	ProcessName       []string    `json:"process_name,omitempty"`
	ProcessPath       []string    `json:"process_path,omitempty"`
	PackageName       []string    `json:"package_name,omitempty"`
	Protocol          []string    `json:"protocol,omitempty"`          // "tcp", "udp", "tls", "quic", "hysteria", "http"
	ClashMode         *string     `json:"clash_mode,omitempty"`        // "global", "rule", "direct"
	WifiSSID          []string    `json:"wifi_ssid,omitempty"`
	WifiBSSID         []string    `json:"wifi_bssid,omitempty"`
	RuleSet           []string    `json:"rule_set,omitempty"`
	Invert            *bool       `json:"invert,omitempty"`
	// Action fields
	Outbound          *string     `json:"outbound,omitempty"`          // Target outbound tag
	Balancer          *string     `json:"balancer,omitempty"`
	// For "logical" rules
	Type              *string     `json:"type,omitempty"` // "logical"
	Mode              *string     `json:"mode,omitempty"` // "and", "or"
	Rules             []*SingBoxRouteRule `json:"rules,omitempty"` // Nested rules
}

// SingBoxRouteConfig defines the structure for routing configurations in SingBox.
type SingBoxRouteConfig struct {
	Rules               []*SingBoxRouteRule    `json:"rules,omitempty"`
	RuleSet             []map[string]interface{} `json:"rule_set,omitempty"` // RuleSet object
	Final               *string                `json:"final,omitempty"`
	AutoDetectInterface *bool                  `json:"auto_detect_interface,omitempty"`
	OverrideAndroidVPN  *bool                  `json:"override_android_vpn,omitempty"`
	DefaultInterface    *string                `json:"default_interface,omitempty"`
	DefaultMark         *int                   `json:"default_mark,omitempty"`
	GeoIP               *map[string]interface{} `json:"geoip,omitempty"`    // GeoIP object
	Geosite             *map[string]interface{} `json:"geosite,omitempty"`  // Geosite object
	DomainStrategy      *string                `json:"domain_strategy,omitempty"` // prefer_ipv4, prefer_ipv6, ipv4_only, ipv6_only
}

// SingBoxConfig is the main configuration structure for SingBox, managed by ProxyPanel.
type SingBoxConfig struct {
	ID          string    `json:"id,omitempty" gorm:"primaryKey"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`

	Log          *SingBoxLogConfig         `json:"log,omitempty"`
	DNS          *SingBoxDNSConfig         `json:"dns,omitempty"`
	NTP          *SingBoxNTPConfig         `json:"ntp,omitempty"`
	Inbounds     []*SingBoxInbound         `json:"inbounds,omitempty"`
	Outbounds    []*SingBoxOutbound        `json:"outbounds,omitempty"`
	Route        *SingBoxRouteConfig       `json:"route,omitempty"`
	Experimental *map[string]interface{}   `json:"experimental,omitempty"`
	Services     []map[string]interface{}  `json:"services,omitempty"`     // Changed to map from []*map
	Endpoints    []map[string]interface{}  `json:"endpoints,omitempty"`    // Changed to map from []*map
	Certificate  []map[string]interface{}  `json:"certificate,omitempty"`  // Changed to map from []*map
}