// Command m2-hip is a runnable ABDM M2 (HIP) reference: it links a care context
// and serves the webhooks Eka sends back.
//
//	export EKA_CLIENT_ID=... EKA_CLIENT_SECRET=... EKA_WEBHOOK_SECRET=...
//	go run ./examples/m2-hip
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	ekasdk "github.com/eka-care/eka-sdk-go"
	"github.com/eka-care/eka-sdk-go/services/abdm/carecontext"
)

// otpStore stands in for the integrator's own storage. The SDK holds no state,
// so an OTP must survive from the link-init webhook to the link-confirm one.
type otpStore struct {
	mu   sync.Mutex
	otps map[string]string // ref_num -> otp
}

func (s *otpStore) put(ref, otp string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.otps[ref] = otp
}

func (s *otpStore) valid(ref, otp string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.otps[ref] != "" && s.otps[ref] == otp
}

func main() {
	ctx := context.Background()

	client := ekasdk.NewFromEnv()
	if err := client.Login(ctx); err != nil {
		log.Fatalf("authentication failed: %v", err)
	}
	cc := client.ABDM.CareContexts()

	secret := os.Getenv("EKA_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("EKA_WEBHOOK_SECRET is required to verify webhook signatures")
	}
	store := &otpStore{otps: map[string]string{}}

	// Outbound: link a care context when a visit closes.
	headers := carecontext.Headers{PatientID: "eka-oid", PartnerUserID: "your-patient-id", HipID: "your-hip-id"}
	if err := cc.Link(ctx, headers, &carecontext.LinkRequest{
		ABHAAddress: "patient@sbx",
		CareContexts: []carecontext.CareContext{{
			CareContextID: "visit-123",
			Display:       "OP Consult " + time.Now().Format("2 Jan"),
			HIType:        carecontext.HITypeOPConsultation,
		}},
	}); err != nil {
		log.Printf("link failed: %v", err)
	}

	// Inbound: your endpoint, your framework, your status codes.
	http.HandleFunc("/webhooks/eka", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unreadable body", http.StatusBadRequest)
			return
		}

		evt, err := carecontext.ParseWebhook(body, r.Header.Get("Eka-Webhook-Signature"), secret)
		switch {
		case errors.Is(err, carecontext.ErrUnknownEvent):
			// Not ours yet — acknowledge so the endpoint keeps working.
			w.WriteHeader(http.StatusOK)
			return
		case err != nil:
			http.Error(w, "rejected", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		switch e := evt.(type) {
		case *carecontext.DataFetchEvent:
			// Look up each requested care context and build its FHIR bundle.
			// ABDM retries, so make this lookup idempotent.
			entries := make([]carecontext.Entry, 0, len(e.CareContexts))
			for _, id := range e.CareContexts {
				entries = append(entries, carecontext.Entry{CareContextID: id, Bundle: fhirBundleFor(id)})
			}
			err = cc.RespondToFetch(ctx, e, entries)

		case *carecontext.LinkStatusEvent:
			log.Printf("care context %s: %s %s", e.CareContextID, e.Status, e.Error)

		case *carecontext.DiscoverEvent:
			// Match conservatively: a loose match leaks another patient's records.
			err = cc.OnDiscover(ctx, e, carecontext.DiscoverResult{
				Patients: matchPatients(e),
			})

		case *carecontext.LinkInitEvent:
			ref, otp := "ref-"+e.TxnID, "123456"
			store.put(ref, otp) // send otp to the patient's registered mobile here
			err = cc.OnLinkInit(ctx, e, carecontext.LinkInitResult{
				RefNum:    ref,
				OTPExpiry: time.Now().Add(10 * time.Minute).UTC().Format(time.RFC3339),
			})

		case *carecontext.LinkConfirmEvent:
			if !store.valid(e.LinkRefNumber, e.Token) {
				err = cc.OnLinkConfirm(ctx, e, carecontext.LinkConfirmResult{
					Error: &carecontext.ErrorDetail{Code: 1402, Message: "invalid OTP"},
				})
				break
			}
			err = cc.OnLinkConfirm(ctx, e, carecontext.LinkConfirmResult{
				Patients: []carecontext.Patient{{
					ReferenceNumber: e.LinkRefNumber,
					CareContexts:    []carecontext.DiscoveredCareContext{{ReferenceNumber: "CC101", Display: "OP Consult"}},
				}},
			})
		}

		if err != nil {
			// Non-2xx makes ABDM retry, which is the recovery path.
			log.Printf("handling %s: %v", evt.EventName(), err)
			http.Error(w, "retry", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080/webhooks/eka")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// fhirBundleFor is where your system builds an ABDM-compliant FHIR R4 bundle.
func fhirBundleFor(careContextID string) []byte {
	return []byte(`{"resourceType":"Bundle","id":"` + careContextID + `","type":"document"}`)
}

// matchPatients is where you query your own patient records.
func matchPatients(e *carecontext.DiscoverEvent) []carecontext.Patient {
	return []carecontext.Patient{{
		ReferenceNumber: "P1",
		Display:         e.PatientName,
		HIType:          carecontext.HITypeOPConsultation,
		CareContexts:    []carecontext.DiscoveredCareContext{{ReferenceNumber: "CC101", Display: "OP Consult"}},
	}}
}
