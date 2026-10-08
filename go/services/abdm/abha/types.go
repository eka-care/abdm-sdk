package abha

// SkipState represents the next screen the user should see in ABDM flows
type SkipState string

const (
	// SkipStateAbhaEnd indicates that the ABHA registration or login process is complete,
	// and no further action is needed.
	SkipStateAbhaEnd SkipState = "abha_end"

	// SkipStateConfirmMobileOTP occurs when a user verifies their Aadhaar OTP but enters a mobile number
	// that is not linked to their Aadhaar. In this case, the user should be taken to
	// the mobile verification page.
	SkipStateConfirmMobileOTP SkipState = "confirm_mobile_otp"

	// SkipStateAbhaSelect indicates the user should be taken to the login page, where they have multiple ABHA
	// addresses. They can choose one to proceed with the login.
	SkipStateAbhaSelect SkipState = "abha_select"

	// SkipStateAbhaCreate indicates that no ABHA address has been created yet, so the user should
	// be directed to the ABHA creation page.
	SkipStateAbhaCreate SkipState = "abha_create"
)

// String returns the string representation of the SkipState
func (s SkipState) String() string {
	return string(s)
}

// IsComplete returns true if the skip state indicates the process is complete
func (s SkipState) IsComplete() bool {
	return s == SkipStateAbhaEnd
}

// RequiresUserAction returns true if the skip state requires further user interaction
func (s SkipState) RequiresUserAction() bool {
	return s != SkipStateAbhaEnd
}
