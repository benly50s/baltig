// internal/config/loader.go
package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultRef = "main"

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".baltig")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func Load() (*Config, error) {
	path := ConfigPath()
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
	if err := os.MkdirAll(ConfigDir(), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
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
	filtered := cfg.Global.Recents[:0]
	for _, r := range cfg.Global.Recents {
		if r.Namespace != namespace {
			filtered = append(filtered, r)
		}
	}
	cfg.Global.Recents = append([]RecentEntry{{Namespace: namespace}}, filtered...)
	if len(cfg.Global.Recents) > 10 {
		cfg.Global.Recents = cfg.Global.Recents[:10]
	}
}

func (cfg *Config) AddProject(p ProjectEntry) {
	for _, existing := range cfg.Projects {
		if existing.ID == p.ID {
			return
		}
	}
	cfg.Projects = append(cfg.Projects, p)
}

func (cfg *Config) RemoveProject(id int) {
	result := cfg.Projects[:0]
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
