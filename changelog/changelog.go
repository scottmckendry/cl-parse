package changelog

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"cl-parse/git"
)

type ChangelogEntry struct {
	Version    string              `json:"version"`
	Date       time.Time           `json:"date"`
	CompareURL string              `json:"compareUrl"`
	Changes    map[string][]Change `json:"changes"`
}

type Change struct {
	Scope       string `json:"scope,omitempty"`
	Description string `json:"description"`
	PR          string `json:"pr,omitempty"`
	Commit      string `json:"commit,omitempty"`
	CommitBody  string `json:"commitBody,omitempty"`
}

type Parser struct {
	entries     []ChangelogEntry
	IncludeBody bool
}

// Create a new Parser
func NewParser() *Parser {
	return &Parser{
		entries: make([]ChangelogEntry, 0),
	}
}

func (p *Parser) GetLatest() (*ChangelogEntry, error) {
	if len(p.entries) == 0 {
		return nil, fmt.Errorf("no changelog entries found")
	}
	return &p.entries[0], nil
}

func (p *Parser) GetVersion(version string) (*ChangelogEntry, error) {
	for _, entry := range p.entries {
		if entry.Version == version {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf("version %s not found", version)
}

// Parse the changelog content and return a slice of ChangelogEntry
func (p *Parser) Parse(content string) ([]ChangelogEntry, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentEntry *ChangelogEntry
	var currentSection string

	versionRegex := regexp.MustCompile(
		`## (?:\[)?(?:v)?([\d.]+(?:-[a-zA-Z0-9]+(?:\.[0-9]+)?)?)\]?(?:\((.*?)\))? \((\d{4}-\d{2}-\d{2})\)`,
	)
	changeRegex := regexp.MustCompile(
		`\* (?:\*\*(.*?)\*\*: )?(.+?)\s*(?:\((?:#(\d+))?\)?)?\s*(?:\((.*?)\))?$`,
	)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || line == "# Changelog" {
			continue
		}

		if matches := versionRegex.FindStringSubmatch(line); matches != nil {
			if currentEntry != nil {
				p.entries = append(p.entries, *currentEntry)
			}

			date, _ := time.Parse("2006-01-02", matches[3])
			currentEntry = &ChangelogEntry{
				Version:    matches[1],
				Date:       date,
				CompareURL: matches[2],
				Changes:    make(map[string][]Change),
			}
			continue
		}

		if strings.HasPrefix(line, "### ") {
			currentSection = strings.TrimPrefix(line, "### ")
			continue
		}

		if strings.HasPrefix(line, "* ") && currentEntry != nil {
			matches := changeRegex.FindStringSubmatch(line)
			if matches != nil {
				change := Change{
					Scope:       matches[1],
					Description: matches[2],
				}

				if matches[3] != "" {
					change.PR = matches[3]
				}
				if matches[4] != "" {
					// extract the hash from the md link
					parts := strings.Split(matches[4], "/")
					change.Commit = parts[len(parts)-1]
					change.Commit = change.Commit[:len(change.Commit)-1] // remove the closing parenthesis

					if p.IncludeBody {
						var err error
						change.CommitBody, err = git.GetCommmitBodyFromSha(".", change.Commit)
						if err != nil {
							return nil, fmt.Errorf("failed to get commit message: %w", err)
						}
					}
				}

				if currentSection != "" {
					currentEntry.Changes[currentSection] = append(
						currentEntry.Changes[currentSection],
						change,
					)
				}
			}
		}
	}

	// add the last entry
	if currentEntry != nil {
		p.entries = append(p.entries, *currentEntry)
	}

	return p.entries, nil
}
