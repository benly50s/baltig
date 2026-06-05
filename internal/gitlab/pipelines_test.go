// internal/gitlab/pipelines_test.go
package gitlab_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/benly/baltig/internal/gitlab"
)

func TestListPipelines(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/5/pipelines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 100, "status": "success", "ref": "main", "web_url": "https://gl.example.com/-/pipelines/100"},
		})
	})
	srv := newTestServer(t, mux)

	client, _ := gitlab.New(srv.URL, "tok")
	pipelines, err := client.ListPipelines(5)
	if err != nil {
		t.Fatalf("ListPipelines() error = %v", err)
	}
	if len(pipelines) != 1 {
		t.Fatalf("len = %d, want 1", len(pipelines))
	}
	if pipelines[0].Status != "success" {
		t.Errorf("status = %q, want success", pipelines[0].Status)
	}
	if pipelines[0].ID != 100 {
		t.Errorf("ID = %d, want 100", pipelines[0].ID)
	}
}

func TestCreatePipeline(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/5/pipeline", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": 101, "status": "pending", "ref": "main",
			"web_url": "https://gl.example.com/-/pipelines/101",
		})
	})
	srv := newTestServer(t, mux)

	client, _ := gitlab.New(srv.URL, "tok")
	p, err := client.CreatePipeline(5, "main", []gitlab.PipelineVariable{
		{Key: "ENV", Value: "staging"},
	})
	if err != nil {
		t.Fatalf("CreatePipeline() error = %v", err)
	}
	if p.ID != 101 {
		t.Errorf("pipeline ID = %d, want 101", p.ID)
	}
	if p.Status != "pending" {
		t.Errorf("status = %q, want pending", p.Status)
	}
}
