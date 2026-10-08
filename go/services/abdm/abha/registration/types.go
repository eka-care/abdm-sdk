package registration

import "github.com/eka-care/eka-sdk-go/services/abdm/abha"

// ===============================
// Aadhaar Registration Types
// ===============================

// InitRequest represents the request to initiate Aadhaar registration
type InitRequest struct {
	AadhaarNumber string `json:"aadhaar_number"`
}

// InitResponse represents the response from Aadhaar init
type InitResponse struct {
	TxnID string  `json:"txn_id"`
	Hint  *string `json:"hint,omitempty"`
}

// VerifyRequest represents the request to verify Aadhaar OTP
type VerifyRequest struct {
	TxnID  string `json:"txn_id"`
	OTP    string `json:"otp"`
	Mobile string `json:"mobile"`
}

// VerifyResponse represents the response from Aadhaar verify
type VerifyResponse struct {
	TxnID        string           `json:"txn_id"`
	SkipState    abha.SkipState   `json:"skip_state"`
	Profile      *ProfileResponse `json:"profile,omitempty"`
	Token        *string          `json:"token,omitempty"`
	RefreshToken *string          `json:"refresh_token,omitempty"`
	Eka          *EkaIds          `json:"eka,omitempty"`
	Hint         *string          `json:"hint,omitempty"`
}

// ResendRequest represents the request to resend Aadhaar OTP
type ResendRequest struct {
	TxnID string `json:"txn_id"`
}

// ResendResponse represents the response from Aadhaar resend OTP
type ResendResponse struct {
	TxnID string  `json:"txn_id"`
	Hint  *string `json:"hint,omitempty"`
}

// MobileVerifyRequest represents the request to verify mobile OTP in Aadhaar flow
type MobileVerifyRequest struct {
	TxnID string `json:"txn_id"`
	OTP   string `json:"otp"`
}

// MobileVerifyResponse represents the response from mobile verify in Aadhaar flow
type MobileVerifyResponse struct {
	TxnID        string           `json:"txn_id"`
	SkipState    abha.SkipState   `json:"skip_state"`
	Profile      *ProfileResponse `json:"profile,omitempty"`
	Token        *string          `json:"token,omitempty"`
	RefreshToken *string          `json:"refresh_token,omitempty"`
	Eka          *EkaIds          `json:"eka,omitempty"`
	Hint         *string          `json:"hint,omitempty"`
}

// MobileResendRequest represents the request to resend mobile OTP in Aadhaar flow
type MobileResendRequest struct {
	TxnID string `json:"txn_id"`
}

// MobileResendResponse represents the response from mobile resend OTP in Aadhaar flow
type MobileResendResponse struct {
	TxnID string  `json:"txn_id"`
	Hint  *string `json:"hint,omitempty"`
}

// CreateRequest represents the request to create ABHA address via Aadhaar
type CreateRequest struct {
	TxnID       string `json:"txn_id"`
	AbhaAddress string `json:"abha_address"`
}

// CreateResponse represents the response from ABHA address creation via Aadhaar
type CreateResponse struct {
	TxnID        string           `json:"txn_id"`
	SkipState    abha.SkipState   `json:"skip_state"`
	Profile      *ProfileResponse `json:"profile,omitempty"`
	Token        *string          `json:"token,omitempty"`
	RefreshToken *string          `json:"refresh_token,omitempty"`
	Eka          *EkaIds          `json:"eka,omitempty"`
	Hint         *string          `json:"hint,omitempty"`
}

// ===============================
// Mobile Registration Types
// ===============================

// MobileInitRequest represents the request to initiate mobile registration
type MobileInitRequest struct {
	MobileNumber string `json:"mobile_number"`
}

// MobileInitResponse represents the response from mobile init
type MobileInitResponse struct {
	TxnID string  `json:"txn_id"`
	Hint  *string `json:"hint,omitempty"`
}

// MobileVerifyOTPRequest represents the request to verify mobile OTP
type MobileVerifyOTPRequest struct {
	TxnID string `json:"txn_id"`
	OTP   string `json:"otp"`
}

