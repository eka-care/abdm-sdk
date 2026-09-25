package registration

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service represents the registration service
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a new registration service
func NewService(config interfaces.Config) *Service {
	httpClient := http.NewClientFromInterface(config)

	return &Service{
		config: config,
		http:   httpClient,
	}
}

// ===============================
// Aadhaar Registration Methods
// ===============================

// AadhaarInit initiates the Aadhaar registration process
func (s *Service) AadhaarInit(ctx context.Context, headers interfaces.Headers, req InitRequest) (*InitResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/init",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result InitResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AadhaarVerify verifies the Aadhaar OTP
func (s *Service) AadhaarVerify(ctx context.Context, headers interfaces.Headers, req VerifyRequest) (*VerifyResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/verify",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result VerifyResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AadhaarResend resends the Aadhaar OTP
func (s *Service) AadhaarResend(ctx context.Context, headers interfaces.Headers, req ResendRequest) (*ResendResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/resend",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result ResendResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AadhaarMobileVerify verifies mobile OTP in Aadhaar registration flow
func (s *Service) AadhaarMobileVerify(ctx context.Context, headers interfaces.Headers, oid string, req MobileVerifyRequest) (*MobileVerifyResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/mobile/verify",
		Headers: headers,
		Body:    req,
		Params:  map[string]string{"oid": oid},
	})
	if err != nil {
		return nil, err
	}

	var result MobileVerifyResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AadhaarMobileResend resends mobile OTP in Aadhaar registration flow
func (s *Service) AadhaarMobileResend(ctx context.Context, headers interfaces.Headers, oid string, req MobileResendRequest) (*MobileResendResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/mobile/resend",
		Headers: headers,
		Body:    req,
		Params:  map[string]string{"oid": oid},
	})
	if err != nil {
		return nil, err
	}

	var result MobileResendResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AadhaarCreatePHR creates a new ABHA address via Aadhaar
func (s *Service) AadhaarCreatePHR(ctx context.Context, headers interfaces.Headers, req CreateRequest) (*CreateResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/aadhaar/create-phr",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result CreateResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ===============================
// Mobile Registration Methods
// ===============================

// MobileInit initiates the mobile registration process
func (s *Service) MobileInit(ctx context.Context, headers interfaces.Headers, req MobileInitRequest) (*MobileInitResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/mobile/init",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result MobileInitResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// MobileVerify verifies the mobile OTP
func (s *Service) MobileVerify(ctx context.Context, headers interfaces.Headers, req MobileVerifyOTPRequest) (*MobileVerifyOTPResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/mobile/verify",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result MobileVerifyOTPResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// MobileResend resends the mobile OTP
func (s *Service) MobileResend(ctx context.Context, headers interfaces.Headers, req MobileResendOTPRequest) (*MobileResendOTPResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/mobile/resend",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result MobileResendOTPResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// MobileCreatePHR creates a new ABHA address via mobile
func (s *Service) MobileCreatePHR(ctx context.Context, headers interfaces.Headers, req MobileCreateRequest) (*MobileCreateResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/mobile/create-phr",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result MobileCreateResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ===============================
// Utility Methods
// ===============================

// CheckAbhaAddressExists checks if an ABHA address already exists
func (s *Service) CheckAbhaAddressExists(ctx context.Context, headers interfaces.Headers, req DoesHealthIdExistRequest) (*DoesHealthIdExistResponse, error) {
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/na/v1/registration/phr/check",
		Headers: headers,
		Body:    req,
	})
	if err != nil {
		return nil, err
	}

	var result DoesHealthIdExistResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// SuggestAbhaAddress gets suggested ABHA addresses based on user details
func (s *Service) SuggestAbhaAddress(ctx context.Context, headers interfaces.Headers, firstName, middleName, lastName, dob, transactionID string) (*SuggestHealthIdResponse, error) {
	params := map[string]string{
		"fn":            firstName,
		"dob":           dob,
		"transactionId": transactionID,
	}
	if middleName != "" {
		params["mn"] = middleName
	}
	if lastName != "" {
		params["ln"] = lastName
	}

	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "GET",
		Path:    "/abdm/na/v1/registration/suggest",
		Headers: headers,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}

	var result SuggestHealthIdResponse
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetPincodeDetails fetches pincode details
func (s *Service) GetPincodeDetails(ctx context.Context, headers interfaces.Headers, pincode string) (*PincodeData, error) {
	path := fmt.Sprintf("/abdm/v1/registration/pincode/%s", pincode)
	resp, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	})
	if err != nil {
		return nil, err
	}

	var result PincodeData
	if err := s.http.UnmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
