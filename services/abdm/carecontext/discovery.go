package carecontext

import (
	"context"
	"fmt"

	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// DiscoverResult answers a discovery request: the patients you matched, or an
// error explaining why you could not.
type DiscoverResult struct {
	Patients []Patient
	Error    *ErrorDetail
}

// LinkInitResult reports that you sent an OTP. RefNum is your reference for this
// linking attempt; it returns as LinkConfirmEvent.LinkRefNumber, so store your
// OTP against it. OTPExpiry is an ISO-8601 timestamp.
type LinkInitResult struct {
	RefNum    string
	OTPExpiry string
	Error     *ErrorDetail
}

// LinkConfirmResult reports the outcome of OTP validation: the care contexts to
// link, or an error if the OTP was wrong or expired.
type LinkConfirmResult struct {
	Patients []Patient
	Error    *ErrorDetail
}

type onDiscoverBody struct {
	RequestID string       `json:"request_id"`
	TxnID     string       `json:"txn_id"`
	Patients  []Patient    `json:"patients,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

type onLinkInitBody struct {
	RequestID string       `json:"request_id"`
	TxnID     string       `json:"txn_id"`
	RefNum    string       `json:"ref_num,omitempty"`
	OTPExpiry string       `json:"otp_expiry,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

type onLinkConfirmBody struct {
	RequestID string       `json:"request_id"`
	Patients  []Patient    `json:"patients,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

// OnDiscover answers an abha.care_context_discover webhook with the unlinked
// care contexts belonging to the patient described in the event.
//
// Matching is yours: query your own patient records using the event's name,
// gender, year of birth and identifiers. Match conservatively — a loose match
// discloses another patient's records to the requester. Set Result.Error when
// nothing matched.
func (s *Service) OnDiscover(ctx context.Context, e *DiscoverEvent, res DiscoverResult) error {
	if e == nil {
		return fmt.Errorf("carecontext: nil discover event")
	}
	body := onDiscoverBody{RequestID: e.RequestID, TxnID: e.TxnID, Error: res.Error}
	if res.Error == nil {
		body.Patients = res.Patients
	}
	return s.postDiscovery(ctx, "/abdm/v1/care-contexts/on-discover", e.OID, e.PartnerPatientID, e.HIPID, &body)
}

// OnLinkInit answers an abha.care_context_discover_link_init webhook, telling
// ABDM you have dispatched an OTP.
//
// Generating the OTP, delivering it, and storing it against res.RefNum are all
// yours — this package holds no state. The OTP must survive until the matching
// link-confirm webhook arrives.
func (s *Service) OnLinkInit(ctx context.Context, e *LinkInitEvent, res LinkInitResult) error {
	if e == nil {
		return fmt.Errorf("carecontext: nil link init event")
	}
	body := onLinkInitBody{
		RequestID: e.RequestID, TxnID: e.TxnID,
		RefNum: res.RefNum, OTPExpiry: res.OTPExpiry, Error: res.Error,
	}
	return s.postDiscovery(ctx, "/abdm/v1/care-contexts/discover/link/on-init", e.OID, e.PartnerPatientID, e.HIPID, &body)
}

// OnLinkConfirm answers an abha.context_discover_link_confirm webhook after you
// have validated the OTP.
//
// Validate e.Token against whatever you stored for e.LinkRefNumber. On success
// send the care contexts to link; on failure set res.Error. This endpoint
// correlates by request_id alone — the webhook carries no txn_id.
func (s *Service) OnLinkConfirm(ctx context.Context, e *LinkConfirmEvent, res LinkConfirmResult) error {
	if e == nil {
		return fmt.Errorf("carecontext: nil link confirm event")
	}
	body := onLinkConfirmBody{RequestID: e.RequestID, Error: res.Error}
	if res.Error == nil {
		body.Patients = res.Patients
	}
	return s.postDiscovery(ctx, "/abdm/v1/care-contexts/discover/link/on-confirm", e.OID, e.PartnerPatientID, e.HIPID, &body)
}

func (s *Service) postDiscovery(ctx context.Context, path, oid, partnerPatientID, hipID string, body any) error {
	_, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    path,
		Headers: interfaces.Headers{PatientID: oid, PartnerUserID: partnerPatientID, HipID: hipID},
		Body:    body,
	})
	return err
}
