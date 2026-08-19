package carecontext

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"
)

func sign(t *testing.T, body, secret string, ts time.Time) string {
	t.Helper()
	sec := fmt.Sprintf("%d", ts.Unix())
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(sec + "." + body))
	return fmt.Sprintf("t=%s,v1=%s", sec, hex.EncodeToString(m.Sum(nil)))
}

const secret = "whsec"

func TestParseWebhookSignature(t *testing.T) {
	body := `{"service":"abdm","event":"abha.link_care_context","data":{"care_context_id":"cc-1","status":"LINKED"}}`

	t.Run("valid", func(t *testing.T) {
		e, err := ParseWebhook([]byte(body), sign(t, body, secret, time.Now()), secret)
		if err != nil {
			t.Fatalf("ParseWebhook: %v", err)
		}
		st, ok := e.(*LinkStatusEvent)
		if !ok {
			t.Fatalf("event type = %T, want *LinkStatusEvent", e)
		}
		if st.CareContextID != "cc-1" || st.Status != "LINKED" {
			t.Errorf("event = %+v", st)
		}
	})

	t.Run("tampered body", func(t *testing.T) {
		sig := sign(t, body, secret, time.Now())
		_, err := ParseWebhook([]byte(body+" "), sig, secret)
		if !errors.Is(err, ErrBadSignature) {
			t.Errorf("err = %v, want ErrBadSignature", err)
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		_, err := ParseWebhook([]byte(body), sign(t, body, "other", time.Now()), secret)
		if !errors.Is(err, ErrBadSignature) {
			t.Errorf("err = %v, want ErrBadSignature", err)
		}
	})

	t.Run("stale timestamp", func(t *testing.T) {
		old := time.Now().Add(-10 * time.Minute)
		_, err := ParseWebhook([]byte(body), sign(t, body, secret, old), secret)
		if !errors.Is(err, ErrStaleTimestamp) {
			t.Errorf("err = %v, want ErrStaleTimestamp", err)
		}
	})

	t.Run("malformed header", func(t *testing.T) {
		_, err := ParseWebhook([]byte(body), "garbage", secret)
		if !errors.Is(err, ErrBadSignature) {
			t.Errorf("err = %v, want ErrBadSignature", err)
		}
	})
}

// The upstream event names are inconsistent (note the third lacks the care_
// prefix), so each is asserted verbatim.
func TestParseWebhookEventNames(t *testing.T) {
	for _, tc := range []struct {
		event string
		want  string
	}{
		{"abha.hip_data_fetch", "*carecontext.DataFetchEvent"},
		{"abha.link_care_context", "*carecontext.LinkStatusEvent"},
		{"abha.care_context_discover", "*carecontext.DiscoverEvent"},
		{"abha.care_context_discover_link_init", "*carecontext.LinkInitEvent"},
		{"abha.context_discover_link_confirm", "*carecontext.LinkConfirmEvent"},
	} {
		body := fmt.Sprintf(`{"service":"abdm","event":%q,"data":{}}`, tc.event)
		e, err := ParseWebhook([]byte(body), sign(t, body, secret, time.Now()), secret)
		if err != nil {
			t.Errorf("%s: %v", tc.event, err)
			continue
		}
		if got := fmt.Sprintf("%T", e); got != tc.want {
			t.Errorf("%s -> %s, want %s", tc.event, got, tc.want)
		}
		if e.EventName() != tc.event {
			t.Errorf("EventName() = %q, want %q", e.EventName(), tc.event)
		}
	}
}

func TestParseWebhookUnknownEvent(t *testing.T) {
	body := `{"service":"abdm","event":"abha.something_new","data":{}}`
	_, err := ParseWebhook([]byte(body), sign(t, body, secret, time.Now()), secret)
	if !errors.Is(err, ErrUnknownEvent) {
		t.Errorf("err = %v, want ErrUnknownEvent", err)
	}
}

func TestParseWebhookDataFetchKeyMaterial(t *testing.T) {
	body := `{"service":"abdm","event":"abha.hip_data_fetch","data":{
      "transaction_id":"txn-1","abha_address":"p@sbx","oid":"o-1",
      "partner_patient_id":"pp-1","hip_id":"hip-1",
      "care_contexts":["cc-1"],"hi_types":["OPConsultation"],
      "key_information":{"crypto_alg":"ECDH","curve":"Curve25519",
        "dh_public_key":{"key_value":"k","parameters":"p","expiry":"e"},"nonce":"n"}}}`

	e, err := ParseWebhook([]byte(body), sign(t, body, secret, time.Now()), secret)
	if err != nil {
		t.Fatal(err)
	}
	df := e.(*DataFetchEvent)
	if df.TransactionID != "txn-1" || df.HIPID != "hip-1" || len(df.CareContexts) != 1 {
		t.Errorf("event = %+v", df)
	}
	if df.keyInfo.Nonce != "n" || df.keyInfo.DHPublicKey.KeyValue != "k" {
		t.Errorf("key material not captured: %+v", df.keyInfo)
	}
}

// An unconfigured secret must not authenticate anything: with an empty key the
// HMAC is computable by anyone, so a forged payload would otherwise pass.
func TestParseWebhookEmptySecretRejected(t *testing.T) {
	body := `{"service":"abdm","event":"abha.link_care_context","data":{"care_context_id":"cc-1"}}`
	// Signature is genuinely valid for the empty key — still must be rejected.
	_, err := ParseWebhook([]byte(body), sign(t, body, "", time.Now()), "")
	if !errors.Is(err, ErrBadSignature) {
		t.Fatalf("err = %v, want ErrBadSignature", err)
	}
}

func TestParseWebhookDataFetchEmptyData(t *testing.T) {
	body := `{"service":"abdm","event":"abha.hip_data_fetch","transaction_id":"txn-9"}`
	e, err := ParseWebhook([]byte(body), sign(t, body, secret, time.Now()), secret)
	if err != nil {
		t.Fatalf("ParseWebhook: %v", err)
	}
	if df := e.(*DataFetchEvent); df.TransactionID != "txn-9" {
		t.Errorf("transaction id = %q", df.TransactionID)
	}
}
