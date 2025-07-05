package models

import "time"

// HAProxyConfig is the top-level structure for an HAProxy configuration.
// It includes metadata for storage and management within ProxyPanel.
type HAProxyConfig struct {
	ID          string    `json:"id" gorm:"primaryKey" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"` // Internal ID for database
	Name        string    `json:"name" gorm:"uniqueIndex" example:"My Default HAProxy Config"`            // User-defined name for the config
	Description string    `json:"description,omitempty" example:"Main HAProxy load balancer configuration"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T13:00:00Z"`

	// Core HAProxy configuration sections
	Global    *HAProxyGlobal    `json:"global,omitempty"`
	Defaults  *HAProxyDefaults  `json:"defaults,omitempty"`
	Frontends []HAProxyFrontend `json:"frontends,omitempty"`
	Backends  []HAProxyBackend  `json:"backends,omitempty"`
	Listens   []HAProxyListen   `json:"listens,omitempty"`
	Stats     *HAProxyStats     `json:"stats,omitempty"`
}

// HAProxyGlobal defines global HAProxy settings.
type HAProxyGlobal struct {
	Daemon          *bool             `json:"daemon,omitempty"`
	User            *string           `json:"user,omitempty"`
	Group           *string           `json:"group,omitempty"`
	Chroot          *string           `json:"chroot,omitempty"`
	Pidfile         *string           `json:"pidfile,omitempty"`
	Maxconn         *int              `json:"maxconn,omitempty"`
	Nbthread        *int              `json:"nbthread,omitempty"`
	SSLDefaultBind  *string           `json:"ssl_default_bind,omitempty"`
	SSLDefaultServer *string          `json:"ssl_default_server,omitempty"`
	Tune            map[string]string `json:"tune,omitempty"`
	CustomDirectives []string         `json:"custom_directives,omitempty"`
}

// HAProxyDefaults defines default settings for HAProxy.
type HAProxyDefaults struct {
	Mode               *string           `json:"mode,omitempty"` // "http", "tcp", "health"
	Timeout            *HAProxyTimeouts  `json:"timeout,omitempty"`
	Retries            *int              `json:"retries,omitempty"`
	Option             []string          `json:"option,omitempty"`
	Errorfile          map[string]string `json:"errorfile,omitempty"`
	CustomDirectives   []string          `json:"custom_directives,omitempty"`
}

// HAProxyTimeouts defines timeout settings.
type HAProxyTimeouts struct {
	Connect     *string `json:"connect,omitempty"`     // e.g., "5000ms", "5s"
	Client      *string `json:"client,omitempty"`      // e.g., "50000ms", "50s"
	Server      *string `json:"server,omitempty"`      // e.g., "50000ms", "50s"
	Check       *string `json:"check,omitempty"`       // e.g., "3500ms", "3.5s"
	Queue       *string `json:"queue,omitempty"`       // e.g., "5000ms", "5s"
	Tunnel      *string `json:"tunnel,omitempty"`      // e.g., "3600s", "1h"
	HTTPRequest *string `json:"http_request,omitempty"` // e.g., "10s"
	HTTPKeepAlive *string `json:"http_keep_alive,omitempty"` // e.g., "2s"
}

// HAProxyFrontend defines a frontend configuration.
type HAProxyFrontend struct {
	Name             string            `json:"name" example:"web_frontend"`
	Bind             []HAProxyBind     `json:"bind,omitempty"`
	Mode             *string           `json:"mode,omitempty"` // "http", "tcp"
	Maxconn          *int              `json:"maxconn,omitempty"`
	ACLs             []HAProxyACL      `json:"acls,omitempty"`
	UseBackend       []HAProxyUseBackend `json:"use_backend,omitempty"`
	DefaultBackend   *string           `json:"default_backend,omitempty"`
	Option           []string          `json:"option,omitempty"`
	CustomDirectives []string          `json:"custom_directives,omitempty"`
}

