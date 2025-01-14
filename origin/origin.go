package origin

import (
	"fmt"
	"net/http"
	"strings"
)

// Issue represents a work item or issue from a Git provider.
type Issue struct {
	Number int    `json:"number"          yaml:"number"          toml:"number"`
	Title  string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Body   string `json:"body,omitempty"  yaml:"body,omitempty"  toml:"body,omitempty"`
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
	if strings.Contains(config.URL, "gitlab.com") {
		return NewGitLabProvider(config), nil
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
