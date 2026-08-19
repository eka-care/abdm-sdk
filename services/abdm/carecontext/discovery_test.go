package carecontext

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func capture(t *testing.T) (*Service, func(), *string, *map[string]any, *http.Header) {
	t.Helper()
	path := new(string)
	body := new(map[string]any)
	hdr := new(http.Header)
	svc, done := testService(t, func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.Path
		*hdr = r.Header
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, body)
		w.WriteHeader(http.StatusNoContent)
	})
	return svc, done, path, body, hdr
}

func TestOnDiscover(t *testing.T) {
	svc, done, path, body, hdr := capture(t)
	defer done()

	evt := &DiscoverEvent{RequestID: "req-1", TxnID: "txn-1", OID: "oid-1", PartnerPatientID: "pp-1", HIPID: "hip-1"}
	err := svc.OnDiscover(context.Background(), evt, DiscoverResult{
		Patients: []Patient{{
			ReferenceNumber: "P1", Display: "Gajendra", HIType: HITypeOPConsultation,
			CareContexts: []DiscoveredCareContext{{ReferenceNumber: "CC101", Display: "OP Consult"}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/abdm/v1/care-contexts/on-discover" {
		t.Errorf("path = %q", *path)
	}
	if (*body)["request_id"] != "req-1" || (*body)["txn_id"] != "txn-1" {
		t.Errorf("correlation = %v", *body)
	}
	p := (*body)["patients"].([]any)[0].(map[string]any)
	if p["ref_num"] != "P1" || p["hi_type"] != "OPConsultation" {
		t.Errorf("patient = %v", p)
	}
	if hdr.Get("X-Pt-Id") != "oid-1" || hdr.Get("X-Partner-Pt-Id") != "pp-1" || hdr.Get("X-Hip-Id") != "hip-1" {
		t.Errorf("headers = %v", *hdr)
	}
}

func TestOnDiscoverReportsError(t *testing.T) {
	svc, done, _, body, hdr := capture(t)
	defer done()

	err := svc.OnDiscover(context.Background(), &DiscoverEvent{RequestID: "r", TxnID: "t", OID: "oid-err", PartnerPatientID: "pp-err", HIPID: "hip-err"},
		DiscoverResult{Error: &ErrorDetail{Code: 1000, Message: "no match"}})
	if err != nil {
		t.Fatal(err)
	}
	e, ok := (*body)["error"].(map[string]any)
	if !ok || e["message"] != "no match" {
		t.Errorf("error = %v", (*body)["error"])
	}
	if _, present := (*body)["patients"]; present {
		t.Error("patients should be omitted when reporting an error")
	}
	if hdr.Get("X-Pt-Id") != "oid-err" || hdr.Get("X-Partner-Pt-Id") != "pp-err" || hdr.Get("X-Hip-Id") != "hip-err" {
		t.Errorf("headers = %v", *hdr)
	}
}

func TestOnLinkInit(t *testing.T) {
	svc, done, path, body, hdr := capture(t)
	defer done()

	err := svc.OnLinkInit(context.Background(),
		&LinkInitEvent{RequestID: "req-2", TxnID: "txn-2", OID: "oid-2", PartnerPatientID: "pp-2", HIPID: "hip-2"},
		LinkInitResult{RefNum: "temp", OTPExpiry: "2026-08-19T10:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/abdm/v1/care-contexts/discover/link/on-init" {
		t.Errorf("path = %q", *path)
	}
	if (*body)["ref_num"] != "temp" || (*body)["otp_expiry"] != "2026-08-19T10:00:00Z" {
		t.Errorf("body = %v", *body)
	}
	if hdr.Get("X-Pt-Id") != "oid-2" || hdr.Get("X-Partner-Pt-Id") != "pp-2" || hdr.Get("X-Hip-Id") != "hip-2" {
		t.Errorf("headers = %v", *hdr)
	}
}

// on-link-confirm correlates by request_id alone; the webhook carries no txn_id.
func TestOnLinkConfirm(t *testing.T) {
	svc, done, path, body, hdr := capture(t)
	defer done()

	err := svc.OnLinkConfirm(context.Background(),
		&LinkConfirmEvent{RequestID: "req-3", LinkRefNumber: "temp", Token: "111111", OID: "oid-3", PartnerPatientID: "pp-3", HIPID: "hip-3"},
		LinkConfirmResult{Patients: []Patient{{ReferenceNumber: "P1"}}})
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/abdm/v1/care-contexts/discover/link/on-confirm" {
		t.Errorf("path = %q", *path)
	}
	if (*body)["request_id"] != "req-3" {
		t.Errorf("request_id = %v", (*body)["request_id"])
	}
	if _, present := (*body)["txn_id"]; present {
		t.Error("txn_id must not be sent on link confirm")
	}
	if hdr.Get("X-Pt-Id") != "oid-3" || hdr.Get("X-Partner-Pt-Id") != "pp-3" || hdr.Get("X-Hip-Id") != "hip-3" {
		t.Errorf("headers = %v", *hdr)
	}
}

// An errored result sends the correlation ids and the error only — same as
// OnDiscover and OnLinkConfirm.
func TestOnLinkInitErrorOmitsRefNum(t *testing.T) {
	svc, done, _, body, _ := capture(t)
	defer done()

	evt := &LinkInitEvent{RequestID: "req-1", TxnID: "txn-1"}
	err := svc.OnLinkInit(context.Background(), evt, LinkInitResult{
		RefNum:    "should-not-be-sent",
		OTPExpiry: "2026-01-01T00:00:00Z",
		Error:     &ErrorDetail{Code: 1000, Message: "no match"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := (*body)["ref_num"]; ok {
		t.Errorf("ref_num sent alongside error: %v", *body)
	}
	if _, ok := (*body)["otp_expiry"]; ok {
		t.Errorf("otp_expiry sent alongside error: %v", *body)
	}
	if (*body)["error"] == nil || (*body)["request_id"] != "req-1" {
		t.Errorf("body = %v", *body)
	}
}
