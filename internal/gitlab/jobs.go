// internal/gitlab/jobs.go
package gitlab

import (
	"fmt"
	"time"
)

// Job represents a pipeline job.
type Job struct {
	ID        int64
	Name      string
	Stage     string
	Status    string
	WebURL    string
	CreatedAt *time.Time
}

// ListJobs returns all jobs for a pipeline.
func (c *Client) ListJobs(projectID int64, pipelineID int64) ([]Job, error) {
	jobs, _, err := c.gl.Jobs.ListPipelineJobs(projectID, pipelineID, nil)
	if err != nil {
		return nil, fmt.Errorf("list jobs for pipeline %d: %w", pipelineID, err)
	}
	result := make([]Job, len(jobs))
	for i, j := range jobs {
		result[i] = Job{
			ID:        j.ID,
			Name:      j.Name,
			Stage:     j.Stage,
			Status:    j.Status,
			WebURL:    j.WebURL,
			CreatedAt: j.CreatedAt,
		}
	}
	return result, nil
}

// GetJobLog fetches the log for a job and returns it as a string.
// For running jobs, callers should poll periodically.
func (c *Client) GetJobLog(projectID int64, jobID int64) (string, error) {
	reader, _, err := c.gl.Jobs.GetTraceFile(projectID, jobID)
	if err != nil {
		return "", fmt.Errorf("get job log %d: %w", jobID, err)
	}
	// bytes.Reader implements io.Reader; read all bytes
	buf := make([]byte, reader.Len())
	_, err = reader.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("read job log %d: %w", jobID, err)
	}
	return string(buf), nil
}
