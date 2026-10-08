package login

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service handles ABHA login operations
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a new login service instance
func NewService(config interfaces.Config) *Service {
	httpClient := http.NewClientFromInterface(config)
	return &Service{
		config: config,
		http:   httpClient,
	}
}

// LoginInit generates OTP for login with different identifier methods
func (s *Service) LoginInit(ctx context.Context, headers interfaces.Headers, req *InitLoginRequest) (*InitLoginResponse, error) {

	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/profile/login/init",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var response InitLoginResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// LoginVerify verifies the login OTP
func (s *Service) LoginVerify(ctx context.Context, headers interfaces.Headers, req *VerifyLoginOTPRequest) (*VerifyLoginOTPResponse, error) {

	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/profile/login/verify",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var response VerifyLoginOTPResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// LoginWithPHRAddress handles login using PHR address
func (s *Service) LoginWithPHRAddress(ctx context.Context, headers interfaces.Headers, req *PhrAddressLoginRequest) (*PhrAddressLoginResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/profile/login/phr",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var response PhrAddressLoginResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
