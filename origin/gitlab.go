package origin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GitLabProvider implements IssueProvider for GitLab repositories
type GitLabProvider struct {
	BaseProvider
	project string // URL-encoded project path with namespace (e.g., "group/project")
}

// NewGitLabProvider creates a new GitLab provider with the given configuration
func NewGitLabProvider(config Config) *GitLabProvider {
	project := parseGitLabURL(config.URL)
	return &GitLabProvider{
		BaseProvider: NewBaseProvider(config),
		project:      project,
	}
}

// createRequest creates a GitLab API request with appropriate headers
func (g *GitLabProvider) createRequest(issueNumber string) (*http.Request, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/issues/%s",
		url.PathEscape(g.project), issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if g.config.Token != "" {
		req.Header.Set("PRIVATE-TOKEN", g.config.Token)
	}

	return req, nil
}

// GetIssue fetches issue details from GitLab
func (g *GitLabProvider) GetIssue(issueNumber string) (*Issue, error) {
	type GitLabIssue struct {
		IID         int    `json:"iid"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	req, err := g.createRequest(issueNumber)
	if err != nil {
		return nil, err
	}

	resp, err := g.doRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	var gitlabIssue GitLabIssue
	if err := json.NewDecoder(resp.Body).Decode(&gitlabIssue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Issue{
		Number: gitlabIssue.IID,
		Title:  gitlabIssue.Title,
		Body:   gitlabIssue.Description,
	}, nil
}

// parseGitLabURL extracts project path from a GitLab URL
func parseGitLabURL(url string) string {
	url = strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "gitlab.com" && i+1 < len(parts) {
			return strings.Join(parts[i+1:], "/")
		}
	}
	return ""
}
