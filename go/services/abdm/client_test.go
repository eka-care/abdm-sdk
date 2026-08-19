package abdm

import (
	"testing"
	"time"

	"github.com/eka-care/abdm-sdk/go/internal/config"
)

func TestClientExposesCareContexts(t *testing.T) {
	cfg := &config.Config{BaseURL: "https://example.invalid", Timeout: time.Second}
	if NewClient(cfg).CareContexts() == nil {
		t.Error("CareContexts() = nil, want a service")
	}
}
