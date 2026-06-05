// internal/gitlab/client_test.go
package gitlab_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benly/baltig/internal/gitlab"
)

func newTestServer(t *testing.T, mux *http.ServeMux) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestPing(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"username": "testuser", "id": 1})
	})
	srv := newTestServer(t, mux)

	client, err := gitlab.New(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	username, err := client.Ping()
	if err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
	if username != "testuser" {
		t.Errorf("username = %q, want testuser", username)
	}
}

func TestPingInvalidToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/user", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"401 Unauthorized"}`, http.StatusUnauthorized)
	})
	srv := newTestServer(t, mux)

	client, err := gitlab.New(srv.URL, "bad-token")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Ping()
	if err == nil {
		t.Error("Ping() expected error for 401, got nil")
	}
}
