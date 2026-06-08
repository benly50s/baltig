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
	Options     []string // if non-empty, variable has predefined choices
}

// GetCIVariables fetches the merged CI config (include: resolved) and returns
// variables that have a description — matching GitLab web "Run pipeline" form.
func (c *Client) GetCIVariables(projectID int64, ref string) ([]CIVariable, error) {
	// ProjectLint resolves all include: directives and returns merged YAML.
	lintOpts := &gl.ProjectLintOptions{
		Ref:         gl.Ptr(ref),
		DryRun:      gl.Ptr(false),
		IncludeJobs: gl.Ptr(false),
	}
	result, _, err := c.gl.Validate.ProjectLint(projectID, lintOpts)
	if err != nil || result.MergedYaml == "" {
		// Fallback: try raw file directly
		return c.getCIVariablesFromRaw(projectID, ref)
	}
	return parseCIVariables([]byte(result.MergedYaml))
}

func (c *Client) getCIVariablesFromRaw(projectID int64, ref string) ([]CIVariable, error) {
	opts := &gl.GetRawFileOptions{Ref: gl.Ptr(ref)}
	content, resp, err := c.gl.RepositoryFiles.GetRawFile(projectID, ".gitlab-ci.yml", opts)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return nil, fmt.Errorf(".gitlab-ci.yml 파일 없음 (ref: %s)", ref)
		}
		return nil, fmt.Errorf("fetch .gitlab-ci.yml: %w", err)
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
		// Only variables with a description are shown in the GitLab web "Run pipeline" form.
		// Scalar variables (no description) are CI-internal config and not user-settable.
		if node.Kind == yaml.MappingNode {
			var complex struct {
				Value       string   `yaml:"value"`
				Description string   `yaml:"description"`
				Options     []string `yaml:"options"`
			}
			if err := node.Decode(&complex); err == nil && complex.Description != "" {
				result = append(result, CIVariable{
					Key:         key,
					Value:       complex.Value,
					Description: complex.Description,
					Options:     complex.Options,
				})
			}
		}
	}
	return result, nil
}
