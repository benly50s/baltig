// internal/gitlab/projects.go
package gitlab

import (
	"fmt"

	gl "gitlab.com/gitlab-org/api/client-go"
)

// Project is a simplified project representation.
type Project struct {
	ID                int64
	Name              string
	NameWithNamespace string // "group/repo"
	WebURL            string
}

// SearchProjects searches for projects matching query (membership only, last activity order).
func (c *Client) SearchProjects(query string) ([]Project, error) {
	opts := &gl.ListProjectsOptions{
		Search:      gl.Ptr(query),
		Membership:  gl.Ptr(true),
		OrderBy:     gl.Ptr("last_activity_at"),
		ListOptions: gl.ListOptions{PerPage: 20},
	}
	projects, _, err := c.gl.Projects.ListProjects(opts)
	if err != nil {
		return nil, fmt.Errorf("search projects: %w", err)
	}
	result := make([]Project, len(projects))
	for i, p := range projects {
		result[i] = Project{
			ID:                p.ID,
			Name:              p.Name,
			NameWithNamespace: p.PathWithNamespace,
			WebURL:            p.WebURL,
		}
	}
	return result, nil
}

// GetProject fetches a project by ID.
func (c *Client) GetProject(id int64) (*Project, error) {
	p, _, err := c.gl.Projects.GetProject(id, nil)
	if err != nil {
		return nil, fmt.Errorf("get project %d: %w", id, err)
	}
	return &Project{
		ID:                p.ID,
		Name:              p.Name,
		NameWithNamespace: p.PathWithNamespace,
		WebURL:            p.WebURL,
	}, nil
}
