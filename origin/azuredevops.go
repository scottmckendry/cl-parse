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
	isPullRequest := len(issueNumber) > 0 && issueNumber[0] == '!'
	if len(issueNumber) > 0 && (issueNumber[0] == '#' || issueNumber[0] == '!') {
		issueNumber = issueNumber[1:]
	}
	url := fmt.Sprintf(
		"https://dev.azure.com/%s/_apis/wit/workitems/%s?api-version=7.1",
		a.org,
		issueNumber,
	)

	if isPullRequest {
		url = fmt.Sprintf(
			"https://dev.azure.com/%s/_apis/git/pullrequests/%s?api-version=7.1",
			a.org,
			issueNumber,
		)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	encodedPat := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, ":%s", a.config.Token))
	if a.config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedPat))
	}

	return req, nil
}

// GetIssue fetches work item or pull request details from Azure DevOps.
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

	var azurePrResponse struct {
		ID          int    `json:"pullRequestId"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	isPullRequest := len(issueNumber) > 0 && issueNumber[0] == '!'
	if isPullRequest {
		if err := json.NewDecoder(resp.Body).Decode(&azurePrResponse); err != nil {
			return nil, fmt.Errorf("failed to decode pull request response: %w", err)
		}
		return &Issue{
			Number: "!" + fmt.Sprintf("%d", azurePrResponse.ID),
			Title:  azurePrResponse.Title,
			Body:   cleanDescription(azurePrResponse.Description),
		}, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&azureResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &Issue{
		Number: "#" + fmt.Sprintf("%d", azureResponse.ID),
		Title:  azureResponse.Fields.Title,
		Body:   cleanDescription(azureResponse.Fields.Description),
	}, nil
}

// parseAzureDevOpsURL extracts organization name from an Azure DevOps URL.
func parseAzureDevOpsURL(url string) (org string) {
	url = strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")

	if strings.HasPrefix(url, "git@ssh.dev.azure.com:") {
		parts := strings.Split(strings.TrimPrefix(url, "git@ssh.dev.azure.com:v3/"), "/")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "dev.azure.com" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// cleanDescription removes HTML tags and whitespace from the issue description.
func cleanDescription(description string) string {
	description = html.UnescapeString(description)
	description = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(description, "")
	return strings.TrimSpace(description)
}
