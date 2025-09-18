// Package authz provides authorization functionality using Casbin for access control
// with RBAC (Role-Based Access Control) model and GORM database integration.
package authz

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/wire"
	"gorm.io/gorm"
)

const (
	// defaultAclModel defines the default Casbin access control model.
	defaultAclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (p.act == "*" || r.act == p.act)`
)

// Authz defines an authorizer that provides authorization functionality.
type Authz struct {
	*casbin.SyncedEnforcer // Uses Casbin's synchronized enforcer
}

// Option defines a function option type for customizing NewAuthz behavior.
type Option func(*authzConfig)

// authzConfig is the configuration structure for the authorizer.
type authzConfig struct {
	aclModel           string        // Casbin model string
	autoLoadPolicyTime time.Duration // Time interval for automatic policy loading
}

// ProviderSet is a Wire Provider set used to declare dependency injection rules.
// Contains the NewAuthz constructor for generating Authz instances.
var ProviderSet = wire.NewSet(NewAuthz, DefaultOptions)

// defaultAuthzConfig returns a default configuration.
func defaultAuthzConfig() *authzConfig {
	return &authzConfig{
		// Default to using the built-in ACL model
		aclModel: defaultAclModel,
		// Default automatic policy loading time interval
		autoLoadPolicyTime: 5 * time.Second,
	}
}

// DefaultOptions provides default authorizer option configuration.
func DefaultOptions() []Option {
	return []Option{
		// Use the default ACL model
		WithAclModel(defaultAclModel),
		// Set automatic policy loading time interval to 10 seconds
		WithAutoLoadPolicyTime(10 * time.Second),
	}
}

// WithAclModel allows customizing the ACL model through options.
func WithAclModel(model string) Option {
	return func(cfg *authzConfig) {
		cfg.aclModel = model
	}
}

// WithAutoLoadPolicyTime allows customizing the automatic policy loading time interval through options.
func WithAutoLoadPolicyTime(interval time.Duration) Option {
	return func(cfg *authzConfig) {
		cfg.autoLoadPolicyTime = interval
	}
}

// NewAuthz creates an authorizer using Casbin for authorization, supporting custom configuration through functional options pattern.
func NewAuthz(db *gorm.DB, opts ...Option) (*Authz, error) {
	// Initialize default configuration
	cfg := defaultAuthzConfig()

	// Apply all options
	for _, opt := range opts {
		opt(cfg)
	}

	// Initialize Gorm adapter for Casbin authorizer
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err // Return error
	}

	// Load Casbin model from configuration
	m, _ := model.NewModelFromString(cfg.aclModel)

	// Initialize the enforcer
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err // Return error
	}

	// Load policies from database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err // Return error
	}

	// Start automatic policy loading with configured time interval
	enforcer.StartAutoLoadPolicy(cfg.autoLoadPolicyTime)

	// Return new authorizer instance
	return &Authz{enforcer}, nil
}

// Authorize performs authorization checks.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	// Call Enforce method for authorization check
	return a.Enforce(sub, obj, act)
}
