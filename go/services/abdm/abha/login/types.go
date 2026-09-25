package login

import "github.com/eka-care/eka-sdk-go/services/abdm/abha"

// InitLoginRequest represents the request for generating login OTP
type InitLoginRequest struct {
	Identifier string      `json:"identifier"`
	Method     LoginMethod `json:"method"` // phr_address, abha_number, mobile, aadhaar_number
}

type LoginMethod string

const (
	LoginMethodPhrAddress    LoginMethod = "phr_address"
	LoginMethodAbhaNumber    LoginMethod = "abha_number"
	LoginMethodMobile        LoginMethod = "mobile"
	LoginMethodAadhaarNumber LoginMethod = "aadhaar_number"
)

// InitLoginResponse represents the response for login OTP generation
type InitLoginResponse struct {
	Hint  string `json:"hint"`
	TxnID string `json:"txn_id"`
}

// VerifyLoginOTPRequest represents the request for verifying login OTP
type VerifyLoginOTPRequest struct {
	OTP   string `json:"otp"`
	TxnID string `json:"txn_id"`
}

// VerifyLoginOTPResponse represents the response for login OTP verification
type VerifyLoginOTPResponse struct {
	AbhaProfiles []AbhaProfile  `json:"abha_profiles"`
	Eka          EkaIDs         `json:"eka"`
	Hint         string         `json:"hint"`
	Profile      Profile        `json:"profile"`
	SkipState    abha.SkipState `json:"skip_state"`
	TxnID        string         `json:"txn_id"`
}

// AbhaProfile represents ABHA profile information
type AbhaProfile struct {
	AbhaAddress string `json:"abha_address"`
	KycVerified string `json:"kyc_verified"`
	Name        string `json:"name"`
}

// EkaIDs represents Eka identification tokens
type EkaIDs struct {
	MinToken string  `json:"min_token"`
	OID      *string `json:"oid,omitempty"`
	UUID     *string `json:"uuid,omitempty"`
}

// Profile represents user profile information
type Profile struct {
	AbhaAddress  string  `json:"abha_address"`
	AbhaNumber   *string `json:"abha_number,omitempty"`
	Address      *string `json:"address,omitempty"`
	DayOfBirth   *int    `json:"day_of_birth,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	Gender       string  `json:"gender"`
	KycVerified  *bool   `json:"kyc_verified,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	MiddleName   *string `json:"middle_name,omitempty"`
	Mobile       *string `json:"mobile,omitempty"`
	MonthOfBirth *int    `json:"month_of_birth,omitempty"`
	Pincode      *string `json:"pincode,omitempty"`
	YearOfBirth  *int    `json:"year_of_birth,omitempty"`
}

// PhrAddressLoginRequest represents the request for PHR address login
type PhrAddressLoginRequest struct {
	PhrAddress string `json:"phr_address"`
	TxnID      string `json:"txn_id"`
}

// PhrAddressLoginResponse represents the response for login
type PhrAddressLoginResponse struct {
	Eka       EkaIDs         `json:"eka"`
	Hint      string         `json:"hint"`
	Profile   Profile        `json:"profile"`
	SkipState abha.SkipState `json:"skip_state"`
	TxnID     string         `json:"txn_id"`
}
