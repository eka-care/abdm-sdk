package carecontext

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	abdmecdh "github.com/eka-care/abdm-ecdh/go"
	"github.com/eka-care/abdm-sdk/go/internal/interfaces"
)

const fhirMedia = "application/fhir+json"

type onFetchEntry struct {
	CareContextID string `json:"care_context_id"`
	Content       string `json:"content"`
	Checksum      string `json:"checksum"`
	Media         string `json:"media"`
}

type onFetchRequest struct {
	TransactionID  string         `json:"transaction_id"`
	PageNumber     int            `json:"page_number"`
	PageCount      int            `json:"page_count"`
	KeyInformation keyMaterial    `json:"key_information"`
	Entries        []onFetchEntry `json:"entries"`
}

// checksum hashes the plaintext bundle.
//
// CONFIRM: the docs say only "Checksum of the non encrypted plain text fhir
// data" without naming an algorithm. ABDM reference implementations have used
// MD5. Confirm with Eka before certification; this is the only place to change.
//
// A wrong checksum fails silently: Eka accepts the push with 202, RespondToFetch
// returns nil, and only the receiving HIU discards the bundle. The records look
// shared from here but never arrive.
func checksum(plaintext []byte) string {
	sum := sha256.Sum256(plaintext)
	return hex.EncodeToString(sum[:])
}

// ourKeyInformation describes the ephemeral key material we generated for one
// response, in the shape the HIU expects.
//
// CONFIRM: the docs do not state what a HIP should send for dh_public_key
// parameters and expiry. We echo the HIU's parameters and set a 24h expiry.
func ourKeyInformation(ours *abdmecdh.KeyMaterial, theirs keyMaterial) keyMaterial {
	params := theirs.DHPublicKey.Parameters
	if params == "" {
		params = "Curve25519/32byte random key"
	}
	curve := theirs.Curve
	if curve == "" {
		curve = "Curve25519"
	}
	return keyMaterial{
		CryptoAlg: "ECDH",
		Curve:     curve,
		Nonce:     ours.Nonce,
		DHPublicKey: dhPublicKey{
			KeyValue:   ours.X509PublicKey,
			Parameters: params,
			Expiry:     timeNow().Add(24 * time.Hour).UTC().Format(time.RFC3339),
		},
	}
}

// RespondToFetch encrypts the given FHIR bundles for the requesting HIU and
// pushes them, answering an abha.hip_data_fetch webhook.
//
// Pass the event exactly as ParseWebhook returned it: it carries the transaction
// id, the HIU's key material, and the identifiers used as request headers. A
// fresh ephemeral key pair is generated per call and discarded, so there is no
// key store to manage.
//
// ABDM retries a fetch it considers unanswered, so this may be called more than
// once for one transaction. Make your own bundle lookup idempotent.
func (s *Service) RespondToFetch(ctx context.Context, e *DataFetchEvent, entries []Entry) error {
	if e == nil {
		return fmt.Errorf("carecontext: nil data fetch event")
	}
	if len(entries) == 0 {
		return fmt.Errorf("carecontext: no entries to send for transaction %q", e.TransactionID)
	}

	ecdh := abdmecdh.New()
	ours, err := ecdh.GenerateKeyMaterial()
	if err != nil {
		return fmt.Errorf("carecontext: generate key material: %w", err)
	}

	out := make([]onFetchEntry, 0, len(entries))
	for _, entry := range entries {
		if len(entry.Bundle) == 0 {
			return fmt.Errorf("carecontext: empty bundle for care context %q", entry.CareContextID)
		}
		enc, err := ecdh.Encrypt(abdmecdh.EncryptionRequest{
			StringToEncrypt:    string(entry.Bundle),
			SenderNonce:        ours.Nonce,
			RequesterNonce:     e.keyInfo.Nonce,
			SenderPrivateKey:   ours.PrivateKey,
			RequesterPublicKey: e.keyInfo.DHPublicKey.KeyValue,
		})
		if err != nil {
			return fmt.Errorf("carecontext: encrypt care context %q: %w", entry.CareContextID, err)
		}
		out = append(out, onFetchEntry{
			CareContextID: entry.CareContextID,
			Content:       enc.EncryptedData,
			Checksum:      checksum(entry.Bundle),
			Media:         fhirMedia,
		})
	}

	// ponytail: single page. Add chunking if bundles ever exceed the payload limit.
	_, err = s.http.Do(ctx, &interfaces.HTTPRequest{
		Method: "POST",
		Path:   "/abdm/v1/hip/care-context/data/on-fetch",
		Headers: interfaces.Headers{
			PatientID:     e.OID,
			PartnerUserID: e.PartnerPatientID,
			HipID:         e.HIPID,
		},
		Body: &onFetchRequest{
			TransactionID:  e.TransactionID,
			PageNumber:     1,
			PageCount:      1,
			KeyInformation: ourKeyInformation(ours, e.keyInfo),
			Entries:        out,
		},
	})
	return err
}
