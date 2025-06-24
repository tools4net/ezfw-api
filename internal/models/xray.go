package models

import "time"

// XrayConfig is the top-level structure for an Xray configuration.
// It also includes metadata for storage and management within ProxyPanel.
type XrayConfig struct {
	ID          string    `json:"id" gorm:"primaryKey"` // Internal ID for database
	Name        string    `json:"name" gorm:"uniqueIndex"` // User-defined name for the config
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Core Xray configuration fields
	Log              *LogObject              `json:"log,omitempty"`
	API              *APIObject              `json:"api,omitempty"`
	DNS              *DNSObject              `json:"dns,omitempty"`
	Routing          *RoutingObject          `json:"routing,omitempty"`
	Policy           *PolicyObject           `json:"policy,omitempty"`
	Inbounds         []InboundObject         `json:"inbounds,omitempty"`
	Outbounds        []OutboundObject        `json:"outbounds,omitempty"`
	Transport        *TransportObject        `json:"transport,omitempty"`
	Stats            *StatsObject            `json:"stats,omitempty"`
	Reverse          *ReverseObject          `json:"reverse,omitempty"`
	FakeDNS          *FakeDNSObject          `json:"fakedns,omitempty"` // Corrected casing
	Metrics          *MetricsObject          `json:"metrics,omitempty"`
	Observatory      *ObservatoryObject      `json:"observatory,omitempty"`
	BurstObservatory *BurstObservatoryObject `json:"burstObservatory,omitempty"` // Corrected casing
}

// LogObject defines logging settings.
type LogObject struct {
	Access   string `json:"access,omitempty"`
	Error    string `json:"error,omitempty"`
	Loglevel string `json:"loglevel,omitempty"` // Typically "debug", "info", "warning", "error", "none"
}

// APIObject defines settings for the Xray API.
type APIObject struct {
	Tag      string   `json:"tag,omitempty"`
	Services []string `json:"services,omitempty"` // e.g., "HandlerService", "LoggerService", "StatsService"
}

// DNSObject defines DNS settings.
type DNSObject struct {
	Hosts    map[string]interface{} `json:"hosts,omitempty"` // Can be string or array of strings
	Servers  []interface{}          `json:"servers,omitempty"` // Can be string (IP) or DnsServerObject
	ClientIP string                 `json:"clientIp,omitempty"`
	Tag      string                 `json:"tag,omitempty"`
	// Other fields like queryStrategy, disableCache, disableFallback etc. can be added
}

// DnsServerObject is used when a server in DNSObject is an object.
type DnsServerObject struct {
	Address      string   `json:"address,omitempty"`
	Port         int      `json:"port,omitempty"`
	Domains      []string `json:"domains,omitempty"`
	ExpectIPs    []string `json:"expectIPs,omitempty"`
	SkipFallback bool     `json:"skipFallback,omitempty"`
}

// RoutingObject defines routing rules.
type RoutingObject struct {
	DomainStrategy string        `json:"domainStrategy,omitempty"` // "AsIs", "IPIfNonMatch", "IPOnDemand"
	Rules          []RoutingRule `json:"rules,omitempty"`
	Balancers      []Balancer    `json:"balancers,omitempty"`
}

