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
	FakeDNS          *FakeDNSObject          `json:"fakedns,omitempty"`
	Metrics          *MetricsObject          `json:"metrics,omitempty"`
	Observatory      *ObservatoryObject      `json:"observatory,omitempty"`
	BurstObservatory *BurstObservatoryObject `json:"burstObservatory,omitempty"` // Project X specific
	Services         map[string]interface{}  `json:"services,omitempty"`         // New, for pluggable services
}

// LogObject defines logging settings.
// Docs: https://xtls.github.io/config/log.html
type LogObject struct {
	Access   *string `json:"access,omitempty"`   // Path to access log file
	Error    *string `json:"error,omitempty"`    // Path to error log file
	Loglevel *string `json:"loglevel,omitempty"` // "debug", "info", "warning", "error", "none"
	DNSLog   *bool   `json:"dnsLog,omitempty"`   // Enable DNS log
}

// APIObject defines settings for the Xray API.
// Docs: https://xtls.github.io/config/api.html
type APIObject struct {
	Tag      *string  `json:"tag,omitempty"`      // Tag for routing requests to this API
	Services []string `json:"services,omitempty"` // List of API services to enable, e.g., "HandlerService", "StatsService"
	Listen   *string  `json:"listen,omitempty"`   // Address to listen on for gRPC (optional)
}

// DNSObject defines DNS settings.
// Docs: https://xtls.github.io/config/dns.html
type DNSObject struct {
	Hosts                  map[string]interface{} `json:"hosts,omitempty"` // Static host mappings. Value can be string IP, []string IPs, or Address object ({"address": "ip"})
	Servers                []interface{}          `json:"servers,omitempty"` // List of DNS servers. Can be string (IP or "localhost") or DnsServerObject.
	ClientIP               *string                `json:"clientIp,omitempty"` // Client IP for EDNS
	QueryStrategy          *string                `json:"queryStrategy,omitempty"` // "UseIP", "UseIPv4", "UseIPv6"
	DisableCache           *bool                  `json:"disableCache,omitempty"`
	DisableFallback        *bool                  `json:"disableFallback,omitempty"`
	DisableFallbackIfMatch *bool                  `json:"disableFallbackIfMatch,omitempty"` // New name for former `disableFeature`
	Tag                    *string                `json:"tag,omitempty"`                    // Tag for routing DNS queries through a specific outbound
}

// DnsServerObject is used when a server in DNSObject is an object.
type DnsServerObject struct {
	Address      *string  `json:"address,omitempty"`      // DNS server address
	Port         *int     `json:"port,omitempty"`         // DNS server port
	ClientIP     *string  `json:"clientIp,omitempty"`     // Client IP for EDNS, specific to this server
	SkipFallback *bool    `json:"skipFallback,omitempty"` // Deprecated, use disableFallbackIfMatch in DNSObject or routing rules
	Domains      []string `json:"domains,omitempty"`      // Domains for which this server should be used
	ExpectIps    []string `json:"expectIps,omitempty"`    // List of IPs to expect for specified domains (replaces expectIPs)
}

// RoutingObject defines routing rules.
// Docs: https://xtls.github.io/config/routing.html
type RoutingObject struct {
	DomainStrategy *string       `json:"domainStrategy,omitempty"` // "AsIs", "IPIfNonMatch", "IPOnDemand", "IPIfNonMatchElseAsIs"
	DomainMatcher  *string       `json:"domainMatcher,omitempty"`  // "linear", "mph" (new)
	Rules          []RoutingRule `json:"rules,omitempty"`
	Balancers      []Balancer    `json:"balancers,omitempty"`
}

