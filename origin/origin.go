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

type Issue struct {
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
}

type IssueProvider interface {
	GetIssue(issueNumber string) (*Issue, error)
}

type Config struct {
	URL   string
	Token string
}

func NewIssueProvider(config Config) (IssueProvider, error) {
	if strings.Contains(config.URL, "github.com") {
		return NewGitHubProvider(config), nil
	}
	if strings.Contains(config.URL, "dev.azure.com") {
		return NewAzureDevOpsProvider(config), nil
	}
	return nil, fmt.Errorf("unsupported git provider for URL: %s", config.URL)
}

type GitHubProvider struct {
	config Config
	owner  string
	repo   string
}

func NewGitHubProvider(config Config) *GitHubProvider {
	owner, repo := parseGitHubURL(config.URL)
	return &GitHubProvider{
		config: config,
		owner:  owner,
		repo:   repo,
	}
}

func (g *GitHubProvider) GetIssue(issueNumber string) (*Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", g.owner, g.repo, issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "go-changelog")
	req.Header.Set("Authorization", "Bearer "+g.config.Token)

	return g.doRequest(req)
}

// AzureDevOpsProvider implements IssueProvider for Azure DevOps
type AzureDevOpsProvider struct {
	config Config
	org    string
}

func NewAzureDevOpsProvider(config Config) *AzureDevOpsProvider {
	org := parseAzureDevOpsURL(config.URL)
	return &AzureDevOpsProvider{
		config: config,
		org:    org,
	}
}

func (a *AzureDevOpsProvider) GetIssue(issueNumber string) (*Issue, error) {
	url := fmt.Sprintf(
		"https://dev.azure.com/%s/_apis/wit/workitems/%s?api-version=7.1",
		a.org,
		issueNumber,
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if a.config.Token != "" {
		encodedPat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(":%s", a.config.Token)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set(
			"Authorization",
			fmt.Sprintf("Basic %s", encodedPat),
		)
	}

	return a.doRequest(req)
}

func (g *GitHubProvider) doRequest(req *http.Request) (*Issue, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get issue details: %s", resp.Status)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

func (a *AzureDevOpsProvider) doRequest(req *http.Request) (*Issue, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get issue details: %s", resp.Status)
	}

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

	// Clean up the description
	description := azureResponse.Fields.Description
	description = html.UnescapeString(description)
	description = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(description, "")
	description = strings.TrimSpace(description)

	return &Issue{
		Number: azureResponse.ID,
		Title:  azureResponse.Fields.Title,
		Body:   description,
	}, nil
}

func parseGitHubURL(url string) (owner, repo string) {
	parts := strings.Split(url, "/")
	return parts[len(parts)-2], parts[len(parts)-1]
}

func parseAzureDevOpsURL(url string) (org string) {
	parts := strings.Split(url, "/")
	return parts[len(parts)-3]
}