// RoutingRule defines a single routing rule.
type RoutingRule struct {
	Type        string   `json:"type"` // "field"
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag,omitempty"`
	Domain      []string `json:"domain,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Port        string   `json:"port,omitempty"` // e.g., "53", "1000-2000"
	Network     string   `json:"network,omitempty"` // "tcp", "udp", "tcp,udp"
	Protocol    []string `json:"protocol,omitempty"` // e.g., "http", "tls", "bittorrent"
	// Other fields like attrs, user, source, etc.
}

// Balancer defines a load balancer.
type Balancer struct {
	Tag      string   `json:"tag"`
	Selector []string `json:"selector,omitempty"`
	// Strategy object can be added if needed
}

// PolicyObject defines local policy settings.
type PolicyObject struct {
	Levels map[string]LevelPolicy `json:"levels,omitempty"` // Key is level (e.g., "0")
	System *SystemPolicy          `json:"system,omitempty"`
}

// LevelPolicy defines policy for a specific user level.
type LevelPolicy struct {
	Handshake      *int `json:"handshake,omitempty"`      // Duration in seconds
	ConnIdle       *int `json:"connIdle,omitempty"`       // Duration in seconds
	UplinkOnly     *int `json:"uplinkOnly,omitempty"`     // Duration in seconds
	DownlinkOnly   *int `json:"downlinkOnly,omitempty"`   // Duration in seconds
	StatsUserUplink bool `json:"statsUserUplink,omitempty"`
	StatsUserDownlink bool `json:"statsUserDownlink,omitempty"`
	BufferSize     *int `json:"bufferSize,omitempty"` // Buffer size in KB, 0 for default
}

// SystemPolicy defines system-wide policies.
type SystemPolicy struct {
	StatsInboundUplink    bool `json:"statsInboundUplink,omitempty"`
	StatsInboundDownlink  bool `json:"statsInboundDownlink,omitempty"`
	StatsOutboundUplink   bool `json:"statsOutboundUplink,omitempty"`
	StatsOutboundDownlink bool `json:"statsOutboundDownlink,omitempty"`
}

// InboundObject defines an inbound connection handler.
type InboundObject struct {
	Tag            string                 `json:"tag,omitempty"`
	Listen         string                 `json:"listen,omitempty"`
	Port           interface{}            `json:"port,omitempty"` // int or string like "1000-2000"
	Protocol       string                 `json:"protocol"`       // e.g., "vmess", "vless", "trojan", "socks", "http"
	Settings       map[string]interface{} `json:"settings,omitempty"` // Protocol-specific settings
	StreamSettings *StreamSettingsObject  `json:"streamSettings,omitempty"`
	Sniffing       *SniffingObject        `json:"sniffing,omitempty"`
	Allocate       *AllocateObject        `json:"allocate,omitempty"`
}

// OutboundObject defines an outbound connection handler.
type OutboundObject struct {
	Tag            string                 `json:"tag,omitempty"`
	Protocol       string                 `json:"protocol"` // e.g., "vmess", "vless", "freedom", "blackhole"
	Settings       map[string]interface{} `json:"settings,omitempty"` // Protocol-specific settings
	StreamSettings *StreamSettingsObject  `json:"streamSettings,omitempty"`
	ProxySettings  *ProxySettings         `json:"proxySettings,omitempty"` // For daisy-chaining
	SendThrough    string                 `json:"sendThrough,omitempty"`
	Mux            *MuxObject             `json:"mux,omitempty"`
}

// StreamSettingsObject defines transport settings.
type StreamSettingsObject struct {
	Network      string             `json:"network,omitempty"` // "tcp", "kcp", "ws", "http", "quic", "grpc"
	Security     string             `json:"security,omitempty"`  // "none", "tls", "xtls"
	TLSSettings  *TLSSettings       `json:"tlsSettings,omitempty"`
	XTLSSettings *XTLSSettings      `json:"xtlsSettings,omitempty"` // Corrected casing
	TCPSettings  *TCPSettings       `json:"tcpSettings,omitempty"`
	KCPSettings  *KCPSettings       `json:"kcpSettings,omitempty"`
	WSSettings   *WSSettings        `json:"wsSettings,omitempty"`
	HTTPSettings *HTTP2Settings     `json:"httpSettings,omitempty"` // For HTTP/2
	QUICSettings *QUICSettings      `json:"quicSettings,omitempty"`
	GRPCSettings *GRPCSettings      `json:"grpcSettings,omitempty"`
	SocketSettings *SocketOptions  `json:"sockopt,omitempty"`
}

// TLSSettings defines TLS settings.
type TLSSettings struct {
	ServerName         string         `json:"serverName,omitempty"`
	AllowInsecure      bool           `json:"allowInsecure,omitempty"`
	ALPN               []string       `json:"alpn,omitempty"`
	Certificates       []Certificate  `json:"certificates,omitempty"`
	DisableSystemRoot  bool           `json:"disableSystemRoot,omitempty"`
	EnableSessionResumption bool		  `json:"enableSessionResumption,omitempty"`
	Fingerprint        string         `json:"fingerprint,omitempty"` // "chrome", "firefox", "safari", "ios", "random"
	// Other fields like minVersion, maxVersion, cipherSuites etc.
}

// XTLSSettings defines XTLS settings.
type XTLSSettings struct {
	ServerName   string        `json:"serverName,omitempty"`
	AllowInsecure bool          `json:"allowInsecure,omitempty"`
	ALPN         []string      `json:"alpn,omitempty"`
	Certificates []Certificate `json:"certificates,omitempty"`
	MinVersion   string        `json:"minVersion,omitempty"`
	MaxVersion   string        `json:"maxVersion,omitempty"`
	CipherSuites string        `json:"cipherSuites,omitempty"`
}


// Certificate defines a TLS certificate.
type Certificate struct {
	Usage         string `json:"usage,omitempty"` // "encipherment", "verify", "issue"
	CertificateFile string `json:"certificateFile,omitempty"`
	KeyFile         string `json:"keyFile,omitempty"`
	Certificate     []string `json:"certificate,omitempty"` // Certificate content as string array
	Key             []string `json:"key,omitempty"`         // Key content as string array
	OcspStapling    uint64   `json:"ocspStapling,omitempty"`
}

// TCPSettings specific configurations for TCP transport.
type TCPSettings struct {
	Header          *HeaderObject `json:"header,omitempty"`
	AcceptProxyProtocol bool       `json:"acceptProxyProtocol,omitempty"`
}

// HeaderObject for TCP header obfuscation.
type HeaderObject struct {
	Type    string      `json:"type"` // "none" or "http"
	Request *RequestConfig `json:"request,omitempty"` // Used if type is "http"
	Response *ResponseConfig `json:"response,omitempty"` // Used if type is "http"
}

// RequestConfig for HTTP header
type RequestConfig struct {
	Version string `json:"version,omitempty"`
	Method string `json:"method,omitempty"`
	Path []string `json:"path,omitempty"`
	Headers map[string]interface{} `json:"headers,omitempty"` // Value can be string or []string
}

// ResponseConfig for HTTP header
type ResponseConfig struct {
	Version string `json:"version,omitempty"`
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
	Headers map[string]interface{} `json:"headers,omitempty"` // Value can be string or []string
}

// KCPSettings specific configurations for mKCP transport.
type KCPSettings struct {
	Mtu             *int          `json:"mtu,omitempty"`
	Tti             *int          `json:"tti,omitempty"`
	UplinkCapacity  *int          `json:"uplinkCapacity,omitempty"`
	DownlinkCapacity *int          `json:"downlinkCapacity,omitempty"`
	Congestion      bool          `json:"congestion,omitempty"`
	ReadBufferSize  *int          `json:"readBufferSize,omitempty"`
	WriteBufferSize *int          `json:"writeBufferSize,omitempty"`
	Header          *HeaderObject `json:"header,omitempty"` // Same HeaderObject as TCP
	Seed            string        `json:"seed,omitempty"`
}

// WSSettings specific configurations for WebSocket transport.
type WSSettings struct {
	Path                string            `json:"path,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol,omitempty"`
	MaxEarlyData        int32             `json:"maxEarlyData,omitempty"`
	EarlyDataHeaderName string            `json:"earlyDataHeaderName,omitempty"`
}

// HTTP2Settings (H2) specific configurations.
type HTTP2Settings struct {
	Host []string `json:"host,omitempty"`
	Path string   `json:"path,omitempty"`
	// Other fields like readIdleTimeout, healthCheckTimeout
}

// QUICSettings specific configurations for QUIC transport.
type QUICSettings struct {
	Security       string        `json:"security,omitempty"` // "none", "aes-128-gcm", "chacha20-poly1305"
	Key            string        `json:"key,omitempty"`
	Header         *HeaderObject `json:"header,omitempty"` // Same HeaderObject as TCP
	KeyFile        string        `json:"keyFile,omitempty"`
	CertFile       string        `json:"certFile,omitempty"`
}

// GRPCSettings specific configurations for gRPC transport.
type GRPCSettings struct {
	ServiceName         string `json:"serviceName,omitempty"`
	MultiMode           bool   `json:"multiMode,omitempty"`
	IdleTimeout         int    `json:"idle_timeout,omitempty"`
	HealthCheckTimeout  int    `json:"health_check_timeout,omitempty"`
	PermitWithoutStream bool   `json:"permit_without_stream,omitempty"`
	InitialWindowsSize  int    `json:"initial_windows_size,omitempty"`
}

// SocketOptions for various transport protocols.
type SocketOptions struct {
	Mark        *int   `json:"mark,omitempty"`         // SO_MARK
	Tfo         *int   `json:"tcpFastOpen,omitempty"`  // TCP Fast Open queue length
	Tproxy      string `json:"tproxy,omitempty"`       // "redirect" or "tproxy"
	DomainStrategy string `json:"domainStrategy,omitempty"` // "UseIP", "UseIPv4", "UseIPv6"
	DialerProxy string `json:"dialerProxy,omitempty"`  // Tag of an outbound proxy
	AcceptProxyProtocol bool `json:"acceptProxyProtocol,omitempty"`
	TCPKeepAliveInterval *int `json:"tcpKeepAliveInterval,omitempty"` // TCP_KEEPINTVL
	TCPKeepAliveIdle *int `json:"tcpKeepAliveIdle,omitempty"` // TCP_KEEPIDLE
	TCPCongestion string `json:"tcpCongestion,omitempty"` // TCP congestion control algorithm
	Interface   string `json:"interface,omitempty"`    // Bind to a specific network interface
	V6Only      bool   `json:"v6only,omitempty"`       // IPV6_V6ONLY
	TCPMptcp    bool   `json:"tcpMptcp,omitempty"`     // TCP MPTCP
	TCPNoDelay  bool   `json:"tcpNoDelay,omitempty"`   // TCP_NODELAY
}


// SniffingObject defines settings for content sniffing.
type SniffingObject struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride,omitempty"` // "http", "tls", "fakedns"
	DomainsExcluded []string `json:"domainsExcluded,omitempty"`
	RouteOnly    bool     `json:"routeOnly,omitempty"`
}

