package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GitHubTool struct {
	token string
}

func NewGitHubTool(token string) *GitHubTool {
	return &GitHubTool{token: token}
}

// retrieves the diff content of a pull request.
func (g *GitHubTool) FetchPRDiff(prURL string) (string, error) {
	// Simple conversion: https://github.com/user/repo/pull/1 -> https://github.com/user/repo/pull/1.diff
	diffURL := prURL
	if !strings.HasSuffix(diffURL, ".diff") {
		diffURL += ".diff"
	}

	req, err := http.NewRequest("GET", diffURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch diff: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func (g *GitHubTool) PostPRComment(prURL string, comment string) error {
	// Simple conversion: https://github.com/user/repo/pull/1 -> https://api.github.com/repos/user/repo/issues/1/comments
	// (GitHub treats PR comments as issue comments for general PR discussion)
	apiURL := strings.Replace(prURL, "github.com", "api.github.com/repos", 1)
	apiURL = strings.Replace(apiURL, "/pull/", "/issues/", 1)
	apiURL += "/comments"

	payload := map[string]string{"body": comment}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal comment payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github returned status %s: %s", resp.Status, string(body))
	}

	return nil
}
