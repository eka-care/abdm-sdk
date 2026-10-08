package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the ABDM API
type APIError struct {
	Code        int    `json:"code"`
	Message     string `json:"error"`
	SourceError *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"source_error,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.SourceError != nil {
		return fmt.Sprintf("API Error %d: %s (Source: %s - %s)", e.Code, e.Message, e.SourceError.Code, e.SourceError.Message)
	}
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// NewAPIError creates a new API error
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewHTTPError creates an API error from HTTP response
func NewHTTPError(statusCode int, body []byte) *APIError {
	return &APIError{
		Code:    statusCode,
		Message: string(body),
	}
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		// Retry on 5xx errors and rate limiting
		return apiErr.Code >= 500 || apiErr.Code == http.StatusTooManyRequests
	}
	return false
}