// MobileVerifyOTPResponse represents the response from mobile verify OTP
type MobileVerifyOTPResponse struct {
	TxnID        string              `json:"txn_id"`
	SkipState    abha.SkipState      `json:"skip_state"`
	AbhaProfiles []VerifyAbhaProfile `json:"abha_profiles,omitempty"`
	Eka          *EkaIds             `json:"eka,omitempty"`
}

// MobileResendOTPRequest represents the request to resend mobile OTP
type MobileResendOTPRequest struct {
	TxnID string `json:"txn_id"`
}

// MobileResendOTPResponse represents the response from mobile resend OTP
type MobileResendOTPResponse struct {
	TxnID string `json:"txn_id"`
	Hint  string `json:"hint"`
}

// MobileCreateRequest represents the request to create ABHA address via mobile
type MobileCreateRequest struct {
	TxnID       string                `json:"txn_id"`
	AbhaAddress string                `json:"abha_address"`
	Profile     ProfileDetailsRequest `json:"profile"`
}

// MobileCreateResponse represents the response from ABHA address creation via mobile
type MobileCreateResponse struct {
	SkipState    abha.SkipState   `json:"skip_state"`
	Success      bool             `json:"success"`
	Profile      *ProfileResponse `json:"profile,omitempty"`
	Token        *string          `json:"token,omitempty"`
	RefreshToken *string          `json:"refresh_token,omitempty"`
	Eka          *EkaIds          `json:"eka,omitempty"`
}

// ===============================
// Common Types
// ===============================

// ProfileResponse represents the ABHA profile
type ProfileResponse struct {
	AbhaAddress  string  `json:"abha_address"`
	AbhaNumber   *string `json:"abha_number,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	MiddleName   *string `json:"middle_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Gender       string  `json:"gender"`
	YearOfBirth  *int    `json:"year_of_birth,omitempty"`
	MonthOfBirth *int    `json:"month_of_birth,omitempty"`
	DayOfBirth   *int    `json:"day_of_birth,omitempty"`
	Mobile       *string `json:"mobile,omitempty"`
	Address      *string `json:"address,omitempty"`
	Pincode      *string `json:"pincode,omitempty"`
	KycVerified  *bool   `json:"kyc_verified,omitempty"`
}

// ProfileDetailsRequest represents profile details for requests
type ProfileDetailsRequest struct {
	FirstName    string  `json:"first_name"`
	MiddleName   *string `json:"middle_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Gender       string  `json:"gender"`
	YearOfBirth  int     `json:"year_of_birth"`
	MonthOfBirth int     `json:"month_of_birth"`
	DayOfBirth   int     `json:"day_of_birth"`
	Address      *string `json:"address,omitempty"`
	Pincode      string  `json:"pincode"`
}

// EkaIds represents the Eka platform identifiers
type EkaIds struct {
	OID      *string `json:"oid,omitempty"`
	UUID     *string `json:"uuid,omitempty"`
	MinToken string  `json:"min_token"`
}

// VerifyAbhaProfile represents ABHA profile info during verification
type VerifyAbhaProfile struct {
	AbhaAddress string `json:"abha_address"`
	Name        string `json:"name"`
	KycVerified string `json:"kyc_verified"`
}

// DoesHealthIdExistRequest represents the request to check if health ID exists
type DoesHealthIdExistRequest struct {
	AbhaAddress string `json:"abha_address"`
}

// DoesHealthIdExistResponse represents the response for health ID existence check
type DoesHealthIdExistResponse struct {
	Exists bool `json:"exists"`
}

// SuggestHealthIdResponse represents the response for suggested ABHA addresses
type SuggestHealthIdResponse struct {
	Suggestions []string `json:"suggestions"`
}

// PincodeData represents pincode resolution data
type PincodeData struct {
	Pincode   string `json:"pincode"`
	DistCode  string `json:"dist_code"`
	DistName  string `json:"dist_name"`
	StateCode string `json:"state_code"`
	StateName string `json:"state_name"`
}
