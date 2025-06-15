package changelog

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cl-parse/git"
	"cl-parse/origin"
)

const (
	dateFormat     = "2006-01-02"
	versionPattern = `## (?:\[)?(?:v)?([\d.]+(?:-[a-zA-Z0-9]+(?:\.[0-9]+)?)?)\]?(?:\((.*?)\))? \((\d{4}-\d{2}-\d{2})\)`
	changePattern  = `\* (?:\*\*(.*?)\*\*: )?(.+?)\s*(?:\((.*?)\))?(?:,\s*closes.*)?$`
)

type ChangelogEntry struct {
	Version    string              `json:"version"    yaml:"version"    toml:"version"`
	Date       time.Time           `json:"date"       yaml:"date"       toml:"date"`
	CompareURL string              `json:"compareUrl" yaml:"compareUrl" toml:"compareUrl"`
	Changes    map[string][]Change `json:"changes"    yaml:"changes"    toml:"changes"`
}

type Change struct {
	Scope        string          `json:"scope,omitempty"        yaml:"scope,omitempty"        toml:"scope,omitempty"`
	Description  string          `json:"description"            yaml:"description"            toml:"description"`
	Commit       string          `json:"commit,omitempty"       yaml:"commit,omitempty"       toml:"commit,omitempty"`
	CommitBody   string          `json:"commitBody,omitempty"   yaml:"commitBody,omitempty"   toml:"commitBody,omitempty"`
	RelatedItems []*origin.Issue `json:"relatedItems,omitempty" yaml:"relatedItems,omitempty" toml:"relatedItems,omitempty"`
}

type Parser struct {
	entries          []ChangelogEntry
	originUrl        string
	OriginToken      string
	IncludeBody      bool
	FetchItemDetails bool
}

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

func (p *Parser) Parse(content string) ([]ChangelogEntry, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentEntry *ChangelogEntry
	var currentSection string
	var err error

	versionRegex := regexp.MustCompile(versionPattern)
	changeRegex := regexp.MustCompile(changePattern)

	if p.FetchItemDetails {
		p.originUrl, err = git.GetOriginURL(".")
		if err != nil {
			return nil, fmt.Errorf("failed to get origin URL: %w", err)
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line == "# Changelog" {
			continue
		}

		if matches := versionRegex.FindStringSubmatch(line); matches != nil {
			if currentEntry != nil {
				p.entries = append(p.entries, *currentEntry)
			}

			currentEntry, err = p.createNewEntry(matches)
			if err != nil {
				return nil, err
			}
			continue
		}

		if strings.HasPrefix(line, "### ") {
			currentSection = strings.TrimPrefix(line, "### ")
			continue
		}

		if strings.HasPrefix(line, "* ") && currentEntry != nil {
			if err := p.parseChange(line, changeRegex, currentSection, currentEntry); err != nil {
				return nil, err
			}
		}
	}

	if currentEntry != nil {
		p.entries = append(p.entries, *currentEntry)
	}

	return p.entries, nil
}

func (p *Parser) createNewEntry(matches []string) (*ChangelogEntry, error) {
	date, err := time.Parse(dateFormat, matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return &ChangelogEntry{
		Version:    matches[1],
		Date:       date,
		CompareURL: matches[2],
		Changes:    make(map[string][]Change),
	}, nil
}

func (p *Parser) parseChange(
	line string,
	changeRegex *regexp.Regexp,
	currentSection string,
	currentEntry *ChangelogEntry,
) error {
	matches := changeRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	relatedItems, err := extractRelatedItems(matches[2], p.originUrl, p.OriginToken)
	if err != nil {
		return err
	}

	change := Change{
		Scope:        matches[1],
		Description:  matches[2],
		RelatedItems: relatedItems,
	}

	if matches[3] != "" {
		change.Commit = parseCommitHashFromLink(matches[3])
		if err := p.addCommitBody(&change); err != nil {
			return err
		}
		if change.CommitBody != "" {
			bodyItems, err := extractRelatedItems(change.CommitBody, p.originUrl, p.OriginToken)
			if err != nil {
				return err
			}
			for _, item := range bodyItems {
				if !containsIssue(change.RelatedItems, item) {
					change.RelatedItems = append(change.RelatedItems, item)
				}
			}
		}
	}

	if currentSection != "" {
		currentEntry.Changes[currentSection] = append(
			currentEntry.Changes[currentSection],
			change,
		)
	}
	return nil
}

func (p *Parser) addCommitBody(change *Change) error {
	if !p.IncludeBody || change.Commit == "" {
		return nil
	}

	body, err := git.GetCommmitBodyFromSha(".", change.Commit)
	if err != nil {
		return fmt.Errorf("failed to get commit message: %w", err)
	}
	change.CommitBody = body
	return nil
}

func parseCommitHashFromLink(link string) string {
	parts := strings.Split(link, "/")
	possibleHash := parts[len(parts)-1]

	if possibleHash[len(possibleHash)-1] == ')' {
		possibleHash = possibleHash[:len(possibleHash)-1]
	}

	if git.IsValidSha(possibleHash) {
		return possibleHash
	}

	return ""
}

func extractRelatedItems(text string, repoUrl string, token string) ([]*origin.Issue, error) {
	issueRegex := regexp.MustCompile(`[#!](\d+)`) // Updated regex to match both # and !
	matches := issueRegex.FindAllStringSubmatch(text, -1)

	seen := make(map[string]bool)
	var items []*origin.Issue

	if repoUrl != "" {
		provider, err := origin.NewIssueProvider(origin.Config{URL: repoUrl, Token: token})
		if err != nil {
			return nil, fmt.Errorf("failed to create issue provider: %w", err)
		}

		for _, match := range matches {
			if !seen[match[1]] {
				isPR := strings.HasPrefix(match[0], "!")

				issue, err := provider.GetIssue(match[1], isPR)
				if err != nil {
					return nil, fmt.Errorf("failed to get %s details for %s: %w",
						map[bool]string{true: "PR", false: "issue"}[isPR],
						match[0], err)
				}
				items = append(items, issue)
				seen[match[1]] = true
			}
		}
	} else {
		// When repoUrl is empty, just create basic Issue objects with numbers
		for _, match := range matches {
			num, _ := strconv.Atoi(match[1])
			if !seen[match[1]] {
				items = append(items, &origin.Issue{Number: num})
				seen[match[1]] = true
			}
		}
	}

	return items, nil
}

func containsIssue(items []*origin.Issue, item *origin.Issue) bool {
	if item == nil {
		return false
	}
	for _, existing := range items {
		if existing != nil && existing.Number == item.Number {
			return true
		}
	}
	return false
}
