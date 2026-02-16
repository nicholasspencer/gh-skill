package internal

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GistFile represents a file within a gist.
type GistFile struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	RawURL   string `json:"raw_url"`
}

// Gist represents a GitHub Gist.
type Gist struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Files       map[string]GistFile `json:"files"`
	HTMLURL     string              `json:"html_url"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
	History []struct {
		Version string `json:"version"`
	} `json:"history"`
}

// FetchGist fetches a gist by ID using the gh CLI.
func FetchGist(gistID string) (*Gist, error) {
	out, err := exec.Command("gh", "api", fmt.Sprintf("/gists/%s", gistID)).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gist %s: %w", gistID, err)
	}
	var g Gist
	if err := json.Unmarshal(out, &g); err != nil {
		return nil, fmt.Errorf("failed to parse gist response: %w", err)
	}
	return &g, nil
}

// CreateGist creates a new gist using the gh CLI.
func CreateGist(description string, files map[string]string, public bool) (*Gist, error) {
	gistFiles := make(map[string]map[string]string)
	for name, content := range files {
		gistFiles[name] = map[string]string{"content": content}
	}
	payload := map[string]interface{}{
		"description": description,
		"public":      public,
		"files":       gistFiles,
	}
	data, _ := json.Marshal(payload)
	cmd := exec.Command("gh", "api", "/gists", "--method", "POST", "--input", "-")
	cmd.Stdin = strings.NewReader(string(data))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to create gist: %w", err)
	}
	var g Gist
	if err := json.Unmarshal(out, &g); err != nil {
		return nil, fmt.Errorf("failed to parse gist response: %w", err)
	}
	return &g, nil
}

// ParseGistID extracts a gist ID from a URL or returns the input if already an ID.
func ParseGistID(input string) string {
	input = strings.TrimSpace(input)
	// Handle URLs like https://gist.github.com/user/abc123
	if strings.Contains(input, "gist.github.com") {
		parts := strings.Split(strings.TrimRight(input, "/"), "/")
		return parts[len(parts)-1]
	}
	return input
}

// SearchGists searches for gists matching a query via GitHub code search.
func SearchGists(query string) ([]Gist, error) {
	// Search gists using gh api - search for gists with #gistskill tag
	searchQuery := fmt.Sprintf("%s #gistskill", query)
	out, err := exec.Command("gh", "api", fmt.Sprintf("/gists/public?per_page=30")).Output()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	var gists []Gist
	if err := json.Unmarshal(out, &gists); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}
	// Filter locally for query match and #gistskill tag
	var results []Gist
	lowerQuery := strings.ToLower(searchQuery)
	_ = lowerQuery
	for _, g := range gists {
		desc := strings.ToLower(g.Description)
		if !strings.Contains(desc, "#gistskill") {
			continue
		}
		if query == "" || strings.Contains(desc, strings.ToLower(query)) {
			results = append(results, g)
		}
		// Also check if SKILL.md exists
		if _, ok := g.Files["SKILL.md"]; ok {
			// Already added above if description matches
		}
	}
	return results, nil
}
