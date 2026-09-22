package auth

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service handles authentication operations for the Eka developer platform
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a new authentication service instance
func NewService(config interfaces.Config) *Service {
	httpClient := http.NewClientFromInterface(config)
	return &Service{
		config: config,
		http:   httpClient,
	}
}

// ClientLogin performs client authentication to get access and refresh tokens
func (s *Service) ClientLogin(ctx context.Context, req *ClientLoginRequest) (*ClientLoginResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method: "POST",
		Path:   "/connect-auth/v1/account/login",
		Body:   req,
	})
	if err != nil {
		return nil, fmt.Errorf("client login request failed: %w", err)
	}

	var response ClientLoginResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal client login response: %w", err)
	}

	return &response, nil
}

// RefreshToken refreshes the access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method: "POST",
		Path:   "/connect-auth/v1/account/refresh",
		Body:   req,
	})
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}

	var response RefreshTokenResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token refresh response: %w", err)
	}

	return &response, nil
}