// AllocateObject defines port allocation strategy.
type AllocateObject struct {
	Strategy    string `json:"strategy,omitempty"` // "always" or "random"
	Refresh     *int   `json:"refresh,omitempty"`  // In minutes
	Concurrency *int   `json:"concurrency,omitempty"`
}

// ProxySettings for outbound proxying.
type ProxySettings struct {
	Tag         string `json:"tag"`
	TransportLayer bool   `json:"transportLayer,omitempty"`
}

// MuxObject for connection multiplexing.
type MuxObject struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency,omitempty"` // 1-1024, default 8
	XUDPConcurrency int `json:"xudpConcurrency,omitempty"`
	XUDPProxyUDP443 string `json:"xudpProxyUDP443,omitempty"` // "disabled", "skip", "prefer"
}


// TransportObject defines global transport settings.
type TransportObject struct {
	TCPSettings  *TCPSettings  `json:"tcpSettings,omitempty"`
	KCPSettings  *KCPSettings  `json:"kcpSettings,omitempty"`
	WSSettings   *WSSettings   `json:"wsSettings,omitempty"`
	HTTPSettings *HTTP2Settings  `json:"httpSettings,omitempty"` // For H2
	QUICSettings *QUICSettings `json:"quicSettings,omitempty"`
	DSSettings   *DomainSocketSettings `json:"dsSettings,omitempty"` // Domain Socket
	GRPCSettings *GRPCSettings `json:"grpcSettings,omitempty"`
	SocketSettings *SocketOptions `json:"sockopt,omitempty"`
}

