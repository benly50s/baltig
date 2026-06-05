// internal/gitlab/jobs_test.go
package gitlab_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/benly/baltig/internal/gitlab"
)

func TestListJobs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/5/pipelines/100/jobs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 200, "name": "build", "stage": "build", "status": "success", "web_url": "https://gl.example.com/-/jobs/200"},
			{"id": 201, "name": "deploy", "stage": "deploy", "status": "failed", "web_url": "https://gl.example.com/-/jobs/201"},
		})
	})
	srv := newTestServer(t, mux)

	client, _ := gitlab.New(srv.URL, "tok")
	jobs, err := client.ListJobs(5, 100)
	if err != nil {
		t.Fatalf("ListJobs() error = %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("len = %d, want 2", len(jobs))
	}
	if jobs[0].Name != "build" {
		t.Errorf("job[0].Name = %q, want build", jobs[0].Name)
	}
	if jobs[1].Status != "failed" {
		t.Errorf("job[1].Status = %q, want failed", jobs[1].Status)
	}
}

func TestGetJobLog(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/projects/5/jobs/200/trace", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Running build...\nDone.\n"))
	})
	srv := newTestServer(t, mux)

	client, _ := gitlab.New(srv.URL, "tok")
	log, err := client.GetJobLog(5, 200)
	if err != nil {
		t.Fatalf("GetJobLog() error = %v", err)
	}
	if log != "Running build...\nDone.\n" {
		t.Errorf("log = %q, want 'Running build...\\nDone.\\n'", log)
	}
}
