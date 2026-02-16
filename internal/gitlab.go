package internal

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GitLabProvider implements Provider using glab CLI for GitLab Snippets.
type GitLabProvider struct{}

func (p *GitLabProvider) Name() string { return "gitlab" }

// gitlabSnippet is the JSON shape returned by the GitLab snippets API.
type gitlabSnippet struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	WebURL      string `json:"web_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Author      struct {
		Username string `json:"username"`
	} `json:"author"`
	Files []struct {
		Path   string `json:"path"`
		RawURL string `json:"raw_url"`
	} `json:"files"`
}

func (s *gitlabSnippet) toGist() *Gist {
	g := &Gist{
		ID:          fmt.Sprintf("%d", s.ID),
		Description: s.Description,
		HTMLURL:     s.WebURL,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Files:       make(map[string]GistFile),
	}
	if s.Description == "" {
		g.Description = s.Title
	}
	g.Owner.Login = s.Author.Username
	for _, f := range s.Files {
		g.Files[f.Path] = GistFile{
			Filename: f.Path,
			RawURL:   f.RawURL,
		}
	}
	return g
}

func (p *GitLabProvider) FetchSnippet(id string) (*Gist, error) {
	out, err := exec.Command("glab", "api", fmt.Sprintf("/snippets/%s", id)).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch snippet %s: %w", id, err)
	}
	var s gitlabSnippet
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, fmt.Errorf("failed to parse snippet response: %w", err)
	}
	g := s.toGist()

	// Fetch raw content for each file
	for _, f := range s.Files {
		rawOut, err := exec.Command("glab", "api", fmt.Sprintf("/snippets/%s/files/main/%s/raw", id, f.Path)).Output()
		if err != nil {
			continue
		}
		gf := g.Files[f.Path]
		gf.Content = string(rawOut)
		g.Files[f.Path] = gf
	}

	return g, nil
}

func (p *GitLabProvider) CreateSnippet(description string, files map[string]string, public bool) (*Gist, error) {
	visibility := "private"
	if public {
		visibility = "public"
	}

	type snippetFile struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}
	var sf []snippetFile
	title := description
	for name, content := range files {
		sf = append(sf, snippetFile{FilePath: name, Content: content})
	}

	payload := map[string]interface{}{
		"title":       title,
		"description": description,
		"visibility":  visibility,
		"files":       sf,
	}
	data, _ := json.Marshal(payload)
	cmd := exec.Command("glab", "api", "/snippets", "--method", "POST", "--input", "-")
	cmd.Stdin = strings.NewReader(string(data))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to create snippet: %w", err)
	}
	var s gitlabSnippet
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, fmt.Errorf("failed to parse snippet response: %w", err)
	}
	return s.toGist(), nil
}

func (p *GitLabProvider) SearchSnippets(query string) ([]Gist, error) {
	encodedQuery := strings.ReplaceAll(query, " ", "+")
	out, err := exec.Command("glab", "api", fmt.Sprintf("/snippets/public?per_page=100&search=%s", encodedQuery)).Output()
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	var snippets []gitlabSnippet
	if err := json.Unmarshal(out, &snippets); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	var results []Gist
	for _, s := range snippets {
		desc := strings.ToLower(s.Description)
		if s.Description == "" {
			desc = strings.ToLower(s.Title)
		}
		if !strings.Contains(desc, "[gh-skill]") {
			continue
		}
		if query == "" || strings.Contains(desc, strings.ToLower(query)) {
			results = append(results, *s.toGist())
		}
	}
	return results, nil
}

func (p *GitLabProvider) AuthenticatedUser() string {
	out, err := exec.Command("glab", "api", "/user", "--jq", ".username").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
