// internal/gitlab/pipelines.go
package gitlab

import (
	"fmt"
	"time"

	gl "gitlab.com/gitlab-org/api/client-go"
)

// Pipeline represents a GitLab pipeline.
type Pipeline struct {
	ID        int64
	Status    string // "running", "success", "failed", "pending", "canceled"
	Ref       string // branch or tag
	WebURL    string
	CreatedAt *time.Time
}

// PipelineVariable is a key-value pair for pipeline creation.
type PipelineVariable struct {
	Key   string
	Value string
}

// ListPipelines returns the 20 most recent pipelines for a project.
func (c *Client) ListPipelines(projectID int64) ([]Pipeline, error) {
	opts := &gl.ListProjectPipelinesOptions{
		ListOptions: gl.ListOptions{PerPage: 20},
	}
	pipelines, _, err := c.gl.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return nil, fmt.Errorf("list pipelines for project %d: %w", projectID, err)
	}
	result := make([]Pipeline, len(pipelines))
	for i, p := range pipelines {
		result[i] = Pipeline{
			ID:        p.ID,
			Status:    p.Status,
			Ref:       p.Ref,
			WebURL:    p.WebURL,
			CreatedAt: p.CreatedAt,
		}
	}
	return result, nil
}

// CreatePipeline triggers a new pipeline on the given ref with optional variables.
func (c *Client) CreatePipeline(projectID int64, ref string, vars []PipelineVariable) (*Pipeline, error) {
	glVars := make([]*gl.PipelineVariableOptions, len(vars))
	for i, v := range vars {
		v := v // capture loop variable
		glVars[i] = &gl.PipelineVariableOptions{
			Key:   &v.Key,
			Value: &v.Value,
		}
	}
	opts := &gl.CreatePipelineOptions{
		Ref:       gl.Ptr(ref),
		Variables: &glVars,
	}
	p, _, err := c.gl.Pipelines.CreatePipeline(projectID, opts)
	if err != nil {
		return nil, fmt.Errorf("create pipeline for project %d on %s: %w", projectID, ref, err)
	}
	return &Pipeline{
		ID:        p.ID,
		Status:    p.Status,
		Ref:       p.Ref,
		WebURL:    p.WebURL,
		CreatedAt: p.CreatedAt,
	}, nil
}

// DeletePipeline deletes a pipeline by ID.
func (c *Client) DeletePipeline(projectID, pipelineID int64) error {
	_, err := c.gl.Pipelines.DeletePipeline(projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("delete pipeline %d: %w", pipelineID, err)
	}
	return nil
}

// GetPipeline fetches a single pipeline's current state.
func (c *Client) GetPipeline(projectID, pipelineID int64) (*Pipeline, error) {
	p, _, err := c.gl.Pipelines.GetPipeline(projectID, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("get pipeline %d: %w", pipelineID, err)
	}
	return &Pipeline{
		ID:        p.ID,
		Status:    p.Status,
		Ref:       p.Ref,
		WebURL:    p.WebURL,
		CreatedAt: p.CreatedAt,
	}, nil
}
