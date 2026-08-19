package carecontext

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/config"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

func testService(t *testing.T, h http.HandlerFunc) (*Service, func()) {
	t.Helper()
	srv := httptest.NewServer(h)
	cfg := &config.Config{BaseURL: srv.URL, Timeout: 5 * time.Second, UserAgent: "test"}
	return NewService(cfg), srv.Close
}

func TestLinkSendsContractShape(t *testing.T) {
	var path string
	var hdr http.Header
	var body map[string]any

	svc, done := testService(t, func(w http.ResponseWriter, r *http.Request) {
		path, hdr = r.URL.Path, r.Header
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusAccepted)
	})
	defer done()

	err := svc.Link(context.Background(),
		interfaces.Headers{PatientID: "oid-1", PartnerUserID: "pt-1", HipID: "hip-1"},
		&LinkRequest{
			ABHAAddress: "patient@sbx",
			CareContexts: []CareContext{
				{CareContextID: "visit-123", Display: "OP Consult", HIType: HITypeOPConsultation},
			},
		})
	if err != nil {
		t.Fatalf("Link: %v", err)
	}

	if path != "/abdm/v1/care-contexts/link" {
		t.Errorf("path = %q", path)
	}
	if hdr.Get("X-Pt-Id") != "oid-1" || hdr.Get("X-Partner-Pt-Id") != "pt-1" || hdr.Get("X-Hip-Id") != "hip-1" {
		t.Errorf("headers = %v", hdr)
	}
	// oid and partner_user_id appear in the body as well as the headers; the SDK
	// fills them from the headers so the caller states them once.
	if body["oid"] != "oid-1" || body["partner_user_id"] != "pt-1" {
		t.Errorf("body identifiers = %v", body)
	}
	if body["abha_address"] != "patient@sbx" {
		t.Errorf("abha_address = %v", body["abha_address"])
	}
	cc := body["care_contexts"].([]any)[0].(map[string]any)
	if cc["care_context_id"] != "visit-123" || cc["hi_type"] != "OPConsultation" {
		t.Errorf("care_contexts[0] = %v", cc)
	}
}
