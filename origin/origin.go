package origin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func GetIssueDetails(issue *Issue, repoUrl, issueNumber string) error {
	owner, repo := getOwnerAndRepo(repoUrl)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s", owner, repo, issueNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "go-changelog")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get issue details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get issue details: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(issue); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func getOwnerAndRepo(repoUrl string) (owner string, repo string) {
	// parse repoUrl to get owner and repo
	parts := strings.Split(repoUrl, "/")
	return parts[len(parts)-2], parts[len(parts)-1]
}
