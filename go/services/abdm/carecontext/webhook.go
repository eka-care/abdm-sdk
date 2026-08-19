package carecontext

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Webhook event names, verbatim from the ABDM Connect docs. These are
// inconsistent upstream — note that the link-confirm event has no care_ prefix.
const (
	EventDataFetch   = "abha.hip_data_fetch"
	EventLinkStatus  = "abha.link_care_context"
	EventDiscover    = "abha.care_context_discover"
	EventLinkInit    = "abha.care_context_discover_link_init"
	EventLinkConfirm = "abha.context_discover_link_confirm"
)

// signatureTolerance is how far a webhook timestamp may drift before rejection.
const signatureTolerance = 3 * time.Minute

var (
	// ErrBadSignature means the payload was not signed by Eka, or was altered.
	ErrBadSignature = errors.New("carecontext: invalid webhook signature")
	// ErrStaleTimestamp means the signature is valid but too old to accept.
	ErrStaleTimestamp = errors.New("carecontext: webhook timestamp outside tolerance")
	// ErrUnknownEvent means the payload is authentic but its event is not one this
	// package handles. Answer 200 — a webhook endpoint must not fail on events it
	// does not yet care about.
	ErrUnknownEvent = errors.New("carecontext: unrecognised webhook event")
)

// timeNow is swapped in tests.
var timeNow = time.Now

// Event is any decoded webhook payload.
type Event interface{ EventName() string }

type dhPublicKey struct {
	Expiry     string `json:"expiry"`
	Parameters string `json:"parameters"`
	KeyValue   string `json:"key_value"`
}

type keyMaterial struct {
	CryptoAlg   string      `json:"crypto_alg"`
	Curve       string      `json:"curve"`
	DHPublicKey dhPublicKey `json:"dh_public_key"`
	Nonce       string      `json:"nonce"`
}

// DataFetchEvent is abha.hip_data_fetch: an HIU wants records. Answer with
// Service.RespondToFetch. The HIU's key material is captured but unexported —
// RespondToFetch uses it so you never handle key material yourself.
type DataFetchEvent struct {
	TransactionID    string
	ABHAAddress      string
	OID              string
	PartnerPatientID string
	HIPID            string
	CareContexts     []string
	HITypes          []HIType

	keyInfo keyMaterial
}

func (*DataFetchEvent) EventName() string { return EventDataFetch }

// LinkStatusEvent is abha.link_care_context: the async outcome of Link.
type LinkStatusEvent struct {
	ABHAAddress      string `json:"abha_address"`
	CareContextID    string `json:"care_context_id"`
	Status           string `json:"status"`
	Error            string `json:"error"`
	RetryCount       int    `json:"retry_count"`
	OID              string `json:"oid"`
	PartnerPatientID string `json:"partner_patient_id"`
	HIPID            string `json:"hip_id"`
}

func (*LinkStatusEvent) EventName() string { return EventLinkStatus }

// Identifier is a verified patient identifier supplied during discovery.
type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DiscoverEvent is abha.care_context_discover: a patient is searching your
// facility for their records. Match them in your own system, conservatively —
// a loose match leaks another patient's records — and answer with OnDiscover.
type DiscoverEvent struct {
	ABHAAddress      string       `json:"abha_address"`
	PatientName      string       `json:"patient_name"`
	Gender           string       `json:"gender"`
	YearOfBirth      int          `json:"year_of_birth"`
	Identifiers      []Identifier `json:"identifiers"`
	RequestID        string       `json:"request_id"`
	TxnID            string       `json:"txn_id"`
	OID              string       `json:"oid"`
	PartnerPatientID string       `json:"partner_patient_id"`
	HIPID            string       `json:"hip_id"`
}

func (*DiscoverEvent) EventName() string { return EventDiscover }

// LinkInitEvent is abha.care_context_discover_link_init: the patient chose care
// contexts to link. Generate an OTP, send it, store it against the ref_num you
// return, then answer with OnLinkInit.
type LinkInitEvent struct {
	ABHAAddress      string          `json:"abha_address"`
	Patient          json.RawMessage `json:"patient"`
	RequestID        string          `json:"request_id"`
	TxnID            string          `json:"txn_id"`
	OID              string          `json:"oid"`
	PartnerPatientID string          `json:"partner_patient_id"`
	HIPID            string          `json:"hip_id"`
}

