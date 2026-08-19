package carecontext

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	abdmecdh "github.com/eka-care/abdm-ecdh/go"
)

// The HIU must be able to decrypt what we push, using only what we send it.
func TestRespondToFetchEncryptsForTheHIU(t *testing.T) {
	e := abdmecdh.New()
	hiu, err := e.GenerateKeyMaterial()
	if err != nil {
		t.Fatal(err)
	}

	var path string
	var body map[string]any
	svc, done := testService(t, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusAccepted)
	})
	defer done()

	bundle := []byte(`{"resourceType":"Bundle","id":"b1"}`)
	evt := &DataFetchEvent{
		TransactionID: "txn-1", OID: "o-1", PartnerPatientID: "pp-1", HIPID: "hip-1",
		keyInfo: keyMaterial{
			CryptoAlg: "ECDH", Curve: "Curve25519", Nonce: hiu.Nonce,
			DHPublicKey: dhPublicKey{KeyValue: hiu.X509PublicKey, Parameters: "p"},
		},
	}

	if err := svc.RespondToFetch(context.Background(), evt,
		[]Entry{{CareContextID: "cc-1", Bundle: bundle}}); err != nil {
		t.Fatalf("RespondToFetch: %v", err)
	}

	if path != "/abdm/v1/hip/care-context/data/on-fetch" {
		t.Errorf("path = %q", path)
	}
	if body["transaction_id"] != "txn-1" {
		t.Errorf("transaction_id = %v", body["transaction_id"])
	}

	entry := body["entries"].([]any)[0].(map[string]any)
	if entry["media"] != "application/fhir+json" || entry["care_context_id"] != "cc-1" {
		t.Errorf("entry = %v", entry)
	}
	sum := sha256.Sum256(bundle)
	if entry["checksum"] != hex.EncodeToString(sum[:]) {
		t.Errorf("checksum = %v", entry["checksum"])
	}

	// Decrypt as the HIU would, using the key material we published.
	ki := body["key_information"].(map[string]any)
	dec, err := e.Decrypt(abdmecdh.DecryptionRequest{
		EncryptedData:       entry["content"].(string),
		SenderNonce:         ki["nonce"].(string),
		RequesterNonce:      hiu.Nonce,
		RequesterPrivateKey: hiu.PrivateKey,
		SenderPublicKey:     ki["dh_public_key"].(map[string]any)["key_value"].(string),
	})
	if err != nil {
		t.Fatalf("HIU could not decrypt: %v", err)
	}
	if dec.DecryptedData != string(bundle) {
		t.Errorf("decrypted = %q, want %q", dec.DecryptedData, bundle)
	}
}

func TestRespondToFetchRejectsEmptyEntries(t *testing.T) {
	svc, done := testService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request should be sent")
	})
	defer done()

	err := svc.RespondToFetch(context.Background(), &DataFetchEvent{TransactionID: "t"}, nil)
	if err == nil {
		t.Error("want error for empty entries")
	}
}
