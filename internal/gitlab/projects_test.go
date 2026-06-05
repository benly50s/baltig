// internal/gitlab/projects_test.go
package gitlab_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/benly/baltig/internal/gitlab"
)

func TestSearchProjects(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":                  10,
				"name":                "frontend",
				"path_with_namespace": "mygroup/frontend",
				"web_url":             "https://gl.example.com/mygroup/frontend",
			},
		})
	})
	// /api/v4/user needed by some versions of the client
	mux.HandleFunc("/api/v4/user", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"username": "u", "id": 1})
	})
	srv := newTestServer(t, mux)

	client, _ := gitlab.New(srv.URL, "tok")
	projects, err := client.SearchProjects("front")
	if err != nil {
		t.Fatalf("SearchProjects() error = %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("len(projects) = %d, want 1", len(projects))
	}
	if projects[0].ID != 10 {
		t.Errorf("project ID = %d, want 10", projects[0].ID)
	}
	if projects[0].NameWithNamespace != "mygroup/frontend" {
		t.Errorf("namespace = %q, want mygroup/frontend", projects[0].NameWithNamespace)
	}
}
