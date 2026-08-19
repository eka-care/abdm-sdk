package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eka-care/eka-sdk-go/auth"
	"github.com/eka-care/eka-sdk-go/internal/config"
	ekahttp "github.com/eka-care/eka-sdk-go/internal/http"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// TestRequestWithExpiredTokenDoesNotDeadlock covers the reentrancy bug where an
// expired token made every request hang forever: Do -> Config.Token ->
// provider.Retrieve (holding its lock) -> Service.RefreshToken -> Do ->
// Config.Token -> Retrieve -> lock held by this same goroutine.
func TestRequestWithExpiredTokenDoesNotDeadlock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/connect-auth/v1/account/refresh":
			_ = json.NewEncoder(w).Encode(auth.RefreshTokenResponse{
				AccessToken:      "fresh-access",
				RefreshToken:     "fresh-refresh",
				ExpiresIn:        3600,
				RefreshExpiresIn: 7200,
			})
		case "/abdm/v1/care-contexts/link":
			if got := r.Header.Get("Authorization"); got != "Bearer fresh-access" {
				t.Errorf("Authorization = %q, want the refreshed token", got)
			}
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Errorf("unexpected path %q", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	cfg := config.NewConfig()
	cfg.BaseURL = srv.URL
	svc := auth.NewService(cfg)

	// Access token already expired; refresh token still good.
	p := auth.NewStaticCredentialsProviderWithService(svc, "stale-access", "good-refresh", -1, 7200)
	cfg.SetTokenFunc(func(ctx context.Context) (string, error) {
		creds, err := p.Retrieve(ctx)
		if err != nil {
			return "", err
		}
		return creds.AccessToken, nil
	})

	client := ekahttp.NewClientFromInterface(cfg)

	done := make(chan error, 1)
	go func() {
		_, err := client.Do(context.Background(), &interfaces.HTTPRequest{
			Method: "POST",
			Path:   "/abdm/v1/care-contexts/link",
			Body:   map[string]string{"abha_address": "x@sbx"},
		})
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request did not complete within 5s: token refresh deadlocked")
	}
}
