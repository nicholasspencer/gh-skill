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

// SearchGists searches for skills by combining the user's own gists with
// a GitHub code search for repo-hosted skills.
func SearchGists(query string) ([]Gist, error) {
	seen := make(map[string]bool)
	var results []Gist

	// 1. User's own gists — GET /gists?per_page=100
	if out, err := exec.Command("gh", "api", "/gists?per_page=100").Output(); err == nil {
		var gists []Gist
		if err := json.Unmarshal(out, &gists); err == nil {
			for _, g := range gists {
				if seen[g.ID] {
					continue
				}
				desc := strings.ToLower(g.Description)
				if !strings.Contains(desc, "[gh-skill]") {
					continue
				}
				if !gistHasSkillFile(g) {
					continue
				}
				if query == "" || strings.Contains(desc, strings.ToLower(query)) {
					seen[g.ID] = true
					results = append(results, g)
				}
			}
		}
	}

	// 2. Code search for repo-hosted skills
	if query != "" {
		q := fmt.Sprintf("%s gh-skill filename:skill.md", query)
		endpoint := fmt.Sprintf("/search/code?q=%s&per_page=30", strings.ReplaceAll(q, " ", "+"))
		if out, err := exec.Command("gh", "api", endpoint).Output(); err == nil {
			var searchResp codeSearchResponse
			if err := json.Unmarshal(out, &searchResp); err == nil {
				for _, item := range searchResp.Items {
					repoFullName := item.Repository.FullName
					if seen[repoFullName] {
						continue
					}
					seen[repoFullName] = true
					// Convert code search hit to a Gist-like result
					results = append(results, Gist{
						ID:          repoFullName,
						Description: item.Repository.Description,
						HTMLURL:     item.Repository.HTMLURL,
						Files: map[string]GistFile{
							item.Name: {Filename: item.Name, RawURL: item.HTMLURL},
						},
						Owner: struct {
							Login string `json:"login"`
						}{Login: item.Repository.Owner.Login},
					})
				}
			}
		}
	}

	return results, nil
}

func gistHasSkillFile(g Gist) bool {
	for name := range g.Files {
		if IsSkillFile(name) {
			return true
		}
	}
	return false
}

// codeSearchResponse represents the GitHub code search API response.
type codeSearchResponse struct {
	Items []codeSearchItem `json:"items"`
}

type codeSearchItem struct {
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
	Repository struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
		Owner       struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}
