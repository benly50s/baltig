// internal/gitlab/ci_vars.go
package gitlab

import (
	"fmt"

	gl "gitlab.com/gitlab-org/api/client-go"
	"gopkg.in/yaml.v3"
)

// CIVariable represents a variable from .gitlab-ci.yml.
type CIVariable struct {
	Key         string
	Value       string
	Description string
}

// GetCIVariables fetches .gitlab-ci.yml and returns the top-level variables section.
// Returns empty slice (no error) if the file doesn't exist or has no variables.
func (c *Client) GetCIVariables(projectID int64, ref string) ([]CIVariable, error) {
	opts := &gl.GetRawFileOptions{
		Ref: gl.Ptr(ref),
	}
	content, _, err := c.gl.RepositoryFiles.GetRawFile(int(projectID), ".gitlab-ci.yml", opts)
	if err != nil {
		// File not found is acceptable — return empty
		return nil, nil
	}

	return parseCIVariables(content)
}

// parseCIVariables parses the variables: section from .gitlab-ci.yml content.
func parseCIVariables(content []byte) ([]CIVariable, error) {
	var ciFile struct {
		Variables map[string]yaml.Node `yaml:"variables"`
	}
	if err := yaml.Unmarshal(content, &ciFile); err != nil {
		return nil, fmt.Errorf("parse .gitlab-ci.yml: %w", err)
	}

	var result []CIVariable
	for key, node := range ciFile.Variables {
		// Only include variables with a description — these are intended for
		// manual pipeline runs (matches GitLab web UI "Run pipeline" behavior).
		// Scalar variables (no description) are internal CI config values.
		if node.Kind == yaml.MappingNode {
			var complex struct {
				Value       string `yaml:"value"`
				Description string `yaml:"description"`
			}
			if err := node.Decode(&complex); err == nil && complex.Description != "" {
				result = append(result, CIVariable{
					Key:         key,
					Value:       complex.Value,
					Description: complex.Description,
				})
			}
		}
	}
	return result, nil
}