// HAProxyBackend defines a backend configuration.
type HAProxyBackend struct {
	Name             string            `json:"name" example:"web_servers"`
	Mode             *string           `json:"mode,omitempty"` // "http", "tcp"
	Balance          *string           `json:"balance,omitempty"` // "roundrobin", "leastconn", "source", etc.
	Servers          []HAProxyServer   `json:"servers,omitempty"`
	Option           []string          `json:"option,omitempty"`
	HTTPCheck        *string           `json:"http_check,omitempty"`
	CustomDirectives []string          `json:"custom_directives,omitempty"`
}

// HAProxyListen defines a listen configuration (combined frontend/backend).
type HAProxyListen struct {
	Name             string            `json:"name" example:"web_service"`
	Bind             []HAProxyBind     `json:"bind,omitempty"`
	Mode             *string           `json:"mode,omitempty"` // "http", "tcp"
	Balance          *string           `json:"balance,omitempty"`
	Servers          []HAProxyServer   `json:"servers,omitempty"`
	Option           []string          `json:"option,omitempty"`
	CustomDirectives []string          `json:"custom_directives,omitempty"`
}

// HAProxyBind defines a bind directive.
type HAProxyBind struct {
	Address    string            `json:"address" example:"*:80"`
	SSL        *bool             `json:"ssl,omitempty"`
	Cert       *string           `json:"cert,omitempty"`
	Ciphers    *string           `json:"ciphers,omitempty"`
	Options    []string          `json:"options,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// HAProxyServer defines a server in a backend.
type HAProxyServer struct {
	Name       string            `json:"name" example:"web1"`
	Address    string            `json:"address" example:"192.168.1.10:80"`
	Weight     *int              `json:"weight,omitempty"`
	Maxconn    *int              `json:"maxconn,omitempty"`
	Check      *bool             `json:"check,omitempty"`
	CheckPort  *int              `json:"check_port,omitempty"`
	CheckInter *string           `json:"check_inter,omitempty"` // e.g., "2000ms", "2s"
	Rise       *int              `json:"rise,omitempty"`
	Fall       *int              `json:"fall,omitempty"`
	Backup     *bool             `json:"backup,omitempty"`
	Disabled   *bool             `json:"disabled,omitempty"`
	SSL        *bool             `json:"ssl,omitempty"`
	Verify     *string           `json:"verify,omitempty"` // "none", "required"
	Options    []string          `json:"options,omitempty"`
}

// HAProxyACL defines an Access Control List rule.
type HAProxyACL struct {
	Name      string `json:"name" example:"is_api"`
	Condition string `json:"condition" example:"path_beg /api"`
}

// HAProxyUseBackend defines a use_backend directive.
type HAProxyUseBackend struct {
	Backend   string `json:"backend" example:"api_servers"`
	Condition string `json:"condition" example:"is_api"`
}

// HAProxyStats defines statistics configuration.
type HAProxyStats struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	URI      *string `json:"uri,omitempty" example:"/stats"`
	Realm    *string `json:"realm,omitempty" example:"HAProxy Statistics"`
	Auth     *string `json:"auth,omitempty" example:"admin:password"`
	Refresh  *string `json:"refresh,omitempty" example:"30s"`
	HideVersion *bool `json:"hide_version,omitempty"`
}

// HAProxyConfigUpdate represents partial updates to an HAProxy configuration.
type HAProxyConfigUpdate struct {
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Global      *HAProxyGlobal     `json:"global,omitempty"`
	Defaults    *HAProxyDefaults   `json:"defaults,omitempty"`
	Frontends   *[]HAProxyFrontend `json:"frontends,omitempty"`
	Backends    *[]HAProxyBackend  `json:"backends,omitempty"`
	Listens     *[]HAProxyListen   `json:"listens,omitempty"`
	Stats       *HAProxyStats      `json:"stats,omitempty"`
}