package internal

import "testing"

func TestDetectProvider(t *testing.T) {
	tests := []struct {
		input        string
		wantProvider string
		wantID       string
	}{
		// GitHub
		{"abc123", "github", "abc123"},
		{"https://gist.github.com/nico/abc123", "github", "abc123"},
		{"https://gist.github.com/nico/abc123/", "github", "abc123"},

		// GitLab
		{"https://gitlab.com/-/snippets/12345", "gitlab", "12345"},
		{"https://gitlab.com/snippets/67890", "gitlab", "67890"},
		{"https://gitlab.com/-/snippets/99999/", "gitlab", "99999"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			provider, id := DetectProvider(tt.input)
			if provider.Name() != tt.wantProvider {
				t.Errorf("DetectProvider(%q) provider = %q, want %q", tt.input, provider.Name(), tt.wantProvider)
			}
			if id != tt.wantID {
				t.Errorf("DetectProvider(%q) id = %q, want %q", tt.input, id, tt.wantID)
			}
		})
	}
}

func TestProviderByName(t *testing.T) {
	if p := ProviderByName("github"); p.Name() != "github" {
		t.Errorf("ProviderByName(github) = %q", p.Name())
	}
	if p := ProviderByName("gitlab"); p.Name() != "gitlab" {
		t.Errorf("ProviderByName(gitlab) = %q", p.Name())
	}
	if p := ProviderByName(""); p.Name() != "github" {
		t.Errorf("ProviderByName('') = %q, want github", p.Name())
	}
	if p := ProviderByName("unknown"); p.Name() != "github" {
		t.Errorf("ProviderByName(unknown) = %q, want github", p.Name())
	}
}
