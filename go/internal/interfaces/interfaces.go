package interfaces

import (
	"context"
	"net/http"
	"time"
)

// Config represents the configuration interface
type Config interface {
	GetBaseURL() string
	GetAPIKey() string
	GetTimeout() time.Duration
	GetMaxRetries() int
	GetUserAgent() string
	GetLogLevel() string
	GetHTTPClient() *http.Client
	GetDisableSSL() bool
	GetRegion() string
	GetRetryMode() string
	GetMaxBackoffDelay() time.Duration
	GetRequestTimeout() time.Duration
	GetResponseTimeout() time.Duration
	GetConnectionTimeout() time.Duration
}

// HTTPClient represents the HTTP client interface
type HTTPClient interface {
	Do(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error)
	UnmarshalResponse(resp *HTTPResponse, v interface{}) error
	AddMiddleware(middleware Middleware)
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string
	Path    string
	Headers Headers
	Body    interface{}
	Params  map[string]string
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Headers represents request headers
type Headers struct {
	PatientID     string
	PartnerUserID string
	HipID         string
}

// Middleware represents a middleware function
type Middleware func(next http.RoundTripper) http.RoundTripper

// Logger represents a logger interface
type Logger interface {
	LogRequest(*http.Request)
	LogResponse(*http.Response, error, time.Duration)
}

// MetricsCollector represents a metrics collector interface
type MetricsCollector interface {
	RecordRequest(*http.Request, *http.Response, error, time.Duration)
}
