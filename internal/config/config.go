package config

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentDevelopment Environment = "development"
)

// GetBaseURL returns the base URL for the environment
func (e Environment) GetBaseURL() string {
	switch e {
	case EnvironmentDevelopment:
		return "https://api.dev.eka.care"
	case EnvironmentProduction:
		return "https://api.eka.care"
	default:
		return "https://api.eka.care" // Default to production
	}
}

// Config holds the internal configuration for the ABDM client
type Config struct {
	Environment        Environment
	BaseURL            string
	ClientID           string // Client ID for authentication
	ClientSecret       string // Client Secret for authentication
	AuthorizationToken string // JWT token for API calls (set after login)
	Timeout            time.Duration
	MaxRetries         int
	UserAgent          string
	LogLevel           string
	HTTPClient         *http.Client
	DisableSSL         bool
	Region             string
	RetryMode          string
	MaxBackoffDelay    time.Duration
	RequestTimeout     time.Duration
	ResponseTimeout    time.Duration
	ConnectionTimeout  time.Duration

	mu      sync.RWMutex
	tokenFn func(context.Context) (string, error)
}

// Ensure Config implements interfaces.Config
var _ interfaces.Config = (*Config)(nil)

// NewConfig creates a new configuration with defaults
func NewConfig() *Config {
	return &Config{
		Environment:       EnvironmentProduction,
		BaseURL:           EnvironmentProduction.GetBaseURL(),
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		UserAgent:         "eka-sdk-go/1.0",
		LogLevel:          "info",
		RetryMode:         "standard",
		MaxBackoffDelay:   20 * time.Second,
		RequestTimeout:    30 * time.Second,
		ResponseTimeout:   30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}

// Interface implementation methods
func (c *Config) GetEnvironment() Environment { return c.Environment }
func (c *Config) GetBaseURL() string          { return c.BaseURL }
func (c *Config) GetAPIKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.AuthorizationToken
}
func (c *Config) GetTimeout() time.Duration           { return c.Timeout }
func (c *Config) GetMaxRetries() int                  { return c.MaxRetries }
func (c *Config) GetUserAgent() string                { return c.UserAgent }
func (c *Config) GetLogLevel() string                 { return c.LogLevel }
func (c *Config) GetHTTPClient() *http.Client         { return c.HTTPClient }
func (c *Config) GetDisableSSL() bool                 { return c.DisableSSL }
func (c *Config) GetRegion() string                   { return c.Region }
func (c *Config) GetRetryMode() string                { return c.RetryMode }
func (c *Config) GetMaxBackoffDelay() time.Duration   { return c.MaxBackoffDelay }
func (c *Config) GetRequestTimeout() time.Duration    { return c.RequestTimeout }
func (c *Config) GetResponseTimeout() time.Duration   { return c.ResponseTimeout }
func (c *Config) GetConnectionTimeout() time.Duration { return c.ConnectionTimeout }

// GetClientID returns the client ID for authentication
func (c *Config) GetClientID() string { return c.ClientID }

// GetClientSecret returns the client secret for authentication
func (c *Config) GetClientSecret() string { return c.ClientSecret }

// SetAuthorizationToken sets a static JWT for API calls.
func (c *Config) SetAuthorizationToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthorizationToken = token
}

// SetTokenFunc installs a resolver consulted before every request. Use this instead
// of SetAuthorizationToken for long-running processes so expired tokens refresh.
func (c *Config) SetTokenFunc(fn func(context.Context) (string, error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tokenFn = fn
}

// Token implements interfaces.TokenProvider.
func (c *Config) Token(ctx context.Context) (string, error) {
	c.mu.RLock()
	fn, static := c.tokenFn, c.AuthorizationToken
	c.mu.RUnlock()
	if fn == nil {
		return static, nil
	}
	return fn(ctx)
}
