// Package ekasdk provides a comprehensive Go SDK for Eka Care APIs.
//
// # Quick Start
//
// The recommended way to create a client is using environment variables:
//
//	// Set environment variables
//	// export EKA_ENVIRONMENT=production
//	// export EKA_CLIENT_ID=your-client-id
//	// export EKA_CLIENT_SECRET=your-client-secret
//
//	client := ekasdk.NewFromEnv()
//	if err := client.Login(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// # Alternative Configuration
//
// For explicit configuration (not recommended for production):
//
//	client := ekasdk.New(
//		ekasdk.WithEnvironment(ekasdk.EnvironmentProduction),
//		ekasdk.WithClientID("your-client-id"),
//		ekasdk.WithClientSecret("your-client-secret"),
//	)
//	if err := client.Login(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// # Security Best Practices
//
// - Use environment variables for all credentials
// - Never commit credentials to version control
// - Use different credentials for different environments
// - Store secrets in secure credential management systems
package ekasdk

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eka-care/eka-sdk-go/auth"
	"github.com/eka-care/eka-sdk-go/internal/config"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
	"github.com/eka-care/eka-sdk-go/services/abdm"
)

// Environment represents the deployment environment
type Environment string

const (
	// EnvironmentProduction represents the production environment
	EnvironmentProduction Environment = "production"
	// EnvironmentDevelopment represents the development environment
	EnvironmentDevelopment Environment = "development"
)

// Client represents the main Eka SDK client
type Client struct {
	config              interfaces.Config
	credentialsProvider auth.CredentialsProvider

	// Service clients
	Auth *auth.Service
	ABDM *abdm.Client
}

// Option represents a configuration option for the client
type Option func(*ClientOptions)

// ClientOptions holds the configuration options for the client
type ClientOptions struct {
	Environment         Environment
	ClientID            string // Client ID for authentication
	ClientSecret        string // Client Secret for authentication
	CredentialsProvider auth.CredentialsProvider
	Timeout             time.Duration
	MaxRetries          int
	UserAgent           string
	LogLevel            string
	HTTPClient          *http.Client
	DisableSSL          bool
	Region              string
	RetryMode           string
	MaxBackoffDelay     time.Duration
	RequestTimeout      time.Duration
	ResponseTimeout     time.Duration
	ConnectionTimeout   time.Duration
}

// DefaultClientOptions returns the default client options
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Environment:       EnvironmentProduction,
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		UserAgent:         "eka-sdk-go/1.0.0",
		LogLevel:          "info",
		HTTPClient:        &http.Client{},
		DisableSSL:        false,
		Region:            "us",
		RetryMode:         "standard",
		MaxBackoffDelay:   20 * time.Second,
		RequestTimeout:    30 * time.Second,
		ResponseTimeout:   30 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}
}

// WithEnvironment sets the environment
func WithEnvironment(env Environment) Option {
	return func(opts *ClientOptions) {
		opts.Environment = env
	}
}

// WithClientID sets the client ID for authentication
func WithClientID(clientID string) Option {
	return func(opts *ClientOptions) {
		opts.ClientID = clientID
	}
}

// WithClientSecret sets the client secret for authentication
func WithClientSecret(clientSecret string) Option {
	return func(opts *ClientOptions) {
		opts.ClientSecret = clientSecret
	}
}

// WithCredentialsProvider sets the credentials provider
func WithCredentialsProvider(provider auth.CredentialsProvider) Option {
	return func(opts *ClientOptions) {
		opts.CredentialsProvider = provider
	}
}

