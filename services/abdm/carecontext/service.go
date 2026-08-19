package carecontext

import (
	"context"

	"github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// Service calls the ABDM care-context APIs.
type Service struct {
	config interfaces.Config
	http   *http.Client
}

// NewService creates a care-context service.
func NewService(config interfaces.Config) *Service {
	return &Service{config: config, http: http.NewClientFromInterface(config)}
}

// Link links care contexts to a patient's ABHA address.
//
// The API is asynchronous: a 202 means accepted, not linked. The outcome
// arrives later as an abha.link_care_context webhook, which ParseWebhook
// decodes into a *LinkStatusEvent.
func (s *Service) Link(ctx context.Context, headers interfaces.Headers, req *LinkRequest) error {
	body := *req
	if body.OID == "" {
		body.OID = headers.PatientID
	}
	if body.PartnerUserID == "" {
		body.PartnerUserID = headers.PartnerUserID
	}
	_, err := s.http.Do(ctx, &interfaces.HTTPRequest{
		Method:  "POST",
		Path:    "/abdm/v1/care-contexts/link",
		Headers: headers,
		Body:    &body,
	})
	return err
}
