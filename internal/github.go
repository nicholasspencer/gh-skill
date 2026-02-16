package internal

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GitHubProvider implements Provider using gh CLI for GitHub Gists.
type GitHubProvider struct{}

func (p *GitHubProvider) Name() string { return "github" }

func (p *GitHubProvider) FetchSnippet(id string) (*Gist, error) {
	return FetchGist(id)
}

func (p *GitHubProvider) CreateSnippet(description string, files map[string]string, public bool) (*Gist, error) {
	return CreateGist(description, files, public)
}

func (p *GitHubProvider) SearchSnippets(query string) ([]Gist, error) {
	return SearchGists(query)
}

func (p *GitHubProvider) AuthenticatedUser() string {
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
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

// SearchGists searches for gists matching a query via GitHub API.
func SearchGists(query string) ([]Gist, error) {
	out, err := exec.Command("gh", "api", "/gists/public?per_page=100").Output()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	var gists []Gist
	if err := json.Unmarshal(out, &gists); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	var results []Gist
	for _, g := range gists {
		desc := strings.ToLower(g.Description)
		if !strings.Contains(desc, "[gh-skill]") {
			continue
		}
		hasSkillFile := false
		for name := range g.Files {
			if IsSkillFile(name) {
				hasSkillFile = true
				break
			}
		}
		if !hasSkillFile {
			continue
		}
		if query == "" || strings.Contains(desc, strings.ToLower(query)) {
			results = append(results, g)
		}
	}
	return results, nil
}
