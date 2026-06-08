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

	var withDesc, withoutDesc []CIVariable
	for key, node := range ciFile.Variables {
		switch node.Kind {
		case yaml.ScalarNode:
			// Simple: ENV: "staging" — pre-fill value, no description
			withoutDesc = append(withoutDesc, CIVariable{Key: key, Value: node.Value})
		case yaml.MappingNode:
			var complex struct {
				Value       string `yaml:"value"`
				Description string `yaml:"description"`
			}
			if err := node.Decode(&complex); err == nil {
				cv := CIVariable{Key: key, Value: complex.Value, Description: complex.Description}
				if complex.Description != "" {
					withDesc = append(withDesc, cv)
				} else {
					withoutDesc = append(withoutDesc, cv)
				}
			}
		}
	}
	// description 있는 변수 먼저, 없는 변수 뒤에
	return append(withDesc, withoutDesc...), nil
}
