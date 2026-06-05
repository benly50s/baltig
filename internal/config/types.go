// internal/config/types.go
package config

import "time"

type Config struct {
	Global   GlobalConfig   `yaml:"global"`
	Projects []ProjectEntry `yaml:"projects"`
}

type GlobalConfig struct {
	GitLabURL  string        `yaml:"gitlab_url"`
	Token      string        `yaml:"token"`
	DefaultRef string        `yaml:"default_ref"`
	Recents    []RecentEntry `yaml:"recents,omitempty"`
}

type ProjectEntry struct {
	ID        int    `yaml:"id"`
	Namespace string `yaml:"namespace"` // stores "group/repo"
	Name      string `yaml:"name"`
	Starred   bool   `yaml:"starred"`
}

type RecentEntry struct {
	Namespace string    `yaml:"namespace"`
	LastUsed  time.Time `yaml:"last_used"`
}
