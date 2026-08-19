package carecontext

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func capture(t *testing.T) (*Service, func(), *string, *map[string]any) {
	t.Helper()
	path := new(string)
	body := new(map[string]any)
	svc, done := testService(t, func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, body)
		w.WriteHeader(http.StatusNoContent)
	})
	return svc, done, path, body
}

func TestOnDiscover(t *testing.T) {
	svc, done, path, body := capture(t)
	defer done()

	evt := &DiscoverEvent{RequestID: "req-1", TxnID: "txn-1", OID: "o", PartnerPatientID: "pp", HIPID: "hip"}
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
}

func TestOnDiscoverReportsError(t *testing.T) {
	svc, done, _, body := capture(t)
	defer done()

	err := svc.OnDiscover(context.Background(), &DiscoverEvent{RequestID: "r", TxnID: "t"},
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
}

func TestOnLinkInit(t *testing.T) {
	svc, done, path, body := capture(t)
	defer done()

	err := svc.OnLinkInit(context.Background(),
		&LinkInitEvent{RequestID: "req-2", TxnID: "txn-2"},
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
}

// on-link-confirm correlates by request_id alone; the webhook carries no txn_id.
func TestOnLinkConfirm(t *testing.T) {
	svc, done, path, body := capture(t)
	defer done()

	err := svc.OnLinkConfirm(context.Background(),
		&LinkConfirmEvent{RequestID: "req-3", LinkRefNumber: "temp", Token: "111111"},
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
}