// DomainSocketSettings for Domain Socket transport.
type DomainSocketSettings struct {
	Path string `json:"path"`
	Abstract bool `json:"abstract,omitempty"`
	Padding  bool `json:"padding,omitempty"`
}


// StatsObject enables or disables statistics. (Empty object means enabled)
type StatsObject struct{}

// ReverseObject defines reverse proxy settings.
type ReverseObject struct {
	Bridges []Bridge `json:"bridges,omitempty"`
	Portals []Portal `json:"portals,omitempty"`
}

// Bridge for reverse proxy.
type Bridge struct {
	Tag    string `json:"tag"`
	Domain string `json:"domain"`
}

// Portal for reverse proxy.
type Portal struct {
	Tag    string `json:"tag"`
	Domain string `json:"domain"`
}

// FakeDNSObject enables FakeDNS.
type FakeDNSObject struct {
	IPPool  string `json:"ipPool"`  // CIDR format, e.g., "198.18.0.0/15"
	PoolSize *int  `json:"poolSize,omitempty"` // Default 65535
}

// MetricsObject defines metrics settings.
type MetricsObject struct {
	Tag    string `json:"tag,omitempty"`
	Listen string `json:"listen,omitempty"`
	Port   *int   `json:"port,omitempty"`
}

// ObservatoryObject defines connection observation settings.
type ObservatoryObject struct {
	SubjectSelector []string `json:"subjectSelector,omitempty"`
	ProbeURL        string   `json:"probeURL,omitempty"`
	ProbeInterval   string   `json:"probeInterval,omitempty"` // e.g. "10m", "1h"
}

// BurstObservatoryObject defines settings for burst connection observation.
type BurstObservatoryObject struct {
	SubjectSelector []string `json:"subjectSelector,omitempty"`
	ProbeURL        string   `json:"probeURL,omitempty"`
	ProbeInterval   string   `json:"probeInterval,omitempty"`
}

// Helper function to get a pointer to an int. Useful for optional int fields.
func IntPtr(i int) *int {
	return &i
}

// GenericConfig is a placeholder for actual configuration data (Xray or Singbox)
// This is used by the handlers to abstract the specific config type initially.
// The actual XrayConfig or SingboxConfig will be stored as JSON string in the DB.
type GenericConfig struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"` // "xray" or "singbox"
	ConfigJSON  string    `json:"config_json"` // The actual JSON configuration
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
