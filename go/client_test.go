package ekasdk

import (
	"context"
	"testing"
	"time"

	"github.com/eka-care/abdm-sdk/go/auth"
	"github.com/eka-care/abdm-sdk/go/internal/interfaces"
)

type stubProvider struct{ token string }

func (p stubProvider) Retrieve(context.Context) (*auth.Credentials, error) {
	return &auth.Credentials{AccessToken: p.token, ExpiresAt: time.Now().Add(time.Hour)}, nil
}

// WithCredentialsProvider must actually be consulted for the bearer token;
// it used to be stored and never used, so requests went out unauthenticated.
func TestWithCredentialsProviderResolvesToken(t *testing.T) {
	c := New(WithCredentialsProvider(stubProvider{token: "from-provider"}))

	tp, ok := c.config.(interfaces.TokenProvider)
	if !ok {
		t.Fatal("config does not implement TokenProvider")
	}
	got, err := tp.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-provider" {
		t.Errorf("Token() = %q, want %q", got, "from-provider")
	}
}
