package internal

import (
	"strings"
)

// GistFile represents a file within a gist or snippet.
type GistFile struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
	RawURL   string `json:"raw_url"`
}

// Gist represents a GitHub Gist or GitLab Snippet (universal struct).
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

// ParseGistID extracts a gist ID from a URL or returns the input if already an ID.
func ParseGistID(input string) string {
	input = strings.TrimSpace(input)
	if strings.Contains(input, "gist.github.com") {
		parts := strings.Split(strings.TrimRight(input, "/"), "/")
		return parts[len(parts)-1]
	}
	return input
}
