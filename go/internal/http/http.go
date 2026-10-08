package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Client represents the HTTP client
type Client struct {
	baseURL    string
	apiKey     string
	userAgent  string
	timeout    time.Duration
	httpClient *http.Client
	middleware []interfaces.Middleware
}

// Config represents HTTP client configuration
type Config struct {
	BaseURL    string
	APIKey     string
	UserAgent  string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// NewClient creates a new HTTP client
func NewClient(cfg *Config) *Client {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	return &Client{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		userAgent:  cfg.UserAgent,
		timeout:    cfg.Timeout,
		httpClient: httpClient,
	}
}

// NewClientFromInterface creates a new HTTP client from an interface
func NewClientFromInterface(config interfaces.Config) *Client {
	httpClient := config.GetHTTPClient()
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.GetTimeout(),
		}
	}

	return &Client{
		baseURL:    config.GetBaseURL(),
		apiKey:     config.GetAPIKey(),
		userAgent:  config.GetUserAgent(),
		timeout:    config.GetTimeout(),
		httpClient: httpClient,
	}
}

// AddMiddleware adds middleware to the client
func (c *Client) AddMiddleware(middleware interfaces.Middleware) {
	c.middleware = append(c.middleware, middleware)
}

// Do performs an HTTP request
func (c *Client) Do(ctx context.Context, req *interfaces.HTTPRequest) (*interfaces.HTTPResponse, error) {
	// Build URL
	u, err := url.Parse(c.baseURL + req.Path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add query parameters
	if req.Params != nil {
		q := u.Query()
		for k, v := range req.Params {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	// Prepare request body
	var reqBody io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)

	if req.Headers.PatientID != "" {
		httpReq.Header.Set("X-Pt-Id", req.Headers.PatientID)
	}
	if req.Headers.PartnerUserID != "" {
		httpReq.Header.Set("X-Partner-Pt-Id", req.Headers.PartnerUserID)
	}
	if req.Headers.HipID != "" {
		httpReq.Header.Set("X-Hip-Id", req.Headers.HipID)
	}

	// Apply middleware
	transport := c.httpClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// Apply custom middleware
	for _, mw := range c.middleware {
		transport = mw(transport)
	}

	// Create client with custom transport
	client := &http.Client{
		Transport: transport,
		Timeout:   c.timeout,
	}

	// Make the request
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, &APIError{
			Code:    resp.StatusCode,
			Message: errorResp.String(),
		}
	}

	return &interfaces.HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// UnmarshalResponse unmarshals the response body into the given type
func (c *Client) UnmarshalResponse(resp *interfaces.HTTPResponse, v interface{}) error {
	if len(resp.Body) == 0 {
		return nil
	}
	return json.Unmarshal(resp.Body, v)
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
	SourceError *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"source_error,omitempty"`
}

func (e *ErrorResponse) String() string {
	if e.SourceError != nil {
		return fmt.Sprintf("Error %d: %s (Source: %s - %s)", e.Code, e.Error, e.SourceError.Code, e.SourceError.Message)
	}
	return fmt.Sprintf("Error %d: %s", e.Code, e.Error)
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