// RoutingRule defines a single routing rule.
type RoutingRule struct {
	Type           *string  `json:"type,omitempty"` // Default "field"
	Domain         []string `json:"domain,omitempty"` // Domains to match
	IP             []string `json:"ip,omitempty"`     // Source/Destination IPs or CIDRs to match
	Port           *string  `json:"port,omitempty"`   // Destination port, e.g., "53", "1000-2000", 443
	Network        *string  `json:"network,omitempty"`  // "tcp", "udp", or "tcp,udp"
	SourceCidr     []string `json:"source,omitempty"`   // Source IPs or CIDRs
	UserEmail      []string `json:"user,omitempty"`     // User emails for authentication based routing
	InboundTag     []string `json:"inboundTag,omitempty"`
	Protocol       []string `json:"protocol,omitempty"`       // "http", "tls", "bittorrent", "dtls" (new)
	Attributes     *string  `json:"attrs,omitempty"`        // HTTP attributes matching
	OutboundTag    *string  `json:"outboundTag,omitempty"`  // Target outbound tag
	BalancerTag    *string  `json:"balancerTag,omitempty"`  // Target balancer tag
	DomainMatcher  *string  `json:"domainMatcher,omitempty"`// "linear", "mph" (per rule)
	Enabled        *bool    `json:"enabled,omitempty"`      // New: enable/disable rule
	SourcePort     *string  `json:"sourcePort,omitempty"`   // Source port, e.g., "53", "1000-2000", 12345 (New)
	TargetAddress  []string `json:"targetAddress,omitempty"`// Target address (domain or IP, New)
	TargetPort     *string  `json:"targetPort,omitempty"`   // Target port (New, complements Port which is often destination)
	TargetUser     []string `json:"targetUser,omitempty"`   // Target user (New)
}

// BalancerStrategyObject defines the strategy for a balancer.
type BalancerStrategyObject struct {
	Type     *string                `json:"type,omitempty"`     // "random", "leastPing"
	Settings map[string]interface{} `json:"settings,omitempty"` // Specific settings for the strategy, e.g., for leastPing: {"observerTag": "tag", "expected": ["ip"], "maxDeviation": 100, "tolerance": 1.5}
}

// Balancer defines a load balancer.
type Balancer struct {
	Tag      *string                 `json:"tag,omitempty"`
	Selector []string                `json:"selector,omitempty"` // Selects outbounds by tag pattern
	Strategy *BalancerStrategyObject `json:"strategy,omitempty"`
}

// PolicyObject defines local policy settings.
// Docs: https://xtls.github.io/config/policy.html
type PolicyObject struct {
	Levels map[string]LevelPolicy `json:"levels,omitempty"` // Key is user level as string (e.g., "0", "1")
	System *SystemPolicy          `json:"system,omitempty"`
}

// LevelPolicy defines policy for a specific user level.
type LevelPolicy struct {
	Handshake         *int  `json:"handshake,omitempty"`          // Connection handshake timeout in seconds, default 4.
	ConnIdle          *int  `json:"connIdle,omitempty"`           // Idle connection timeout in seconds, default 300.
	UplinkOnly        *int  `json:"uplinkOnly,omitempty"`         // Uplink only traffic duration in seconds, default 0 (disabled).
	DownlinkOnly      *int  `json:"downlinkOnly,omitempty"`       // Downlink only traffic duration in seconds, default 0 (disabled).
	StatsUserUplink   *bool `json:"statsUserUplink,omitempty"`    // Enable uplink stats for users of this level, default false.
	StatsUserDownlink *bool `json:"statsUserDownlink,omitempty"`  // Enable downlink stats for users of this level, default false.
	BufferSize        *int  `json:"bufferSize,omitempty"`         // Buffer size in KB. 0 for default Xray buffer. -1 for no buffer (direct copy).
	MaxConnections    *int  `json:"maxConnections,omitempty"`     // Maximum number of connections for a user of this level, default 0 (unlimited). (Newer field)
}

