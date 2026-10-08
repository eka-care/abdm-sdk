package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service represents the utilities service
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a new utilities service
func NewService(config interfaces.Config) *Service {
	httpClient := http.NewClientFromInterface(config)

	return &Service{
		config: config,
		http:   httpClient,
	}
}

// GenerateTransactionID generates a unique transaction ID
func (s *Service) GenerateTransactionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ValidateAadhaarNumber validates an Aadhaar number
func (s *Service) ValidateAadhaarNumber(aadhaar string) error {
	if len(aadhaar) != 12 {
		return fmt.Errorf("Aadhaar number must be 12 digits")
	}

	for _, char := range aadhaar {
		if char < '0' || char > '9' {
			return fmt.Errorf("Aadhaar number must contain only digits")
		}
	}

	return nil
}

// ValidateMobileNumber validates a mobile number
func (s *Service) ValidateMobileNumber(mobile string) error {
	if len(mobile) != 10 {
		return fmt.Errorf("Mobile number must be 10 digits")
	}

	for _, char := range mobile {
		if char < '0' || char > '9' {
			return fmt.Errorf("Mobile number must contain only digits")
		}
	}

	return nil
}

// ValidateABHAAddress validates an ABHA address
func (s *Service) ValidateABHAAddress(address string) error {
	if len(address) < 3 {
		return fmt.Errorf("ABHA address must be at least 3 characters")
	}

	// Check if it ends with @abdm
	if len(address) >= 5 && address[len(address)-5:] != "@abdm" {
		return fmt.Errorf("ABHA address must end with @abdm")
	}

	return nil
}

// FormatDate formats a date for API requests
func (s *Service) FormatDate(year, month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// ParseDate parses a date string from API responses
func (s *Service) ParseDate(dateStr string) (year, month, day int, err error) {
	_, err = fmt.Sscanf(dateStr, "%d-%d-%d", &year, &month, &day)
	return
}

// RetryWithBackoff retries a function with exponential backoff
func (s *Service) RetryWithBackoff(ctx context.Context, fn func() error, maxRetries int, initialDelay time.Duration) error {
	var lastErr error
	delay := initialDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			}
		}
	}

	return lastErr
}

// IsRetryableError checks if an error is retryable
func (s *Service) IsRetryableError(err error) bool {
	// Check for network errors, 5xx status codes, etc.
	// This is a simplified implementation
	return err != nil
}
