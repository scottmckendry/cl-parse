package changelog

import (
	"bufio"
	"regexp"
	"strings"
	"time"
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
}

type Parser struct {
	entries []ChangelogEntry
}

// Create a new Parser
func NewParser() *Parser {
	return &Parser{
		entries: make([]ChangelogEntry, 0),
	}
}

// Parse the changelog content and return a slice of ChangelogEntry
func (p *Parser) Parse(content string) ([]ChangelogEntry, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentEntry *ChangelogEntry
	var currentSection string

	versionRegex := regexp.MustCompile(
		`## \[(?:v)?([\d.]+(?:-[a-zA-Z0-9]+(?:\.[0-9]+)?)?)\]\((.*?)\) \((\d{4}-\d{2}-\d{2})\)`,
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
					change.Commit = matches[4]
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
