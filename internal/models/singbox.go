package models

import "time"

type SingBoxLogConfig struct {
Disabled  *bool   `json:"disabled,omitempty"`
Level     *string `json:"level,omitempty"` // e.g., "trace", "debug", "info", "warn", "error", "fatal", "panic"
Output    *string `json:"output,omitempty"` // file path
Timestamp *bool   `json:"timestamp,omitempty"`
}

// SingBoxDialFields encapsulates common dialing options.
type SingBoxDialFields struct {
	Detour              *string     `json:"detour,omitempty"`
	BindInterface       *string     `json:"bind_interface,omitempty"`
	Inet4BindAddress    *string     `json:"inet4_bind_address,omitempty"`
	Inet6BindAddress    *string     `json:"inet6_bind_address,omitempty"`
	RoutingMark         interface{} `json:"routing_mark,omitempty"` // int or hex string
	ReuseAddr           *bool       `json:"reuse_addr,omitempty"`
	NetNS               *string     `json:"netns,omitempty"`
	ConnectTimeout      *string     `json:"connect_timeout,omitempty"`      // duration string
	TCPFastOpen         *bool       `json:"tcp_fast_open,omitempty"`
	TCPMultiPath        *bool       `json:"tcp_multi_path,omitempty"`
	UDPFragment         *bool       `json:"udp_fragment,omitempty"`
	DomainResolver      interface{} `json:"domain_resolver,omitempty"`      // string or object
	NetworkStrategy     *string     `json:"network_strategy,omitempty"`
	NetworkType         []string    `json:"network_type,omitempty"`
	FallbackNetworkType []string    `json:"fallback_network_type,omitempty"`
	FallbackDelay       *string     `json:"fallback_delay,omitempty"` // duration string
	// DomainStrategy is deprecated and removed in favor of DomainResolver
}

// SingBoxECHConfig is a placeholder for ECH configuration.
type SingBoxECHConfig struct {
	Enabled                   *bool    `json:"enabled,omitempty"`
	PQSignatureSchemesEnabled *bool    `json:"pq_signature_schemes_enabled,omitempty"`
	ConfigPath                *string  `json:"config_path,omitempty"`    // Path to ECHKeyContents
	ConfigEncoded             *string  `json:"config_encoded,omitempty"` // Base64 encoded ECHKeyContents
}

// SingBoxDNSTLSSettings defines TLS settings specific to a DNS server.
type SingBoxDNSTLSSettings struct {
	ServerName    *string           `json:"server_name,omitempty"`
	Insecure      *bool             `json:"insecure,omitempty"`
	DisableSNI    *bool             `json:"disable_sni,omitempty"`
	ALPN          []string          `json:"alpn,omitempty"`
	MinVersion    *string           `json:"min_version,omitempty"` // e.g., "tls1.2", "tls1.3"
	MaxVersion    *string           `json:"max_version,omitempty"`
	CipherSuites  []string          `json:"cipher_suites,omitempty"`
	Ech           *SingBoxECHConfig `json:"ech,omitempty"`
	Certificate   []string          `json:"certificate,omitempty"`      // PEM format client certificate chain
	CertificatePath *string         `json:"certificate_path,omitempty"` // Path to client certificate chain file
	Key           *string           `json:"key,omitempty"`              // PEM format client private key
	KeyPath       *string           `json:"key_path,omitempty"`         // Path to client private key file
}

// SingBoxDNSHTTPSettings defines settings for DoH (DNS over HTTPS) or DoH3.
type SingBoxDNSHTTPSettings struct {
	Method  *string           `json:"method,omitempty"` // "GET" or "POST"
	Path    *string           `json:"path,omitempty"`   // Defaults to "/dns-query"
	Headers map[string]string `json:"headers,omitempty"`
}

// SingBoxDNSQUICSettings defines settings for DoQ (DNS over QUIC).
type SingBoxDNSQUICSettings struct {
	InitialPacketLength *int `json:"initial_packet_length,omitempty"` // Default 1200
}

