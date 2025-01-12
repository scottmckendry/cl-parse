package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func setupTestRepo(t *testing.T) (string, func()) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize a new git repository
	_, err = git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

func TestIsGitRepo(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (string, func())
		expected bool
	}{
		{
			name: "valid git repository",
			setup: func() (string, func()) {
				return setupTestRepo(t)
			},
			expected: true,
		},
		{
			name: "non-git directory",
			setup: func() (string, func()) {
				dir, err := os.MkdirTemp("", "non-git-*")
				if err != nil {
					t.Fatal(err)
				}
				return dir, func() { os.RemoveAll(dir) }
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := tt.setup()
			defer cleanup()

			result := IsGitRepo(dir)
			if result != tt.expected {
				t.Errorf("IsGitRepo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetCommmitBodyFromSha(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test repository with a commit
	repo, err := git.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	// Create a test file
	testFile := filepath.Join(dir, "test.txt")

	tests := []struct {
		name          string
		commitMsg     string
		fileContent   string
		expectedBody  string
		expectedError bool
	}{
		{
			name:          "commit with body",
			commitMsg:     "Initial commit\n\nThis is the body\nMultiple lines",
			fileContent:   "test content 1",
			expectedBody:  "This is the body\nMultiple lines",
			expectedError: false,
		},
		{
			name:          "commit without body",
			commitMsg:     "Single line commit",
			fileContent:   "test content 2",
			expectedBody:  "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write different content for each test
			err = os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			_, err = w.Add("test.txt")
			if err != nil {
				t.Fatal(err)
			}

			hash, err := w.Commit(tt.commitMsg, &git.CommitOptions{
				Author: &object.Signature{
					Name:  "test",
					Email: "test@example.com",
					When:  time.Now(),
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			body, err := GetCommmitBodyFromSha(dir, hash.String())
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if body != tt.expectedBody {
					t.Errorf("GetCommmitBodyFromSha() = %q, want %q", body, tt.expectedBody)
				}
			}
		})
	}
}
