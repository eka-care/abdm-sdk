package auth

import (
	"context"
	"sync"
	"time"
)

// CredentialsProvider represents a generic credential provider interface
type CredentialsProvider interface {
	// Retrieve retrieves the credentials. Implementations should handle
	// automatic refreshing when possible.
	Retrieve(ctx context.Context) (*Credentials, error)
}

// Credentials represents the authentication credentials for API access
type Credentials struct {
	// AccessToken is the JWT token for API authentication
	AccessToken string

	// RefreshToken can be used to refresh the access token
	RefreshToken string

	// ExpiresAt is when the access token expires
	ExpiresAt time.Time

	// RefreshExpiresAt is when the refresh token expires
	RefreshExpiresAt time.Time

	// Source indicates how the credentials were obtained
	Source string
}

// Expired returns true if the credentials are expired
func (c *Credentials) Expired() bool {
	return time.Now().After(c.ExpiresAt.Add(-5 * time.Minute)) // 5 minute buffer
}

// CanRefresh returns true if the credentials can be refreshed
func (c *Credentials) CanRefresh() bool {
	return c.RefreshToken != "" && time.Now().Before(c.RefreshExpiresAt)
}

// StaticCredentialsProvider provides static credentials
type StaticCredentialsProvider struct {
	credentials *Credentials
}

// NewStaticCredentialsProvider creates a new static credentials provider
func NewStaticCredentialsProvider(accessToken, refreshToken string, expiresIn, refreshExpiresIn int) *StaticCredentialsProvider {
	now := time.Now()
	return &StaticCredentialsProvider{
		credentials: &Credentials{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			ExpiresAt:        now.Add(time.Duration(expiresIn) * time.Second),
			RefreshExpiresAt: now.Add(time.Duration(refreshExpiresIn) * time.Second),
			Source:           "StaticCredentialsProvider",
		},
	}
}

// Retrieve returns the static credentials
func (p *StaticCredentialsProvider) Retrieve(ctx context.Context) (*Credentials, error) {
	return p.credentials, nil
}

// ClientCredentialsProvider handles client-based authentication
type ClientCredentialsProvider struct {
	client  *Service
	request *ClientLoginRequest
	cache   *Credentials
	mu      sync.RWMutex
}

// NewClientCredentialsProvider creates a new client credentials provider
func NewClientCredentialsProvider(client *Service, req *ClientLoginRequest) *ClientCredentialsProvider {
	return &ClientCredentialsProvider{
		client:  client,
		request: req,
	}
}

// Retrieve gets credentials using client authentication, with automatic refresh
func (p *ClientCredentialsProvider) Retrieve(ctx context.Context) (*Credentials, error) {
	p.mu.RLock()
	if p.cache != nil && !p.cache.Expired() {
		creds := p.cache
		p.mu.RUnlock()
		return creds, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check pattern
	if p.cache != nil && !p.cache.Expired() {
		return p.cache, nil
	}

	// Try to refresh if possible
	if p.cache != nil && p.cache.CanRefresh() {
		refreshReq := &RefreshTokenRequest{
			AccessToken:  p.cache.AccessToken,
			RefreshToken: p.cache.RefreshToken,
		}

		resp, err := p.client.RefreshToken(ctx, refreshReq)
		if err == nil {
			p.cache = &Credentials{
				AccessToken:      resp.AccessToken,
				RefreshToken:     resp.RefreshToken,
				ExpiresAt:        time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second),
				RefreshExpiresAt: time.Now().Add(time.Duration(resp.RefreshExpiresIn) * time.Second),
				Source:           "ClientCredentialsProvider(refresh)",
			}
			return p.cache, nil
		}
		// If refresh fails, fall through to login
	}

	// Perform initial login or re-login
	resp, err := p.client.ClientLogin(ctx, p.request)
	if err != nil {
		return nil, err
	}

	p.cache = &Credentials{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		ExpiresAt:        time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second),
		RefreshExpiresAt: time.Now().Add(time.Duration(resp.RefreshExpiresIn) * time.Second),
		Source:           "ClientCredentialsProvider(login)",
	}

	return p.cache, nil
}

// CredentialsCache wraps a credentials provider with caching capabilities
type CredentialsCache struct {
	provider CredentialsProvider
	cache    *Credentials
	mu       sync.RWMutex
}

// NewCredentialsCache creates a new credentials cache
func NewCredentialsCache(provider CredentialsProvider) *CredentialsCache {
	return &CredentialsCache{
		provider: provider,
	}
}

// Retrieve gets credentials from the cache or underlying provider
func (c *CredentialsCache) Retrieve(ctx context.Context) (*Credentials, error) {
	c.mu.RLock()
	if c.cache != nil && !c.cache.Expired() {
		creds := c.cache
		c.mu.RUnlock()
		return creds, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check pattern
	if c.cache != nil && !c.cache.Expired() {
		return c.cache, nil
	}

	creds, err := c.provider.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	c.cache = creds
	return creds, nil
}
