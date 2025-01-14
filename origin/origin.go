package origin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
)

// Issue represents a work item or issue from a Git provider.
type Issue struct {
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
}

// IssueProvider defines the interface for fetching issue details from Git providers.
type IssueProvider interface {
	GetIssue(issueNumber string) (*Issue, error)
}

// Config contains the configuration needed to connect to a Git provider.
type Config struct {
	URL   string // Repository URL
	Token string // Authentication token
}

// NewIssueProvider creates an appropriate IssueProvider based on the repository URL.
// Currently supports GitHub and Azure DevOps.
func NewIssueProvider(config Config) (IssueProvider, error) {
	if strings.Contains(config.URL, "github.com") {
		return NewGitHubProvider(config), nil
	}
	if strings.Contains(config.URL, "dev.azure.com") {
		return NewAzureDevOpsProvider(config), nil
	}
	return nil, fmt.Errorf("unsupported git provider for URL: %s", config.URL)
}

// BaseProvider implements common functionality for all Git providers.
type BaseProvider struct {
	config Config
	client *http.Client
}

// NewBaseProvider creates a new BaseProvider with the given configuration.
func NewBaseProvider(config Config) BaseProvider {
	return BaseProvider{
		config: config,
		client: &http.Client{},
	}
}

// doRequest performs an HTTP request and handles common response scenarios.
// Returns nil response for 404 status and error for other non-200 statuses.
func (b *BaseProvider) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get issue details: %s", resp.Status)
	}

	return resp, nil
}

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

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

// AzureDevOpsProvider implements IssueProvider for Azure DevOps repositories.
type AzureDevOpsProvider struct {
	BaseProvider
	org string // Azure DevOps organization
}

// NewAzureDevOpsProvider creates a new Azure DevOps provider with the given configuration.
func NewAzureDevOpsProvider(config Config) *AzureDevOpsProvider {
	org := parseAzureDevOpsURL(config.URL)
	return &AzureDevOpsProvider{
		BaseProvider: NewBaseProvider(config),
		org:          org,
	}
}

// createRequest creates an Azure DevOps API request with appropriate headers.
func (a *AzureDevOpsProvider) createRequest(issueNumber string) (*http.Request, error) {
	url := fmt.Sprintf(
		"https://dev.azure.com/%s/_apis/wit/workitems/%s?api-version=7.1",
		a.org,
		issueNumber,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	encodedPat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(":%s", a.config.Token)))
	if a.config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedPat))
	}

	return req, nil
}

// GetIssue fetches work item details from Azure DevOps.
func (a *AzureDevOpsProvider) GetIssue(issueNumber string) (*Issue, error) {
	req, err := a.createRequest(issueNumber)
	if err != nil {
		return nil, err
	}

	resp, err := a.doRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	var azureResponse struct {
		ID     int `json:"id"`
		Fields struct {
			Title       string `json:"System.Title"`
			Description string `json:"System.Description"`
		} `json:"fields"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&azureResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Issue{
		Number: azureResponse.ID,
		Title:  azureResponse.Fields.Title,
		Body:   cleanDescription(azureResponse.Fields.Description),
	}, nil
}

// cleanDescription removes HTML tags and whitespace from the issue description.
func cleanDescription(description string) string {
	description = html.UnescapeString(description)
	description = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(description, "")
	return strings.TrimSpace(description)
}

// parseGitHubURL extracts owner and repository name from a GitHub URL.
func parseGitHubURL(url string) (owner, repo string) {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2], parts[len(parts)-1]
}

// parseAzureDevOpsURL extracts organization name from an Azure DevOps URL.
func parseAzureDevOpsURL(url string) (org string) {
	parts := strings.Split(url, "/")
	return parts[len(parts)-3]
}