// SingBoxDNSHostsServerSettings defines settings for a "hosts" type DNS server.
type SingBoxDNSHostsServerSettings struct {
	Path   *string                `json:"path,omitempty"`   // Path to hosts file
	Format *string                `json:"format,omitempty"` // "plain_text", "domain_set"
	Data   map[string]interface{} `json:"data,omitempty"`   // Inline hosts data, domain: IP or []IP
}

// SingBoxDNSDhcpServerSettings defines settings for a "dhcp" type DNS server.
type SingBoxDNSDhcpServerSettings struct {
	Interface *string `json:"interface,omitempty"` // "auto" or specific interface name like "en0"
}

// SingBoxDNSServer defines a single DNS server configuration.
// It's a large struct to accommodate various types and their fields as per Sing-box documentation.
type SingBoxDNSServer struct {
	Tag  *string `json:"tag,omitempty"`
	Type *string `json:"type,omitempty"` // "udp", "tcp", "tls", "https", "http3", "quic", "local", "hosts", "dhcp", "fakeip", "resolved", "tailscale", or empty for legacy (deprecated)

	// --- Legacy Fields (used if type is empty or "legacy") ---
	Address         *string `json:"address,omitempty"`          // Legacy server address (e.g., "8.8.8.8", "tls://1.1.1.1")
	AddressResolver *string `json:"address_resolver,omitempty"` // Legacy: tag of another DNS server
	AddressStrategy *string `json:"address_strategy,omitempty"` // Legacy: "prefer_ipv4", "prefer_ipv6", "ipv4_only", "ipv6_only"
	Strategy        *string `json:"strategy,omitempty"`         // Legacy: Default domain strategy for this server. Also a global DNS config field.
	// Detour is part of SingBoxDialFields, ClientSubnet is also part of global DNS config.

	// --- Fields for New Typed Servers (UDP, TCP, TLS, HTTPS, QUIC, HTTP3) ---
	Server     *string `json:"server,omitempty"`      // New: Server address (IP or domain)
	ServerPort *int    `json:"server_port,omitempty"` // New: Server port

	// --- Embedded Dial Fields (common for many new server types) ---
	// Note: Detour from DialFields takes precedence over legacy Detour for new types.
	//       ClientSubnet from DialFields/server-specific context takes precedence.
	SingBoxDialFields

	// --- Type-Specific Settings ---
	TLS   *SingBoxDNSTLSSettings      `json:"tls,omitempty"`   // For types: "tls", "https", "quic", "http3"
	HTTP  *SingBoxDNSHTTPSettings     `json:"http,omitempty"`  // For types: "https", "http3" (DoH/DoH3)
	QUIC  *SingBoxDNSQUICSettings     `json:"quic,omitempty"`  // For types: "quic", "http3" (underlying QUIC transport)
	Hosts *SingBoxDNSHostsServerSettings `json:"hosts,omitempty"` // For type: "hosts"
	DHCP  *SingBoxDNSDhcpServerSettings `json:"dhcp,omitempty"`  // For type: "dhcp"
	// For "local", "fakeip" (as server type), "resolved", "tailscale", they generally don't have extra unique fields beyond Type/Tag and DialFields.
	// "fakeip" as a server type just needs `type: "fakeip"`. It uses the global fakeip config.
	// "local" just needs `type: "local"`.
	// "resolved" just needs `type: "resolved"`.
	// "tailscale" just needs `type: "tailscale"`.

	ClientSubnet *string `json:"client_subnet,omitempty"` // This can be set per-server, overriding global.
}

