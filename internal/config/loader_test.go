// internal/config/loader_test.go
package config_test

import (
	"path/filepath"
	"testing"

	"github.com/benly/baltig/internal/config"
)

func TestLoadNotExist(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Global.DefaultRef != "main" {
		t.Errorf("default ref = %q, want %q", cfg.Global.DefaultRef, "main")
	}
}

func TestSaveAndLoad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := &config.Config{
		Global: config.GlobalConfig{
			GitLabURL:  "https://gitlab.example.com",
			Token:      "glpat-test",
			DefaultRef: "main",
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Global.GitLabURL != cfg.Global.GitLabURL {
		t.Errorf("GitLabURL = %q, want %q", loaded.Global.GitLabURL, cfg.Global.GitLabURL)
	}
	if loaded.Global.Token != cfg.Global.Token {
		t.Errorf("Token = %q, want %q", loaded.Global.Token, cfg.Global.Token)
	}

	// Verify project round-trip
	cfg.AddProject(config.ProjectEntry{ID: 42, Namespace: "mygroup/myrepo", Name: "myrepo", Starred: true})
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save() with project error = %v", err)
	}
	loaded2, err := config.Load()
	if err != nil {
		t.Fatalf("Load() after project save error = %v", err)
	}
	if len(loaded2.Projects) != 1 {
		t.Fatalf("len(projects) = %d, want 1", len(loaded2.Projects))
	}
	if loaded2.Projects[0].ID != 42 {
		t.Errorf("project ID = %d, want 42", loaded2.Projects[0].ID)
	}
	if !loaded2.Projects[0].Starred {
		t.Error("project Starred = false, want true")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "missing url",
			cfg:     &config.Config{Global: config.GlobalConfig{Token: "tok"}},
			wantErr: true,
		},
		{
			name:    "missing token",
			cfg:     &config.Config{Global: config.GlobalConfig{GitLabURL: "https://gl.example.com"}},
			wantErr: true,
		},
		{
			name: "valid",
			cfg: &config.Config{Global: config.GlobalConfig{
				GitLabURL: "https://gl.example.com",
				Token:     "tok",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddRecent(t *testing.T) {
	cfg := &config.Config{}
	cfg.AddRecent("group/repo-a")
	cfg.AddRecent("group/repo-b")
	cfg.AddRecent("group/repo-a") // duplicate → moves to front

	if cfg.Global.Recents[0].Namespace != "group/repo-a" {
		t.Errorf("first recent = %q, want group/repo-a", cfg.Global.Recents[0].Namespace)
	}
	if len(cfg.Global.Recents) != 2 {
		t.Errorf("len(recents) = %d, want 2", len(cfg.Global.Recents))
	}
}

func TestAddRecentCap(t *testing.T) {
	cfg := &config.Config{}
	for i := 0; i < 12; i++ {
		cfg.AddRecent(filepath.Join("group", string(rune('a'+i))))
	}
	if len(cfg.Global.Recents) != 10 {
		t.Errorf("len(recents) = %d, want 10", len(cfg.Global.Recents))
	}
}

func TestAddAndRemoveProject(t *testing.T) {
	cfg := &config.Config{}
	cfg.AddProject(config.ProjectEntry{ID: 1, Namespace: "g/a", Name: "a"})
	cfg.AddProject(config.ProjectEntry{ID: 1, Namespace: "g/a", Name: "a-updated", Starred: true}) // upsert
	if len(cfg.Projects) != 1 {
		t.Errorf("len(projects) = %d, want 1 (upsert should not add duplicate)", len(cfg.Projects))
	}
	if cfg.Projects[0].Name != "a-updated" {
		t.Errorf("project Name = %q, want a-updated (upsert should update)", cfg.Projects[0].Name)
	}
	cfg.RemoveProject(1)
	if len(cfg.Projects) != 0 {
		t.Errorf("len(projects) = %d after remove, want 0", len(cfg.Projects))
	}
}