func (*LinkInitEvent) EventName() string { return EventLinkInit }

// LinkConfirmEvent is abha.context_discover_link_confirm: the patient submitted
// the OTP. Validate Token against what you stored for LinkRefNumber, then answer
// with OnLinkConfirm. This event carries no txn_id — LinkRefNumber is the
// correlation key, and matches the RefNum you sent from OnLinkInit.
type LinkConfirmEvent struct {
	LinkRefNumber    string `json:"linkRefNumber"`
	Token            string `json:"token"`
	RequestID        string `json:"request_id"`
	OID              string `json:"oid"`
	PartnerPatientID string `json:"partner_patient_id"`
	HIPID            string `json:"hip_id"`
}

func (*LinkConfirmEvent) EventName() string { return EventLinkConfirm }

// ParseWebhook verifies an Eka webhook and decodes it into a typed event.
//
// signatureHeader is the raw Eka-Webhook-Signature header; secret is the signing
// key from your webhook subscription; body must be the exact bytes received,
// because the signature covers them verbatim.
//
// Returns ErrBadSignature, ErrStaleTimestamp or ErrUnknownEvent. Treat the first
// two as 401 and ErrUnknownEvent as 200.
func ParseWebhook(body []byte, signatureHeader, secret string) (Event, error) {
	if err := verifySignature(body, signatureHeader, secret); err != nil {
		return nil, err
	}

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("carecontext: decode webhook envelope: %w", err)
	}

	switch env.Event {
	case EventDataFetch:
		var d struct {
			TransactionID    string      `json:"transaction_id"`
			ABHAAddress      string      `json:"abha_address"`
			OID              string      `json:"oid"`
			PartnerPatientID string      `json:"partner_patient_id"`
			HIPID            string      `json:"hip_id"`
			CareContexts     []string    `json:"care_contexts"`
			HITypes          []HIType    `json:"hi_types"`
			KeyInformation   keyMaterial `json:"key_information"`
		}
		if len(env.Data) > 0 {
			if err := json.Unmarshal(env.Data, &d); err != nil {
				return nil, fmt.Errorf("carecontext: decode %s: %w", env.Event, err)
			}
		}
		// The envelope's transaction_id is the fallback when data omits it.
		txn := d.TransactionID
		if txn == "" {
			txn = env.TransactionID
		}
		return &DataFetchEvent{
			TransactionID: txn, ABHAAddress: d.ABHAAddress, OID: d.OID,
			PartnerPatientID: d.PartnerPatientID, HIPID: d.HIPID,
			CareContexts: d.CareContexts, HITypes: d.HITypes, keyInfo: d.KeyInformation,
		}, nil
	case EventLinkStatus:
		return decodeInto(env, &LinkStatusEvent{})
	case EventDiscover:
		return decodeInto(env, &DiscoverEvent{})
	case EventLinkInit:
		return decodeInto(env, &LinkInitEvent{})
	case EventLinkConfirm:
		return decodeInto(env, &LinkConfirmEvent{})
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownEvent, env.Event)
	}
}

func decodeInto[T Event](env envelope, into T) (Event, error) {
	if len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, into); err != nil {
			return nil, fmt.Errorf("carecontext: decode %s: %w", env.Event, err)
		}
	}
	return into, nil
}

func verifySignature(body []byte, header, secret string) error {
	var ts, v1 string
	for _, part := range strings.Split(header, ",") {
		k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		switch k {
		case "t":
			ts = v
		case "v1":
			v1 = v
		}
	}
	if ts == "" || v1 == "" {
		return fmt.Errorf("%w: malformed header", ErrBadSignature)
	}
	// An empty secret would key the HMAC with nothing, which anyone can compute.
	// Refuse rather than authenticate every forgery.
	if secret == "" {
		return fmt.Errorf("%w: no webhook secret configured", ErrBadSignature)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts))
	mac.Write([]byte("."))
	mac.Write(body)
	if !hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(v1)) {
		return ErrBadSignature
	}

	// Timestamp is checked only after the signature, so an attacker cannot use
	// timing here to learn anything about the key.
	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: unparseable timestamp", ErrBadSignature)
	}
	if drift := timeNow().Sub(time.Unix(sec, 0)); drift > signatureTolerance || drift < -signatureTolerance {
		return ErrStaleTimestamp
	}
	return nil
}