// SingBoxFakeIPConfig defines FakeIP settings within DNS.
// Documentation: https://sing-box.sagernet.org/configuration/dns/fakeip/
type SingBoxFakeIPConfig struct {
	Enabled            *bool    `json:"enabled,omitempty"`
	Inet4Range         *string  `json:"inet4_range,omitempty"`          // CIDR, e.g., "198.18.0.0/15"
	Inet6Range         *string  `json:"inet6_range,omitempty"`          // CIDR, e.g., "fc00::/18"
	Inet4Mask          *string  `json:"inet4_mask,omitempty"`           // New in 1.9, alternative to inet4_range
	Inet6Mask          *string  `json:"inet6_mask,omitempty"`           // New in 1.9, alternative to inet6_range
	HttpTimeout        *string  `json:"http_timeout,omitempty"`         // Duration string, default 10s
	UDPTimeout         *string  `json:"udp_timeout,omitempty"`          // Duration string, default 10s
	ExcludedDomain     []string `json:"excluded_domain,omitempty"`      // New in 1.9
	ExcludedDomainFile *string  `json:"excluded_domain_file,omitempty"` // New in 1.9, path to file
	Store              *string  `json:"store,omitempty"`                // New in 1.9, "memory" (default) or "bolt"
	StorePath          *string  `json:"store_path,omitempty"`           // New in 1.9, path for "bolt" store
}

// SingBoxDNSRule defines a DNS rule.
// Documentation: https://sing-box.sagernet.org/configuration/dns/rule/
// and Action: https://sing-box.sagernet.org/configuration/dns/rule_action/
type SingBoxDNSRule struct {
	// Action fields (mutually exclusive, only one should be set)
	Server       *string  `json:"server,omitempty"`        // Tag of a DNS server
	Outbound     *string  `json:"outbound,omitempty"`      // Tag of an outbound (becomes a DNS transport)
	Type         *string  `json:"type,omitempty"`          // For "logical" rules: "logical"
	ClientSubnet *string  `json:"client_subnet,omitempty"` // Override client subnet for this rule
	DisableCache *bool    `json:"disable_cache,omitempty"` // Disable DNS cache for this rule
	RewriteTTL   *uint32  `json:"rewrite_ttl,omitempty"`   // Rewrite TTL for matched queries

	// Matching fields
	AuthUser        []string          `json:"auth_user,omitempty"`
	ClashMode       *string           `json:"clash_mode,omitempty"` // "global", "rule", "direct"
	Default         *bool             `json:"default,omitempty"`    // New in 1.9, matches if no other rule matched
	Domain          []string          `json:"domain,omitempty"`
	DomainKeyword   []string          `json:"domain_keyword,omitempty"`
	DomainRegex     []string          `json:"domain_regex,omitempty"`
	DomainSuffix    []string          `json:"domain_suffix,omitempty"`
	Executable      []string          `json:"executable,omitempty"` // New in 1.9
	GeoIP           []string          `json:"geoip,omitempty"`      // Country code
	Geosite         []string          `json:"geosite,omitempty"`
	Inbound         []string          `json:"inbound,omitempty"` // Inbound tags
	Invert          *bool             `json:"invert,omitempty"`
	IPCidr          []string          `json:"ip_cidr,omitempty"`
	IPVersion       *int              `json:"ip_version,omitempty"` // 0, 4, 6
	Network         *string           `json:"network,omitempty"`    // "tcp", "udp"
	PackageName     []string          `json:"package_name,omitempty"`
	Port            []string          `json:"port,omitempty"` // Port or port range string, e.g., "53", "1000-2000"
	PortRange       []string          `json:"port_range,omitempty"` // Deprecated by "port"
	ProcessName     []string          `json:"process_name,omitempty"`
	ProcessPath     []string          `json:"process_path,omitempty"` // New in 1.9
	Protocol        []string          `json:"protocol,omitempty"`     // "tls", "http", "quic", "ssh", "stun"
	QueryType       []string          `json:"query_type,omitempty"`   // "A", "AAAA", "CNAME", "MX", "TXT", "ANY", etc. or numeric
	RuleSet         []string          `json:"rule_set,omitempty"`
	SourceGeoIP     []string          `json:"source_geoip,omitempty"`
	SourceIPCidr    []string          `json:"source_ip_cidr,omitempty"`
	SourcePort      []string          `json:"source_port,omitempty"`      // Port or port range string
	SourcePortRange []string          `json:"source_port_range,omitempty"`// Deprecated by "source_port"
	User            []string          `json:"user,omitempty"`             // New in 1.9
	WIFISSID        []string          `json:"wifi_ssid,omitempty"`        // New in 1.9
	WIFIBSSID       []string          `json:"wifi_bssid,omitempty"`       // New in 1.9

	// For "logical" rules
	Mode  *string           `json:"mode,omitempty"`  // "and" or "or"
	Rules []*SingBoxDNSRule `json:"rules,omitempty"` // Nested rules
}