// SystemPolicy defines system-wide policies.
type SystemPolicy struct {
	StatsInboundUplink       *bool `json:"statsInboundUplink,omitempty"`       // Enable uplink stats for all inbounds, default false.
	StatsInboundDownlink     *bool `json:"statsInboundDownlink,omitempty"`     // Enable downlink stats for all inbounds, default false.
	StatsOutboundUplink      *bool `json:"statsOutboundUplink,omitempty"`      // Enable uplink stats for all outbounds, default false.
	StatsOutboundDownlink    *bool `json:"statsOutboundDownlink,omitempty"`    // Enable downlink stats for all outbounds, default false.
	OverrideAccessLogAddress *bool `json:"overrideAccessLogAddress,omitempty"` // Override inbound's access log address with Xray's address. (Newer field)
	OverrideAccessLogPort    *bool `json:"overrideAccessLogPort,omitempty"`    // Override inbound's access log port with Xray's port. (Newer field)
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
// Docs: https://xtls.github.io/config/outbounds.html
type OutboundObject struct {
	Tag            *string                `json:"tag,omitempty"`
	Protocol       *string                `json:"protocol,omitempty"` // e.g., "vmess", "vless", "freedom", "blackhole"
	Settings       map[string]interface{} `json:"settings,omitempty"` // Protocol-specific settings
	StreamSettings *StreamSettingsObject  `json:"streamSettings,omitempty"`
	ProxySettings  *ProxySettings         `json:"proxySettings,omitempty"` // For daisy-chaining
	SendThrough    *string                `json:"sendThrough,omitempty"`   // IP address to send traffic through
	Mux            *MuxObject             `json:"mux,omitempty"`
}

// StreamSettingsObject defines transport settings.
// Docs: https://xtls.github.io/config/transport.html#streamsettingsobject
type StreamSettingsObject struct {
	Network        *string                `json:"network,omitempty"`     // "tcp", "kcp", "ws", "http" (HTTP/2), "domainsocket", "quic", "grpc"
	Security       *string                `json:"security,omitempty"`    // "none", "tls", "xtls", "reality"
	TLSSettings    *TLSSettings           `json:"tlsSettings,omitempty"` // Used for security: "tls" or "reality"
	XTLSSettings   *XTLSSettings          `json:"xtlsSettings,omitempty"`// Used for security: "xtls"
	TCPSettings    *TCPSettings           `json:"tcpSettings,omitempty"`
	KCPSettings    *KCPSettings           `json:"kcpSettings,omitempty"`
	WSSettings     *WSSettings            `json:"wsSettings,omitempty"`
	HTTPSettings   *HTTP2Settings         `json:"httpSettings,omitempty"` // For HTTP/2
	DSSettings     *DomainSocketSettings  `json:"dsSettings,omitempty"`   // Domain Socket settings
	QUICSettings   *QUICSettings          `json:"quicSettings,omitempty"`
	GRPCSettings   *GRPCSettings          `json:"grpcSettings,omitempty"`
	SocketSettings *SocketOptions         `json:"sockopt,omitempty"`
}

// RealitySettingsObject defines settings for REALITY.
// Docs: https://xtls.github.io/config/transport.html#realitysettingsObject
type RealitySettingsObject struct {
	Show         *bool    `json:"show,omitempty"`     // If true, REALITY debug info will be shown in log
	Dest         *string  `json:"dest,omitempty"`     // Destination server address and port, e.g., "example.com:443"
	Xver         *int     `json:"xver,omitempty"`     // Proxy version, 0 or 1. Default 0.
	ServerNames  []string `json:"serverNames,omitempty"`// List of server names to mimic
	PrivateKey   *string  `json:"privateKey,omitempty"` // Base64 encoded private key (generated by `xray x25519`)
	MinClientVer *string  `json:"minClientVer,omitempty"`// Minimum client Xray version, e.g., "1.8.0"
	MaxClientVer *string  `json:"maxClientVer,omitempty"`// Maximum client Xray version
	MaxTimeDiff  *int64   `json:"maxTimeDiff,omitempty"`// Max time difference in ms, default 0 (disabled)
	ShortIds     []string `json:"shortIds,omitempty"`   // List of short IDs (0-15 byte hex strings)
	SpiderX      *string  `json:"spiderX,omitempty"`   // Path for crawling destination server, default "/"
}

// TLSSettings defines TLS settings.
// Docs: https://xtls.github.io/config/transport.html#tlssettingsobject
type TLSSettings struct {
	ServerName                     *string                `json:"serverName,omitempty"` // SNI
	AllowInsecure                  *bool                  `json:"allowInsecure,omitempty"`
	ALPN                           []string               `json:"alpn,omitempty"` // e.g., "http/1.1", "h2"
	MinVersion                     *string                `json:"minVersion,omitempty"` // "1.0", "1.1", "1.2", "1.3"
	MaxVersion                     *string                `json:"maxVersion,omitempty"`
	CipherSuites                   *string                `json:"cipherSuites,omitempty"` // Colon-separated list
	Certificates                   []Certificate          `json:"certificates,omitempty"`
	DisableSystemRoot              *bool                  `json:"disableSystemRoot,omitempty"`         // Default false
	EnableSessionResumption        *bool                  `json:"enableSessionResumption,omitempty"` // Default false
	Fingerprint                    *string                `json:"fingerprint,omitempty"`             // "chrome", "firefox", "safari", "ios", "android", "edge", "360", "qq", "random", "randomized"
	PinnedPeerCertificateChainSha256 []string            `json:"pinnedPeerCertificateChainSha256,omitempty"` // New
	MasterKeyLog                   *string                `json:"masterKeyLog,omitempty"`                   // New, path to log TLS master key
	RejectUnknownSni               *bool                  `json:"rejectUnknownSni,omitempty"`               // New, default false
	RealitySettings                *RealitySettingsObject `json:"realitySettings,omitempty"`              // New, for security: "reality"
}

// XTLSSettings defines XTLS settings (for security: "xtls").
// Docs: https://xtls.github.io/config/transport.html#xtlssettingsobject
type XTLSSettings struct {
	ServerName                     *string       `json:"serverName,omitempty"` // SNI
	AllowInsecure                  *bool         `json:"allowInsecure,omitempty"`
	ALPN                           []string      `json:"alpn,omitempty"`
	MinVersion                     *string       `json:"minVersion,omitempty"` // "1.0", "1.1", "1.2", "1.3"
	MaxVersion                     *string       `json:"maxVersion,omitempty"`
	CipherSuites                   *string       `json:"cipherSuites,omitempty"`
	Certificates                   []Certificate `json:"certificates,omitempty"`
	DisableSystemRoot              *bool         `json:"disableSystemRoot,omitempty"` // Default false
	Fingerprint                    *string       `json:"fingerprint,omitempty"`       // Not typically used with XTLS server-side
	PinnedPeerCertificateChainSha256 []string   `json:"pinnedPeerCertificateChainSha256,omitempty"`
	// EnableSessionResumption is NOT part of XTLSSettings as per docs
}


// Certificate defines a TLS certificate.
// Docs: https://xtls.github.io/config/transport.html#certificateobject
type Certificate struct {
	Usage           *string  `json:"usage,omitempty"` // "encipherment", "verify", "issue", "ignored" (new)
	CertificateFile *string  `json:"certificateFile,omitempty"`
	KeyFile         *string  `json:"keyFile,omitempty"`
	Certificate     []string `json:"certificate,omitempty"`     // Certificate content as string array (PEM format)
	Key             []string `json:"key,omitempty"`             // Key content as string array (PEM format)
	OcspStapling    *uint32  `json:"ocspStapling,omitempty"`    // Refresh interval in seconds, default 3600
	OneTimeLoading  *bool    `json:"oneTimeLoading,omitempty"`  // New, default false
}

// TCPSettings specific configurations for TCP transport.
// Docs: https://xtls.github.io/config/transport.html#tcpsettingsobject
type TCPSettings struct {
	Header                *HeaderObject `json:"header,omitempty"`
	AcceptProxyProtocol   *bool         `json:"acceptProxyProtocol,omitempty"` // For HAProxy's PROXY protocol
	TCPNoDelay            *bool         `json:"tcpNoDelay,omitempty"`          // Corresponds to sockopt TCP_NODELAY. Often set here for convenience.
	TCPKeepAliveInterval  *int          `json:"tcpKeepAliveInterval,omitempty"`// Corresponds to sockopt TCP_KEEPINTVL. Often set here.
	Congestion            *string       `json:"congestion,omitempty"`          // Corresponds to sockopt TCP_CONGESTION. Often set here.
}

// HeaderObject for TCP header obfuscation.
// Docs: https://xtls.github.io/config/transport.html#headerobject
type HeaderObject struct {
	Type     *string         `json:"type,omitempty"` // "none", "http", "srtp", "utp", "wechat-video", "dtls", "wireguard"
	Request  *RequestConfig  `json:"request,omitempty"`  // Used if type is "http"
	Response *ResponseConfig `json:"response,omitempty"` // Used if type is "http"
}

// RequestConfig for HTTP header
type RequestConfig struct {
	Version *string                `json:"version,omitempty"`
	Method  *string                `json:"method,omitempty"`
	Path    []string               `json:"path,omitempty"`
	Headers map[string]interface{} `json:"headers,omitempty"` // Value can be string or []string
}

// ResponseConfig for HTTP header
type ResponseConfig struct {
	Version *string                `json:"version,omitempty"`
	Status  *string                `json:"status,omitempty"`
	Reason  *string                `json:"reason,omitempty"`
	Headers map[string]interface{} `json:"headers,omitempty"` // Value can be string or []string
}

// KCPSettings specific configurations for mKCP transport.
// Docs: https://xtls.github.io/config/transport.html#kcpsettingsobject
type KCPSettings struct {
	Mtu              *int          `json:"mtu,omitempty"`
	Tti              *int          `json:"tti,omitempty"`
	UplinkCapacity   *int          `json:"uplinkCapacity,omitempty"`
	DownlinkCapacity *int          `json:"downlinkCapacity,omitempty"`
	Congestion       *bool         `json:"congestion,omitempty"`
	ReadBufferSize   *int          `json:"readBufferSize,omitempty"`
	WriteBufferSize  *int          `json:"writeBufferSize,omitempty"`
	Header           *HeaderObject `json:"header,omitempty"`
	Seed             *string       `json:"seed,omitempty"`
}

// WSSettings specific configurations for WebSocket transport.
// Docs: https://xtls.github.io/config/transport.html#websocketsettingsobject
type WSSettings struct {
	Path                 *string           `json:"path,omitempty"`
	Headers              map[string]string `json:"headers,omitempty"`
	AcceptProxyProtocol  *bool             `json:"acceptProxyProtocol,omitempty"`
	MaxEarlyData         *int32            `json:"maxEarlyData,omitempty"`         // New
	EarlyDataHeaderName  *string           `json:"earlyDataHeaderName,omitempty"`  // New
	UseBrowserForwarding *bool             `json:"useBrowserForwarding,omitempty"` // New, default false
}

// HTTP2Settings (H2) specific configurations.
// Docs: https://xtls.github.io/config/transport.html#http2settingsobject
type HTTP2Settings struct {
	Host               []string `json:"host,omitempty"` // List of domains
	Path               *string  `json:"path,omitempty"` // Path for HTTP/2 requests
	ReadIdleTimeout    *int     `json:"readIdleTimeout,omitempty"`    // New
	HealthCheckTimeout *int     `json:"healthCheckTimeout,omitempty"` // New
	Method             *string  `json:"method,omitempty"`             // New, default "PUT"
}

// QUICSettings specific configurations for QUIC transport.
// Docs: https://xtls.github.io/config/transport.html#quicsettingsobject
type QUICSettings struct {
	Security *string       `json:"security,omitempty"` // "none", "aes-128-gcm", "chacha20-poly1305"
	Key      *string       `json:"key,omitempty"`
	Header   *HeaderObject `json:"header,omitempty"` // Header obfuscation
}

// GRPCSettings specific configurations for gRPC transport.
// Docs: https://xtls.github.io/config/transport.html#grpcsettingsobject
type GRPCSettings struct {
	ServiceName         *string `json:"serviceName,omitempty"`
	MultiMode           *bool   `json:"multiMode,omitempty"`           // Default false
	IdleTimeout         *int    `json:"idle_timeout,omitempty"`        // Seconds
	HealthCheckTimeout  *int    `json:"health_check_timeout,omitempty"`// Seconds
	PermitWithoutStream *bool   `json:"permit_without_stream,omitempty"`// Default false
	InitialWindowsSize  *int32  `json:"initial_windows_size,omitempty"`// Bytes
	UserAgent           *string `json:"user_agent,omitempty"`          // New
}

// SocketOptions for various transport protocols.
// Docs: https://xtls.github.io/config/transport.html#socketoptions
type SocketOptions struct {
	Mark                 *int    `json:"mark,omitempty"`                   // SO_MARK
	TCPFastOpen          *int    `json:"tcpFastOpen,omitempty"`            // TCP Fast Open queue length
	Tproxy               *string `json:"tproxy,omitempty"`                 // "redirect", "tproxy", "off"
	DomainStrategy       *string `json:"domainStrategy,omitempty"`         // "UseIP", "UseIPv4", "UseIPv6", "AsIs"
	DialerProxy          *string `json:"dialerProxy,omitempty"`            // Tag of an outbound proxy
	AcceptProxyProtocol  *bool   `json:"acceptProxyProtocol,omitempty"`
	TCPKeepAliveInterval *int    `json:"tcpKeepAliveInterval,omitempty"`   // TCP_KEEPINTVL (seconds)
	TCPKeepAliveIdle     *int    `json:"tcpKeepAliveIdle,omitempty"`       // TCP_KEEPIDLE (seconds)
	TCPUserTimeout       *int    `json:"tcpUserTimeout,omitempty"`         // New, TCP_USER_TIMEOUT (milliseconds)
	TCPCongestion        *string `json:"tcpCongestion,omitempty"`          // TCP congestion control algorithm
	Interface            *string `json:"interface,omitempty"`              // Bind to a specific network interface
	V6Only               *bool   `json:"v6only,omitempty"`                 // IPV6_V6ONLY
	TCPMptcp             *bool   `json:"tcpMptcp,omitempty"`               // TCP MPTCP
	TCPNoDelay           *bool   `json:"tcpNoDelay,omitempty"`             // TCP_NODELAY
	UDPReusable          *bool   `json:"udpReusable,omitempty"`            // New, SO_REUSEADDR for UDP
	UDPTimeout           *int    `json:"udpTimeout,omitempty"`             // New, timeout for UDP connections in seconds
}


// SniffingObject defines settings for content sniffing.
// Docs: https://xtls.github.io/config/inbounds.html#sniffingobject
type SniffingObject struct {
	Enabled         *bool    `json:"enabled,omitempty"`                  // Default true if object exists
	DestOverride    []string `json:"destOverride,omitempty"`           // "http", "tls", "fakedns", "quic" (new)
	DomainsExcluded []string `json:"domainsExcluded,omitempty"`        // New
	MetadataOnly    *bool    `json:"metadataOnly,omitempty"`           // New, for fakedns to sniff SNI/ALPN only
	RouteOnly       *bool    `json:"routeOnly,omitempty"`              // New, if true, sniffing result is only for routing
	AppProtocol     []string `json:"appProtocol,omitempty"`            // New in v1.8.7, for application layer protocol sniffing
	AppProtocolPort *string  `json:"appProtocolPort,omitempty"`        // New in v1.8.7, port range for app protocol sniffing
}

// AllocateObject defines port allocation strategy.
// Docs: https://xtls.github.io/config/inbounds.html#allocateobject
type AllocateObject struct {
	Strategy    *string `json:"strategy,omitempty"`    // "always" or "random", default "always"
	Refresh     *int    `json:"refresh,omitempty"`     // In minutes, default 5
	Concurrency *int    `json:"concurrency,omitempty"` // Default 3
}

// ProxySettings for outbound proxying.
// Docs: https://xtls.github.io/config/outbounds.html#proxysettings
type ProxySettings struct {
	Tag            *string `json:"tag,omitempty"` // Tag of another outbound to use as transport proxy
	TransportLayer *bool   `json:"transportLayer,omitempty"` // VLESS specific, handle proxying at transport layer
}

// MuxObject for connection multiplexing.
// Docs: https://xtls.github.io/config/outbounds.html#muxobject
type MuxObject struct {
	Enabled         *bool   `json:"enabled,omitempty"`         // Default false
	Concurrency     *int    `json:"concurrency,omitempty"`     // 1-1024, default 8
	XUDPConcurrency *int    `json:"xudpConcurrency,omitempty"` // For QUIC, default 8
	XUDPProxyUDP443 *string `json:"xudpProxyUDP443,omitempty"`// "disabled", "skip", "prefer" (for QUIC)
	Padding         *bool   `json:"padding,omitempty"`         // New in v1.8.4, enable padding for mux frames
}


// TransportObject defines global transport settings.
// Docs: https://xtls.github.io/config/transport.html#transportobject
type TransportObject struct {
	TCPSettings    *TCPSettings          `json:"tcpSettings,omitempty"`
	KCPSettings    *KCPSettings          `json:"kcpSettings,omitempty"`
	WSSettings     *WSSettings           `json:"wsSettings,omitempty"`
	HTTPSettings   *HTTP2Settings        `json:"httpSettings,omitempty"`
	DSSettings     *DomainSocketSettings `json:"dsSettings,omitempty"`
	QUICSettings   *QUICSettings         `json:"quicSettings,omitempty"`
	GRPCSettings   *GRPCSettings         `json:"grpcSettings,omitempty"`
	SocketSettings *SocketOptions        `json:"sockopt,omitempty"`
}

// DomainSocketSettings for Domain Socket transport.
// Docs: https://xtls.github.io/config/transport.html#domainsocketsettingsobject
type DomainSocketSettings struct {
	Path     *string `json:"path,omitempty"`      // Required
	Abstract *bool   `json:"abstract,omitempty"`  // Default false, Linux only
	Padding  *bool   `json:"padding,omitempty"`   // Default false, Linux abstract only
}


// StatsObject enables or disables statistics. (Empty object means enabled)
type StatsObject struct{}

// ReverseObject defines reverse proxy settings.
// Docs: https://xtls.github.io/config/reverse.html
type ReverseObject struct {
	Bridges []Bridge `json:"bridges,omitempty"`
	Portals []Portal `json:"portals,omitempty"`
}

// Bridge for reverse proxy.
type Bridge struct {
	Tag    *string `json:"tag,omitempty"`    // Required
	Domain *string `json:"domain,omitempty"` // Required
}

// Portal for reverse proxy.
type Portal struct {
	Tag    *string `json:"tag,omitempty"`    // Required
	Domain *string `json:"domain,omitempty"` // Required, or use `domains` for multiple subdomains
	// Domains []string `json:"domains,omitempty"` // Alternative to single domain
}

// FakeDNSObject enables FakeDNS.
// Docs: https://xtls.github.io/config/fakedns.html
type FakeDNSObject struct {
	IPPool   *string `json:"ipPool,omitempty"`  // CIDR format, e.g., "198.18.0.0/15", Required
	PoolSize *int64  `json:"poolSize,omitempty"` // Default 65535
	UDPExtra *string `json:"udpExtra,omitempty"` // Deprecated: "disable"
}

// MetricsObject defines metrics settings.
// Docs: https://xtls.github.io/config/metrics.html
type MetricsObject struct {
	Tag    *string `json:"tag,omitempty"`    // Required
	Listen *string `json:"listen,omitempty"` // Optional, address to listen on (Prometheus format)
	Port   *uint16 `json:"port,omitempty"`   // Optional, port to listen on (Prometheus format)
}

// ObservatoryObject defines connection observation settings.
// Docs: https://xtls.github.io/config/observatory.html
type ObservatoryObject struct {
	SubjectSelector []string `json:"subjectSelector,omitempty"` // Required
	ProbeURL        *string  `json:"probeURL,omitempty"`        // Required
	ProbeInterval   *string  `json:"probeInterval,omitempty"`   // Required, duration string e.g. "10m", "1h"
}

// BurstObservatoryObject defines settings for burst connection observation. (Project X specific)
type BurstObservatoryObject struct {
	SubjectSelector []string `json:"subjectSelector,omitempty"`
	ProbeURL        *string  `json:"probeURL,omitempty"`
	ProbeInterval   *string  `json:"probeInterval,omitempty"`
}

// Helper function to get a pointer to an int. Useful for optional int fields.
func IntPtr(i int) *int {
	return &i
}
