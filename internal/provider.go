package internal

import (
	"regexp"
	"strings"
)

// Provider abstracts snippet storage backends (GitHub Gists, GitLab Snippets).
type Provider interface {
	Name() string
	FetchSnippet(id string) (*Gist, error)
	CreateSnippet(description string, files map[string]string, public bool) (*Gist, error)
	SearchSnippets(query string) ([]Gist, error)
	AuthenticatedUser() string
}

var gitlabSnippetRe = regexp.MustCompile(`gitlab\.com/(?:-/)?snippets/(\d+)`)

// DetectProvider examines a URL or ID and returns the appropriate provider and extracted ID.
func DetectProvider(input string) (Provider, string) {
	input = strings.TrimSpace(input)

	if strings.Contains(input, "gitlab.com") {
		if m := gitlabSnippetRe.FindStringSubmatch(input); len(m) == 2 {
			return &GitLabProvider{}, m[1]
		}
	}

	// Default to GitHub — ParseGistID handles gist.github.com URLs and bare IDs
	return &GitHubProvider{}, ParseGistID(input)
}

// ProviderByName returns a provider by name string. Defaults to GitHub.
func ProviderByName(name string) Provider {
	switch strings.ToLower(name) {
	case "gitlab":
		return &GitLabProvider{}
	default:
		return &GitHubProvider{}
	}
}