// SingBoxDNSConfig defines the structure for DNS configurations in SingBox.
type SingBoxDNSConfig struct {
	Servers          []*SingBoxDNSServer    `json:"servers,omitempty"`
	Rules            []*SingBoxDNSRule      `json:"rules,omitempty"`
	Final            *string                `json:"final,omitempty"`   // Tag of final DNS server
	Strategy         *string                `json:"strategy,omitempty"` // Default DNS strategy: "prefer_ipv4", "prefer_ipv6", "ipv4_only", "ipv6_only"
	DisableCache     *bool                  `json:"disable_cache,omitempty"`
	DisableExpire    *bool                  `json:"disable_expire,omitempty"`
	IndependentCache *bool                  `json:"independent_cache,omitempty"`
	ReverseMapping   *bool                  `json:"reverse_mapping,omitempty"`
	FakeIP           *SingBoxFakeIPConfig   `json:"fakeip,omitempty"`
	Hosts            map[string]interface{} `json:"hosts,omitempty"` // domain: IP or IP array
	ClientSubnet     *string                `json:"client_subnet,omitempty"` // Global client_subnet default
	CacheCapacity    *int                   `json:"cache_capacity,omitempty"` // New in 1.11.0
}

// SingBoxNTPConfig defines the structure for NTP configurations in SingBox.
// Documentation: https://sing-box.sagernet.org/configuration/ntp/
type SingBoxNTPConfig struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	Server     *string `json:"server,omitempty"`     // Default: "time.apple.com"
	ServerPort *int    `json:"server_port,omitempty"`// Default: 123
	Interval   *string `json:"interval,omitempty"`   // Duration string, default: "1h"
	Detour     *string `json:"detour,omitempty"`     // Outbound tag for NTP requests (legacy, part of dial_fields now but often shown standalone)

	DialFields *SingBoxDialFields `json:"dial_fields,omitempty"` // New in 1.12 for extended dial options
}

