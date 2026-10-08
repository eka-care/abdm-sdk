package profile

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service handles ABHA profile operations
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a new profile service instance
func NewService(config interfaces.Config) *Service {
	httpClient := http.NewClientFromInterface(config)
	return &Service{
		config: config,
		http:   httpClient,
	}
}

// GetProfile retrieves the user's ABHA profile information
func (s *Service) GetProfile(ctx context.Context, headers interfaces.Headers) (*ProfileResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "GET",
		Path:    "/abdm/v1/profile",
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}

	var response ProfileResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile response: %w", err)
	}

	return &response, nil
}

// GetAssetCard retrieves the ABHA card as a binary image
func (s *Service) GetAssetCard(ctx context.Context, headers interfaces.Headers, req *AssetRequest) (*AssetCardResponse, error) {
	httpReq := &interfaces.HTTPRequest{
		Method:  "GET",
		Path:    "/abdm/v1/profile/asset/card",
		Headers: headers,
	}

	// Add query parameters if provided
	if req != nil && req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		httpReq.Params = queryParams
	}

	resp, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	// The API returns binary image data (image/png), not JSON
	return &AssetCardResponse{
		Data:        resp.Body,
		ContentType: "image/png",
	}, nil
}

// GetAssetQR retrieves the ABHA QR code data as JSON
func (s *Service) GetAssetQR(ctx context.Context, headers interfaces.Headers, req *AssetRequest) (*AssetQRResponse, error) {
	httpReq := &interfaces.HTTPRequest{
		Method:  "GET",
		Path:    "/abdm/v1/profile/asset/qr",
		Headers: headers,
	}

	// Add query parameters if provided
	if req != nil && req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		// The format query parameter for QR is optional and defaults to JSON
		queryParams["format"] = "json"
		httpReq.Params = queryParams
	}

	resp, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var response AssetQRResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal QR asset response: %w", err)
	}

	return &response, nil
}

// UpdateProfile updates the user's ABHA profile information
func (s *Service) UpdateProfile(ctx context.Context, headers interfaces.Headers, req *UpdateProfileRequest) error {
	httpReq := &interfaces.HTTPRequest{
		Method:  "PATCH",
		Path:    "/abdm/v1/profile",
		Headers: headers,
		Body:    req,
	}

	// Add query parameters if OID is provided
	if req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		httpReq.Params = queryParams
	}

	_, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return err
	}

	// Returns 204 No Content on success
	return nil
}

// DeleteProfile deletes the user's ABHA profile and all associated data
func (s *Service) DeleteProfile(ctx context.Context, headers interfaces.Headers, oid string) error {
	httpReq := &interfaces.HTTPRequest{
		Method:  "DELETE",
		Path:    "/abdm/v1/profile",
		Headers: headers,
	}

	// Add query parameters if OID is provided
	if oid != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = oid
		httpReq.Params = queryParams
	}

	_, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return err
	}

	// Returns 204 No Content on success
	return nil
}

// KYCInit initializes the KYC process by requesting an OTP
func (s *Service) KYCInit(ctx context.Context, headers interfaces.Headers, req *KYCInitRequest) (*KYCInitResponse, error) {
	httpReq := &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/profile/kyc/init",
		Headers: headers,
		Body:    req,
	}

	// Add query parameters if OID is provided
	if req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		httpReq.Params = queryParams
	}

	resp, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var response KYCInitResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal KYC init response: %w", err)
	}

	return &response, nil
}

// KYCResend resends the OTP for KYC verification
func (s *Service) KYCResend(ctx context.Context, headers interfaces.Headers, req *KYCResendRequest) (*KYCResendResponse, error) {
	httpReq := &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/profile/kyc/resend",
		Headers: headers,
		Body:    req,
	}

	// Add query parameters if OID is provided
	if req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		httpReq.Params = queryParams
	}

	resp, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var response KYCResendResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal KYC resend response: %w", err)
	}

	return &response, nil
}

// KYCVerify verifies the OTP to complete the KYC process
func (s *Service) KYCVerify(ctx context.Context, headers interfaces.Headers, req *KYCVerifyRequest) (*KYCVerifyResponse, error) {
	httpReq := &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/profile/kyc/verify",
		Headers: headers,
		Body:    req,
	}

	// Add query parameters if OID is provided
	if req.OID != "" {
		queryParams := make(map[string]string)
		queryParams["oid"] = req.OID
		httpReq.Params = queryParams
	}

	resp, err := s.http.Do(ctx, httpReq)
	if err != nil {
		return nil, err
	}

	var response KYCVerifyResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal KYC verify response: %w", err)
	}

	return &response, nil
}

// SessionInit initializes a new session for the user
func (s *Service) SessionInit(ctx context.Context, headers interfaces.Headers, req *SessionInitRequest) (*SessionInitResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/session/init",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var response SessionInitResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session init response: %w", err)
	}

	return &response, nil
}

// SessionVerify verifies the session using OTP
func (s *Service) SessionVerify(ctx context.Context, headers interfaces.Headers, req *SessionVerifyRequest) (*SessionVerifyResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/session/verify",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var response SessionVerifyResponse
	if err := s.http.UnmarshalResponse(resp, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session verify response: %w", err)
	}

	return &response, nil
}
