// internal/gitlab/client.go
package gitlab

import (
	"fmt"

	gl "gitlab.com/gitlab-org/api/client-go"
)

// Client wraps the official GitLab API client.
type Client struct {
	gl  *gl.Client
	url string
}

// New creates a GitLab client from base URL and PAT.
func New(baseURL, token string) (*Client, error) {
	c, err := gl.NewClient(token, gl.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("create gitlab client: %w", err)
	}
	return &Client{gl: c, url: baseURL}, nil
}

// Ping verifies the connection and returns the authenticated username.
func (c *Client) Ping() (string, error) {
	user, _, err := c.gl.Users.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("ping gitlab: %w", err)
	}
	return user.Username, nil
}
