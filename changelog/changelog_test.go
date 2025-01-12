package changelog

import (
	"reflect"
	"testing"
	"time"
)

// Test helpers
func createTestEntry(version, date string, compareURL string, changes map[string][]Change) ChangelogEntry {
	return ChangelogEntry{
		Version:    version,
		Date:       mustParseTime(date),
		CompareURL: compareURL,
		Changes:    changes,
	}
}

func createTestChange(description string, scope string, pr string, commit string) Change {
	return Change{
		Description: description,
		Scope:       scope,
		PR:          pr,
		Commit:      commit,
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ChangelogEntry
		wantErr bool
	}{
		{
			name: "parses basic changelog",
			input: `# Changelog

## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)

### Features

* **api**: add new endpoint (#123)
* basic feature (1a196c09283903991da080552e3aa980ac64fec9)

### Bug Fixes

* **ui**: fix button alignment
`,
			want: []ChangelogEntry{
				createTestEntry("1.0.0", "2025-01-01", "https://github.com/user/repo/compare/v0.1.0...v1.0.0", map[string][]Change{
					"Features": {
						createTestChange("add new endpoint", "api", "123", ""),
						createTestChange("basic feature", "", "", "1a196c09283903991da080552e3aa980ac64fec9"),
					},
					"Bug Fixes": {
						createTestChange("fix button alignment", "ui", "", ""),
					},
				}),
			},
		},
		{
			name: "parses version with prerelease",
			input: `# Changelog
## [v1.0.0-alpha.1](https://github.com/user/repo/compare/v0.1.0...v1.0.0-alpha.1) (2025-01-01)

### Features

* basic feature
`,
			want: []ChangelogEntry{
				createTestEntry("1.0.0-alpha.1", "2025-01-01", "https://github.com/user/repo/compare/v0.1.0...v1.0.0-alpha.1", map[string][]Change{
					"Features": {
						createTestChange("basic feature", "", "", ""),
					},
				}),
			},
		},
		{
			name: "parses version without URL",
			input: `# Changelog
## 1.0.0 (2025-01-01)

### Features

* basic feature
`,
			want: []ChangelogEntry{
				createTestEntry("1.0.0", "2025-01-01", "", map[string][]Change{
					"Features": {
						createTestChange("basic feature", "", "", ""),
					},
				}),
			},
		},
		{
			name: "extracts hashes from commit links",
			input: `# Changelog
## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)

### Features

* basic feature ([commit](https://github.com/user/repo/commit/8f5b75c6ba6c525e29463e2a96fec119e426e283))
* another feature ([link text](https://github.com/user/repo/commit/22822a9f19442b51d952b550e73ad3c229583371))
* some docs ([docs link text](https://example.com/docs))
`,
			want: []ChangelogEntry{
				createTestEntry("1.0.0", "2025-01-01", "https://github.com/user/repo/compare/v0.1.0...v1.0.0", map[string][]Change{
					"Features": {
						createTestChange("basic feature", "", "", "8f5b75c6ba6c525e29463e2a96fec119e426e283"),
						createTestChange("another feature", "", "", "22822a9f19442b51d952b550e73ad3c229583371"),
						createTestChange("some docs", "", "", ""),
					},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}

func TestGetLatest(t *testing.T) {
	const testChangelog = `# Changelog

## [v2.0.0](https://github.com/user/repo/compare/v1.0.0...v2.0.0) (2025-02-01)
### Features
* new feature

## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)
### Features
* basic feature`

	t.Run("returns latest version", func(t *testing.T) {
		p := NewParser()
		_, err := p.Parse(testChangelog)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		latest, err := p.GetLatest()
		if err != nil {
			t.Fatalf("GetLatest failed: %v", err)
		}

		if latest.Version != "2.0.0" {
			t.Errorf("Expected version 2.0.0, got %s", latest.Version)
		}
	})
}

func TestGetVersion(t *testing.T) {
	const testChangelog = `# Changelog

## [v2.0.0](https://github.com/user/repo/compare/v1.0.0...v2.0.0) (2025-02-01)
### Features
* new feature

## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)
### Features
* basic feature`

	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "existing version",
			version: "1.0.0",
			want:    "1.0.0",
			wantErr: false,
		},
		{
			name:    "non-existent version",
			version: "3.0.0",
			want:    "",
			wantErr: true,
		},
	}

	p := NewParser()
	if _, err := p.Parse(testChangelog); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := p.GetVersion(tt.version)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if entry.Version != tt.want {
				t.Errorf("Expected version %s, got %s", tt.want, entry.Version)
			}
		})
	}
}

func TestEmptyChangelog(t *testing.T) {
	t.Run("handles empty changelog", func(t *testing.T) {
		p := NewParser()
		_, err := p.Parse("# Changelog\n")
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if _, err := p.GetLatest(); err == nil {
			t.Error("Expected error for empty changelog but got none")
		}

		if _, err := p.GetVersion("1.0.0"); err == nil {
			t.Error("Expected error for non-existent version but got none")
		}
	})
}

func mustParseTime(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t
}
