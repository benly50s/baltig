// internal/config/loader.go
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultRef = "main"

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config: cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".baltig"), nil
}

// ConfigPath returns the path to the config file.
func ConfigPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaults(), nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	applyDefaults(&cfg)
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	path := filepath.Join(dir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func Validate(cfg *Config) error {
	if cfg.Global.GitLabURL == "" {
		return errors.New("gitlab_url is required; run 'baltig onboard'")
	}
	if cfg.Global.Token == "" {
		return errors.New("token is required; run 'baltig onboard'")
	}
	return nil
}

func (cfg *Config) AddRecent(namespace string) {
	// Remove if already present, then prepend
	filtered := make([]RecentEntry, 0, len(cfg.Global.Recents))
	for _, r := range cfg.Global.Recents {
		if r.Namespace != namespace {
			filtered = append(filtered, r)
		}
	}
	cfg.Global.Recents = append([]RecentEntry{{Namespace: namespace, LastUsed: time.Now().UTC()}}, filtered...)
	if len(cfg.Global.Recents) > 10 {
		cfg.Global.Recents = cfg.Global.Recents[:10]
	}
}

func (cfg *Config) AddProject(p ProjectEntry) {
	for i, existing := range cfg.Projects {
		if existing.ID == p.ID {
			cfg.Projects[i] = p // update in place
			return
		}
	}
	cfg.Projects = append(cfg.Projects, p)
}

func (cfg *Config) RemoveProject(id int64) {
	result := make([]ProjectEntry, 0, len(cfg.Projects))
	for _, p := range cfg.Projects {
		if p.ID != id {
			result = append(result, p)
		}
	}
	cfg.Projects = result
}

func defaults() *Config {
	return &Config{
		Global: GlobalConfig{DefaultRef: defaultRef},
	}
}

func applyDefaults(cfg *Config) {
	if cfg.Global.DefaultRef == "" {
		cfg.Global.DefaultRef = defaultRef
	}
}
