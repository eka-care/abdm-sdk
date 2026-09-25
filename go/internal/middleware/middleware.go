package middleware

import (
	"net/http"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// RetryMiddleware creates a retry middleware
func RetryMiddleware(maxRetries int, backoff time.Duration) interfaces.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &retryTransport{
			next:       next,
			maxRetries: maxRetries,
			backoff:    backoff,
		}
	}
}

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware(logger interfaces.Logger) interfaces.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &loggingTransport{
			next:   next,
			logger: logger,
		}
	}
}

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authFunc func(*http.Request) error) interfaces.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &authTransport{
			next:     next,
			authFunc: authFunc,
		}
	}
}

// MetricsMiddleware creates a metrics middleware
func MetricsMiddleware(metrics interfaces.MetricsCollector) interfaces.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &metricsTransport{
			next:    next,
			metrics: metrics,
		}
	}
}

// retryTransport implements retry logic
type retryTransport struct {
	next       http.RoundTripper
	maxRetries int
	backoff    time.Duration
}

func (r *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		resp, err := r.next.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		}

		if attempt < r.maxRetries {
			time.Sleep(r.backoff * time.Duration(attempt+1))
		}
	}

	return nil, lastErr
}

// loggingTransport implements logging
type loggingTransport struct {
	next   http.RoundTripper
	logger interfaces.Logger
}

func (l *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	l.logger.LogRequest(req)

	resp, err := l.next.RoundTrip(req)

	duration := time.Since(start)
	l.logger.LogResponse(resp, err, duration)

	return resp, err
}

// authTransport implements authentication
type authTransport struct {
	next     http.RoundTripper
	authFunc func(*http.Request) error
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := a.authFunc(req); err != nil {
		return nil, err
	}
	return a.next.RoundTrip(req)
}

// metricsTransport implements metrics collection
type metricsTransport struct {
	next    http.RoundTripper
	metrics interfaces.MetricsCollector
}

func (m *metricsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := m.next.RoundTrip(req)

	duration := time.Since(start)
	m.metrics.RecordRequest(req, resp, err, duration)

	return resp, err
}
