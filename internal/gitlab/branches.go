// internal/gitlab/branches.go
package gitlab

import (
	"fmt"

	gl "gitlab.com/gitlab-org/api/client-go"
)

// ListBranches returns branch names for a project, sorted by last commit date.
func (c *Client) ListBranches(projectID int64) ([]string, error) {
	opts := &gl.ListBranchesOptions{
		ListOptions: gl.ListOptions{PerPage: 50},
	}
	branches, _, err := c.gl.Branches.ListBranches(int(projectID), opts)
	if err != nil {
		return nil, fmt.Errorf("list branches for project %d: %w", projectID, err)
	}
	names := make([]string, len(branches))
	for i, b := range branches {
		names[i] = b.Name
	}
	return names, nil
}
