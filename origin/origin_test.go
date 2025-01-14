package origin

import (
	"net/http"
	"net/http/httptest"
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
			name:    "unsupported provider",
			url:     "https://gitlab.com/owner/repo",
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/issues/1" {
			t.Errorf("Expected path /repos/owner/repo/issues/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"number": 1, "title": "Test Issue", "body": "Test Body"}`))
	}))
	defer server.Close()

	provider := &GitHubProvider{
		BaseProvider: BaseProvider{
			client: server.Client(),
		},
		owner: "scottmckendry",
		repo:  "cl-parse",
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
