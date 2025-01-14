package origin

import (
	"net/http"
	"testing"
)

func TestNewIssueProvider(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name: "github provider",
			url:  "https://github.com/owner/repo",
			want: "*origin.GitHubProvider",
		},
		{
			name: "azure devops provider",
			url:  "https://dev.azure.com/org/project/repo",
			want: "*origin.AzureDevOpsProvider",
		},
		{
			name: "gitlab provider",
			url:  "https://gitlab.com/owner/repo",
			want: "*origin.GitLabProvider",
		},
		{
			name:    "unsupported provider",
			url:     "https://bitbucket.org/owner/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{URL: tt.url}
			got, err := NewIssueProvider(config)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewIssueProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Error("NewIssueProvider() returned nil provider")
			}
		})
	}
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		url       string
		wantOwner string
		wantRepo  string
	}{
		{
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			url:       "https://github.com/owner/repo/",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			owner, repo := parseGitHubURL(tt.url)
			if owner != tt.wantOwner {
				t.Errorf("parseGitHubURL() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("parseGitHubURL() repo = %v, want %v", repo, tt.wantRepo)
			}
		})
	}
}

func TestGitHubProvider_GetIssue(t *testing.T) {
	provider := &GitHubProvider{
		BaseProvider: BaseProvider{
			client: &http.Client{},
		},
		owner: "scottmckendry",
		repo:  "cl-parse",
	}

	issue, err := provider.GetIssue("9")
	if err != nil {
		t.Fatalf("GetIssue() error = %v", err)
	}

	if issue.Number != 9 || issue.Title != "Test Issue" || issue.Body != "Test Body" {
		t.Errorf(
			"GetIssue() = %+v, want {Number: 9, Title: 'Test Issue', Body: 'Test Body'}",
			issue,
		)
	}
}

func TestParseGitLabURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantProject string
	}{
		{
			name:        "simple project path",
			url:         "https://gitlab.com/owner/repo",
			wantProject: "owner/repo",
		},
		{
			name:        "nested group project path",
			url:         "https://gitlab.com/group/subgroup/repo",
			wantProject: "group/subgroup/repo",
		},
		{
			name:        "with .git suffix",
			url:         "https://gitlab.com/owner/repo.git",
			wantProject: "owner/repo",
		},
		{
			name:        "with trailing slash",
			url:         "https://gitlab.com/owner/repo/",
			wantProject: "owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGitLabURL(tt.url)
			if got != tt.wantProject {
				t.Errorf("parseGitLabURL() = %v, want %v", got, tt.wantProject)
			}
		})
	}
}

func TestGitLabProvider_GetIssue(t *testing.T) {
	provider := &GitLabProvider{
		BaseProvider: BaseProvider{
			client: &http.Client{},
		},
		project: "scottmckendry/test",
	}

	issue, err := provider.GetIssue("1")
	if err != nil {
		t.Fatalf("GetIssue() error = %v", err)
	}

	if issue.Number != 1 || issue.Title != "Test Issue" || issue.Body != "Test Body" {
		t.Errorf(
			"GetIssue() = %+v, want {Number: 1, Title: 'Test Issue', Body: 'Test Body'}",
			issue,
		)
	}
}

func TestParseAzureDevOpsOrganization(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantOrg string
	}{
		{
			name:    "simple project path",
			url:     "https://dev.azure.com/org/project/repo",
			wantOrg: "org",
		},
		{
			name:    "nested project path",
			url:     "https://dev.azure.com/org/group/project/repo",
			wantOrg: "org",
		},
		{
			name:    "with trailing slash",
			url:     "https://dev.azure.com/org/project/repo/",
			wantOrg: "org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAzureDevOpsURL(tt.url)
			if got != tt.wantOrg {
				t.Errorf("parseAzureDevOpsOrganization() = %v, want %v", got, tt.wantOrg)
			}
		})
	}
}