// SingBoxInbound defines the structure for inbound connections in SingBox.
// Common fields from shared/listen: https://sing-box.sagernet.org/configuration/shared/listen/
type SingBoxInbound struct {
	Type                     string                 `json:"type"` // Protocol type, e.g., "mixed", "socks", "http", "vmess", etc.
	Tag                      string                 `json:"tag"`  // Unique tag for this inbound
	Listen                   *string                `json:"listen,omitempty"`
	ListenPort               *int                   `json:"listen_port,omitempty"`
	TCPFastOpen              *bool                  `json:"tcp_fast_open,omitempty"`
	UDPFragment              *bool                  `json:"udp_fragment,omitempty"` // Deprecated in 1.9, use udp_mtu
	UDPMTU                   *int                   `json:"udp_mtu,omitempty"`      // New in 1.9
	Sniff                    *bool                  `json:"sniff,omitempty"`
	SniffOverrideDestination *bool                  `json:"sniff_override_destination,omitempty"`
	SniffTimeout             *string                `json:"sniff_timeout,omitempty"` // Duration string
	DomainStrategy           *string                `json:"domain_strategy,omitempty"` // "prefer_ipv4", "prefer_ipv6", etc.
	UDPTimeout               *string                `json:"udp_timeout,omitempty"`     // Duration string
	ProxyProtocol            *int                   `json:"proxy_protocol,omitempty"`  // PROXY protocol version: 0 (disable), 1, 2
	BindInterface            *string                `json:"bind_interface,omitempty"`
	Inet4BindAddress         *string                `json:"inet4_bind_address,omitempty"`
	Inet6BindAddress         *string                `json:"inet6_bind_address,omitempty"`
	RoutingMark              interface{}            `json:"routing_mark,omitempty"` // int or hex string
	ReuseAddr                *bool                  `json:"reuse_addr,omitempty"`
	NetNS                    *string                `json:"netns,omitempty"` // New in 1.12

	// Protocol-specific settings, TLS, and Transport settings are kept generic for now.
	Settings  map[string]interface{} `json:"settings,omitempty"`  // Protocol-specific settings
	TLS       map[string]interface{} `json:"tls,omitempty"`       // TLS settings object, see shared/tls
	Transport map[string]interface{} `json:"transport,omitempty"` // Transport settings object, see shared/v2ray-transport
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
// Documentation: https://sing-box.sagernet.org/configuration/route/rule/
type SingBoxRouteRule struct {
	// --- Action fields (only one of these should be set per rule) ---
	Outbound *string `json:"outbound,omitempty"` // Tag of an outbound to use
	Balancer *string `json:"balancer,omitempty"` // Tag of a balancer outbound to use (New in 1.9)
	// Type "logical" is implicitly handled by the presence of "rules" and "mode"

	// --- Matching fields ---
	AuthUser      []string    `json:"auth_user,omitempty"`      // New in 1.9
	ClashMode     *string     `json:"clash_mode,omitempty"`     // "global", "rule", "direct"
	Default       *bool       `json:"default,omitempty"`        // New in 1.9, matches if no other preceding rule matched
	Domain        []string    `json:"domain,omitempty"`
	DomainKeyword []string    `json:"domain_keyword,omitempty"`
	DomainRegex   []string    `json:"domain_regex,omitempty"`
	DomainSuffix  []string    `json:"domain_suffix,omitempty"`
	Email         []string    `json:"email,omitempty"`          // New in 1.9
	Executable    []string    `json:"executable,omitempty"`     // New in 1.9
	GeoIP         []string    `json:"geoip,omitempty"`          // Country code, e.g., "CN"
	Geosite       []string    `json:"geosite,omitempty"`
	Inbound       interface{} `json:"inbound,omitempty"`        // string or []string of inbound tags
	Invert        *bool       `json:"invert,omitempty"`
	IPCidr        []string    `json:"ip_cidr,omitempty"`        // Renamed from 'ip'
	IPVersion     *int        `json:"ip_version,omitempty"`     // 0 (any), 4, 6
	Network       interface{} `json:"network,omitempty"`        // string or []string: "tcp", "udp"
	Notes         *string     `json:"notes,omitempty"`          // New in 1.9 (comment field, not for matching)
	PackageName   []string    `json:"package_name,omitempty"`
	Port          interface{} `json:"port,omitempty"`           // string, int, or list of these (e.g., "80", 443, "1000-2000")
	ProcessName   []string    `json:"process_name,omitempty"`
	ProcessPath   []string    `json:"process_path,omitempty"`
	Protocol      interface{} `json:"protocol,omitempty"`       // string or []string: "tls", "http", "quic", "ssh", "stun", etc.
	RuleSet       []string    `json:"rule_set,omitempty"`
	SourceGeoIP   []string    `json:"source_geoip,omitempty"`
	SourceIPCidr  []string    `json:"source_ip_cidr,omitempty"` // Renamed from 'source_ip'
	SourcePort    interface{} `json:"source_port,omitempty"`    // string, int, or list of these
	User          []string    `json:"user,omitempty"`           // New in 1.9
	WIFIBSSID     []string    `json:"wifi_bssid,omitempty"`     // New in 1.9
	WIFISSID      []string    `json:"wifi_ssid,omitempty"`      // New in 1.9

	// --- For "logical" rules ---
	Type  *string             `json:"type,omitempty"`  // Should be "logical" if Mode and Rules are set
	Mode  *string             `json:"mode,omitempty"`  // "and" or "or"
	Rules []*SingBoxRouteRule `json:"rules,omitempty"` // Nested rules
}

// SingBoxRouteConfig defines the structure for routing configurations in SingBox.
// Documentation: https://sing-box.sagernet.org/configuration/route/
type SingBoxRouteConfig struct {
	Rules                 []*SingBoxRouteRule      `json:"rules,omitempty"`
	RuleSet               []map[string]interface{} `json:"rule_set,omitempty"`         // List of RuleSet objects
	Final                 *string                  `json:"final,omitempty"`              // Default outbound tag
	AutoDetectInterface   *bool                    `json:"auto_detect_interface,omitempty"`// Default: true
	OverrideAndroidVPN    *bool                    `json:"override_android_vpn,omitempty"` // Default: true
	DefaultInterface      *string                  `json:"default_interface,omitempty"`    // Default physical outbound interface
	DefaultMark           *int                     `json:"default_mark,omitempty"`         // SO_MARK for outbound packets
	GeoIP                 *map[string]interface{}  `json:"geoip,omitempty"`                // GeoIP configuration object
	Geosite               *map[string]interface{}  `json:"geosite,omitempty"`              // Geosite configuration object
	DomainStrategy        *string                  `json:"domain_strategy,omitempty"`      // "prefer_ipv4", "prefer_ipv6", "ipv4_only", "ipv6_only", "asio" (deprecated), UseIP (Xray term?)
	IndependentCache      *bool                    `json:"independent_cache,omitempty"`    // For rule_set matching result cache (New in 1.9)
	DefaultDomainResolver *string                  `json:"default_domain_resolver,omitempty"`// Tag of a DNS server (New in 1.12)
}

// SingBoxConfig is the main configuration structure for SingBox, managed by ProxyPanel.
type SingBoxConfig struct {
	ID          string    `json:"id,omitempty" gorm:"primaryKey" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	Name        string    `json:"name,omitempty" example:"My SingBox Test"`
	Description string    `json:"description,omitempty" example:"Experimental Sing-box setup"`
	CreatedAt   time.Time `json:"createdAt,omitempty" example:"2023-01-02T10:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" example:"2023-01-02T11:00:00Z"`

	Log          *SingBoxLogConfig         `json:"log,omitempty"`
	DNS          *SingBoxDNSConfig         `json:"dns,omitempty"`
	NTP          *SingBoxNTPConfig         `json:"ntp,omitempty"`
	Inbounds     []*SingBoxInbound         `json:"inbounds,omitempty"`
	Outbounds    []*SingBoxOutbound        `json:"outbounds,omitempty"`
	Route        *SingBoxRouteConfig       `json:"route,omitempty"`
	Experimental *map[string]interface{}   `json:"experimental,omitempty"`
	Services     []map[string]interface{}  `json:"services,omitempty"`     // Generic map for various service types
	Endpoints    []map[string]interface{}  `json:"endpoints,omitempty"`    // Generic map for various endpoint types
	Certificate  []*SingBoxCertificate     `json:"certificate,omitempty"`  // List of certificate objects
}

// SingBoxCertificate defines a certificate for use in TLS configurations.
// Documentation: https://sing-box.sagernet.org/configuration/certificate/
type SingBoxCertificate struct {
	Certificate     []string `json:"certificate,omitempty"`      // PEM format certificate chain (list of strings)
	CertificatePath *string  `json:"certificate_path,omitempty"` // Path to PEM certificate chain file
	Key             []string `json:"key,omitempty"`              // PEM format private key (list of strings)
	KeyPath         *string  `json:"key_path,omitempty"`         // Path to PEM private key file
}