// WithTimeout sets the timeout
func WithTimeout(timeout time.Duration) Option {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(maxRetries int) Option {
	return func(opts *ClientOptions) {
		opts.MaxRetries = maxRetries
	}
}

// WithUserAgent sets the user agent
func WithUserAgent(userAgent string) Option {
	return func(opts *ClientOptions) {
		opts.UserAgent = userAgent
	}
}

// WithLogLevel sets the log level
func WithLogLevel(logLevel string) Option {
	return func(opts *ClientOptions) {
		opts.LogLevel = logLevel
	}
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(opts *ClientOptions) {
		opts.HTTPClient = httpClient
	}
}

// WithDisableSSL sets whether to disable SSL verification
func WithDisableSSL(disableSSL bool) Option {
	return func(opts *ClientOptions) {
		opts.DisableSSL = disableSSL
	}
}

// New creates a new Eka SDK client with the given options
func New(opts ...Option) *Client {
	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Create internal config manually
	internalConfig := &config.Config{
		Environment:       config.Environment(options.Environment),
		BaseURL:           getBaseURL(options.Environment),
		ClientID:          options.ClientID,
		ClientSecret:      options.ClientSecret,
		Timeout:           options.Timeout,
		MaxRetries:        options.MaxRetries,
		UserAgent:         options.UserAgent,
		LogLevel:          options.LogLevel,
		HTTPClient:        options.HTTPClient,
		DisableSSL:        options.DisableSSL,
		Region:            options.Region,
		RetryMode:         options.RetryMode,
		MaxBackoffDelay:   options.MaxBackoffDelay,
		RequestTimeout:    options.RequestTimeout,
		ResponseTimeout:   options.ResponseTimeout,
		ConnectionTimeout: options.ConnectionTimeout,
	}

	return &Client{
		config:              internalConfig,
		credentialsProvider: options.CredentialsProvider,
		Auth:                auth.NewService(internalConfig),
		ABDM:                createABDMClient(internalConfig),
	}
}

// NewFromEnv creates a new client using environment variables
func NewFromEnv() *Client {
	options := DefaultClientOptions()

	// Read environment variables
	if env := os.Getenv("EKA_ENVIRONMENT"); env != "" {
		options.Environment = Environment(env)
	}

	if clientID := os.Getenv("EKA_CLIENT_ID"); clientID != "" {
		options.ClientID = clientID
	}

	if clientSecret := os.Getenv("EKA_CLIENT_SECRET"); clientSecret != "" {
		options.ClientSecret = clientSecret
	}

	if timeout := os.Getenv("EKA_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			options.Timeout = time.Duration(t) * time.Second
		}
	}

	if maxRetries := os.Getenv("EKA_MAX_RETRIES"); maxRetries != "" {
		if mr, err := strconv.Atoi(maxRetries); err == nil {
			options.MaxRetries = mr
		}
	}

	if userAgent := os.Getenv("EKA_USER_AGENT"); userAgent != "" {
		options.UserAgent = userAgent
	}

	if logLevel := os.Getenv("EKA_LOG_LEVEL"); logLevel != "" {
		options.LogLevel = logLevel
	}

	if disableSSL := os.Getenv("EKA_DISABLE_SSL"); disableSSL != "" {
		if ds, err := strconv.ParseBool(disableSSL); err == nil {
			options.DisableSSL = ds
		}
	}

	if region := os.Getenv("EKA_REGION"); region != "" {
		options.Region = region
	}

	return New(
		WithEnvironment(options.Environment),
		WithClientID(options.ClientID),
		WithClientSecret(options.ClientSecret),
		WithTimeout(options.Timeout),
		WithMaxRetries(options.MaxRetries),
		WithUserAgent(options.UserAgent),
		WithLogLevel(options.LogLevel),
		WithDisableSSL(options.DisableSSL),
	)
}

// getBaseURL returns the base URL for the given environment
func getBaseURL(env Environment) string {
	switch env {
	case EnvironmentProduction:
		return "https://api.eka.care"
	case EnvironmentDevelopment:
		return "https://api-dev.eka.care"
	default:
		return "https://api.eka.care"
	}
}

// createABDMClient creates an ABDM client from the internal config
func createABDMClient(cfg *config.Config) *abdm.Client {
	// The ABDM client now just organizes services and uses the main config
	return abdm.NewClient(cfg)
}

// GetCredentials retrieves the current credentials using the configured provider
func (c *Client) GetCredentials(ctx context.Context) (*auth.Credentials, error) {
	if c.credentialsProvider == nil {
		return nil, fmt.Errorf("no credentials provider configured")
	}
	return c.credentialsProvider.Retrieve(ctx)
}

// SetCredentialsProvider sets a new credentials provider
func (c *Client) SetCredentialsProvider(provider auth.CredentialsProvider) {
	c.credentialsProvider = provider
}

// NewClientCredentialsProvider creates a client credentials provider using this client's auth service
func (c *Client) NewClientCredentialsProvider(req *auth.ClientLoginRequest) *auth.ClientCredentialsProvider {
	return auth.NewClientCredentialsProvider(c.Auth, req)
}

// Login performs authentication using client credentials and sets up the client for API calls
func (c *Client) Login(ctx context.Context) error {
	cfg := c.config.(*config.Config)

	// Check if we have required client credentials
	if cfg.ClientID == "" {
		return fmt.Errorf("client ID is required for authentication. Set EKA_CLIENT_ID environment variable or use WithClientID() option")
	}

	if cfg.ClientSecret == "" {
		return fmt.Errorf("client secret is required for authentication. Set EKA_CLIENT_SECRET environment variable or use WithClientSecret() option")
	}

	// Create a client credentials provider
	loginRequest := &auth.ClientLoginRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	}

	provider := auth.NewClientCredentialsProvider(c.Auth, loginRequest)
	c.credentialsProvider = provider

	// Get credentials to trigger initial login
	credentials, err := provider.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate with provided credentials: %w", err)
	}

	// Set the authorization token in config for ABDM client
	cfg.SetAuthorizationToken(credentials.AccessToken)

	// Recreate ABDM client with the new token
	c.ABDM = createABDMClient(cfg)

	return nil
}
