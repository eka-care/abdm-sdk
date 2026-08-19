package profile

// ProfileResponse represents user profile information
type ProfileResponse struct {
	AbhaAddress  string  `json:"abha_address"`
	AbhaNumber   *string `json:"abha_number,omitempty"`
	Name         *string `json:"name,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	MiddleName   *string `json:"middle_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Gender       string  `json:"gender"`
	DateOfBirth  *string `json:"date_of_birth,omitempty"`
	YearOfBirth  *int    `json:"year_of_birth,omitempty"`
	MonthOfBirth *int    `json:"month_of_birth,omitempty"`
	DayOfBirth   *int    `json:"day_of_birth,omitempty"`
	Email        *string `json:"email,omitempty"`
	Mobile       *string `json:"mobile,omitempty"`
	Address      *string `json:"address,omitempty"`
	Pincode      *string `json:"pincode,omitempty"`
	KycVerified  *bool   `json:"kyc_verified,omitempty"`
}

// AssetCardResponse represents ABHA card asset response
// Note: The API returns image/png (binary data), not JSON
type AssetCardResponse struct {
	Data        []byte `json:"-"` // Raw image data (PNG format)
	ContentType string `json:"-"` // Content type (image/png)
}

// AssetQRResponse represents ABHA QR code asset response
type AssetQRResponse struct {
	AbhaAddress  string `json:"abha_address,omitempty"`
	DistName     string `json:"dist name,omitempty"`
	DistrictLGD  string `json:"distlgd,omitempty"`
	DistrictName string `json:"district_name,omitempty"`
	DOB          string `json:"dob,omitempty"`
	Gender       string `json:"gender,omitempty"`
	HID          string `json:"hid,omitempty"`
	HIDN         string `json:"hidn,omitempty"`
	Mobile       string `json:"mobile,omitempty"`
	Name         string `json:"name,omitempty"`
	PHR          string `json:"phr,omitempty"`
	StateName    string `json:"state name,omitempty"`
	StateLGD     string `json:"statelgd,omitempty"`
}

// SessionInitRequest represents session initialization request
type SessionInitRequest struct {
	AbhaAddress string `json:"abha_address"`
}

// SessionInitResponse represents session initialization response
type SessionInitResponse struct {
	TxnID string `json:"txn_id"`
}

// SessionVerifyRequest represents session verification request
type SessionVerifyRequest struct {
	OTP   string `json:"otp"`
	TxnID string `json:"txn_id"`
}

// SessionVerifyResponse represents session verification response
type SessionVerifyResponse struct {
	Token        string  `json:"token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

// AssetRequest represents request parameters for asset generation
type AssetRequest struct {
	OID string `json:"oid,omitempty"` // OID is used to identify the user
}

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	OID          string `json:"-"` // OID is passed as query parameter, not in body
	Address      string `json:"address,omitempty"`
	DayOfBirth   int    `json:"day_of_birth,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	Gender       string `json:"gender,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	MiddleName   string `json:"middle_name,omitempty"`
	MonthOfBirth int    `json:"month_of_birth,omitempty"`
	Pincode      string `json:"pincode,omitempty"`
	YearOfBirth  int    `json:"year_of_birth,omitempty"`
}

// KYCInitRequest represents the request body for KYC initialization
type KYCInitRequest struct {
	OID        string `json:"-"`            // OID is passed as query parameter, not in body
	Identifier string `json:"identifier"`   // The identifier for KYC
	Method     string `json:"method"`       // KYC method (e.g., "abha-number")
	UserXToken string `json:"user_x_token"` // User X token
}

// KYCInitResponse represents the response from KYC initialization
type KYCInitResponse struct {
	TxnID string `json:"txn_id"`
}

// KYCResendRequest represents the request body for KYC OTP resend
type KYCResendRequest struct {
	OID   string `json:"-"`      // OID is passed as query parameter, not in body
	TxnID string `json:"txn_id"` // Transaction ID from KYC init
}

// KYCResendResponse represents the response from KYC OTP resend
type KYCResendResponse struct {
	TxnID string `json:"txn_id"`
}

// KYCVerifyRequest represents the request body for KYC verification
type KYCVerifyRequest struct {
	OID        string `json:"-"`            // OID is passed as query parameter, not in body
	OTP        string `json:"otp"`          // OTP for verification
	TxnID      string `json:"txn_id"`       // Transaction ID from KYC init
	UserXToken string `json:"user_x_token"` // User X token
}

// KYCVerifyResponse represents the response from KYC verification
type KYCVerifyResponse struct {
	TxnID string `json:"txn_id"`
}
