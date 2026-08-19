// Package carecontext implements ABDM Milestone 2 for a HIP: linking care
// contexts, serving health data on demand, and answering patient-initiated
// discovery.
//
// The package is boilerplate only. It carries bytes and correlation identifiers;
// it never decides what your data means. Matching patients, generating OTPs and
// building FHIR bundles stay in your code.
//
// It is a library, not a server: you own your webhook endpoint and register its
// URL with Eka. Inside your handler, call ParseWebhook and then the responder
// for the event you received.
package carecontext

import "encoding/json"

// HIType is an ABDM health information type.
type HIType string

const (
	HITypeOPConsultation       HIType = "OPConsultation"
	HITypePrescription         HIType = "Prescription"
	HITypeDischargeSummary     HIType = "DischargeSummary"
	HITypeDiagnosticReport     HIType = "DiagnosticReport"
	HITypeImmunizationRecord   HIType = "ImmunizationRecord"
	HITypeHealthDocumentRecord HIType = "HealthDocumentRecord"
	HITypeWellnessRecord       HIType = "WellnessRecord"
)

// CareContext is one logical group of records — a visit, a test, a document.
type CareContext struct {
	CareContextID string   `json:"care_context_id"`
	Display       string   `json:"display"`
	HIType        HIType   `json:"hi_type,omitempty"`
	HITypes       []HIType `json:"hi_types,omitempty"`
}

// LinkRequest is the body of POST /abdm/v1/care-contexts/link.
// OID and PartnerUserID are filled from the request headers when left empty.
type LinkRequest struct {
	ABHAAddress   string        `json:"abha_address"`
	CareContexts  []CareContext `json:"care_contexts"`
	OID           string        `json:"oid,omitempty"`
	PartnerUserID string        `json:"partner_user_id,omitempty"`
}

// Entry is one care context's FHIR bundle, supplied by you and encrypted by
// RespondToFetch. Bundle is ABDM-compliant FHIR R4 JSON.
type Entry struct {
	CareContextID string
	Bundle        []byte
}

// ErrorDetail reports a failure back to ABDM through an on-* responder — no
// patient matched, OTP invalid, and so on.
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DiscoveredCareContext is one unlinked care context offered during discovery.
type DiscoveredCareContext struct {
	ReferenceNumber string `json:"ref_num"`
	Display         string `json:"display"`
}

// Patient is a matched patient and their unlinked care contexts.
type Patient struct {
	ReferenceNumber string                  `json:"ref_num"`
	Display         string                  `json:"display"`
	HIType          HIType                  `json:"hi_type,omitempty"`
	CareContexts    []DiscoveredCareContext `json:"care_contexts"`
}

// envelope is the common webhook wrapper. Data is decoded per event.
type envelope struct {
	Service       string          `json:"service"`
	Event         string          `json:"event"`
	EventTime     int64           `json:"event_time"`
	TransactionID string          `json:"transaction_id"`
	Timestamp     int64           `json:"timestamp"`
	BusinessID    string          `json:"business_id"`
	ClientID      string          `json:"client_id"`
	Data          json.RawMessage `json:"data"`
}
