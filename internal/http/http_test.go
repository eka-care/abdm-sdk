package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eka-care/eka-sdk-go/internal/config"
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
)

// A long-running client must pick up a refreshed token, not the one it was built with.
func TestClientResolvesTokenPerRequest(t *testing.T) {
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	var n int64
	cfg := &config.Config{BaseURL: srv.URL, Timeout: 5 * time.Second, UserAgent: "test"}
	cfg.SetTokenFunc(func(ctx context.Context) (string, error) {
		return fmt.Sprintf("token-%d", atomic.AddInt64(&n, 1)), nil
	})

	c := NewClientFromInterface(cfg)
	for i := 0; i < 2; i++ {
		if _, err := c.Do(context.Background(), &interfaces.HTTPRequest{Method: "GET", Path: "/x"}); err != nil {
			t.Fatalf("request %d: %v", i, err)
		}
	}

	want := []string{"Bearer token-1", "Bearer token-2"}
	if len(seen) != 2 || seen[0] != want[0] || seen[1] != want[1] {
		t.Errorf("Authorization headers = %q, want %q", seen, want)
	}
}

// With no token func set, the static config token is still used (M1 back-compat).
func TestClientFallsBackToConfigToken(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	cfg := &config.Config{BaseURL: srv.URL, Timeout: 5 * time.Second, UserAgent: "test"}
	cfg.SetAuthorizationToken("static-token")

	if _, err := NewClientFromInterface(cfg).Do(context.Background(),
		&interfaces.HTTPRequest{Method: "GET", Path: "/x"}); err != nil {
		t.Fatal(err)
	}
	if got != "Bearer static-token" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer static-token")
	}
}
