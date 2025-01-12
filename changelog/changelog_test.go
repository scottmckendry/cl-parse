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

func mustParseTime(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t
}
