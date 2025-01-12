package changelog

import (
	"reflect"
	"testing"
	"time"
)

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
* basic feature (abc123)

### Bug Fixes

* **ui**: fix button alignment
`,
			want: []ChangelogEntry{
				{
					Version:    "1.0.0",
					Date:       mustParseTime("2025-01-01"),
					CompareURL: "https://github.com/user/repo/compare/v0.1.0...v1.0.0",
					Changes: map[string][]Change{
						"Features": {
							{
								Scope:       "api",
								Description: "add new endpoint",
								PR:          "123",
							},
							{
								Description: "basic feature",
								Commit:      "abc123",
							},
						},
						"Bug Fixes": {
							{
								Scope:       "ui",
								Description: "fix button alignment",
							},
						},
					},
				},
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
				{
					Version:    "1.0.0-alpha.1",
					Date:       mustParseTime("2025-01-01"),
					CompareURL: "https://github.com/user/repo/compare/v0.1.0...v1.0.0-alpha.1",
					Changes: map[string][]Change{
						"Features": {
							{
								Description: "basic feature",
							},
						},
					},
				},
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
				{
					Version:    "1.0.0",
					Date:       mustParseTime("2025-01-01"),
					CompareURL: "",
					Changes: map[string][]Change{
						"Features": {
							{
								Description: "basic feature",
							},
						},
					},
				},
			},
		},
		{
			name: "extracts hashes from commit links",
			input: `# Changelog
## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)

### Features

* basic feature ([commit](https://github.com/user/repo/commit/abc123))
* another feature ([link text](https://github.com/user/repo/commit/def456))
`,
			want: []ChangelogEntry{
				{
					Version:    "1.0.0",
					Date:       mustParseTime("2025-01-01"),
					CompareURL: "https://github.com/user/repo/compare/v0.1.0...v1.0.0",
					Changes: map[string][]Change{
						"Features": {
							{
								Description: "basic feature",
								Commit:      "abc123",
							},
							{
								Description: "another feature",
								Commit:      "def456",
							},
						},
					},
				},
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
	input := `# Changelog

## [v2.0.0](https://github.com/user/repo/compare/v1.0.0...v2.0.0) (2025-02-01)
### Features
* new feature

## [v1.0.0](https://github.com/user/repo/compare/v0.1.0...v1.0.0) (2025-01-01)
### Features
* basic feature`

	p := NewParser()
	_, err := p.Parse(input)
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
}

func TestGetVersion(t *testing.T) {
	input := `# Changelog

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
	_, err := p.Parse(input)
	if err != nil {
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
	input := "# Changelog\n"
	p := NewParser()
	_, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = p.GetLatest()
	if err == nil {
		t.Error("Expected error for empty changelog but got none")
	}

	_, err = p.GetVersion("1.0.0")
	if err == nil {
		t.Error("Expected error for non-existent version but got none")
	}
}

func mustParseTime(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t
}
