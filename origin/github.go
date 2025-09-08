package origin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GitHubProvider implements IssueProvider for GitHub repositories.
type GitHubProvider struct {
	BaseProvider
	owner string // GitHub repository owner
	repo  string // GitHub repository name
}

// NewGitHubProvider creates a new GitHub provider with the given configuration.
func NewGitHubProvider(config Config) *GitHubProvider {
	owner, repo := parseGitHubURL(config.URL)
	return &GitHubProvider{
		BaseProvider: NewBaseProvider(config),
		owner:        owner,
		repo:         repo,
	}
}

// createRequest creates a GitHub API request with appropriate headers.
func (g *GitHubProvider) createRequest(issueNumber string) (*http.Request, error) {
	if len(issueNumber) > 0 && (issueNumber[0] == '#' || issueNumber[0] == '!') {
		issueNumber = issueNumber[1:]
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s",
		g.owner, g.repo, issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "go-changelog")

	if g.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+g.config.Token)
	}

	return req, nil
}

// GetIssue fetches issue details from GitHub.
func (g *GitHubProvider) GetIssue(issueNumber string) (*Issue, error) {
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

	var raw struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	issue := &Issue{Number: "#" + fmt.Sprintf("%d", raw.Number), Title: raw.Title, Body: raw.Body}
	return issue, nil
}

// parseGitHubURL extracts owner and repository name from a GitHub URL.
func parseGitHubURL(url string) (owner, repo string) {
	url = strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")

	if strings.HasPrefix(url, "git@github.com:") {
		parts := strings.Split(strings.TrimPrefix(url, "git@github.com:"), "/")
		if len(parts) >= 2 {
			return parts[0], parts[1]
		}
		return "", ""
	}

	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "github.com" && i+2 < len(parts) {
			return parts[i+1], parts[i+2]
		}
	}
	return "", ""
}
