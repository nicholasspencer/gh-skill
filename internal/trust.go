package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const trustedAuthorsFile = "trusted-authors.json"

// TrustedAuthor represents a trusted gist author.
type TrustedAuthor struct {
	Username  string `json:"username"`
	TrustedAt string `json:"trusted_at"`
}

// TrustStore manages trusted authors.
type TrustStore struct {
	Authors []TrustedAuthor `json:"authors"`
}

func trustStorePath() string {
	return filepath.Join(SkillsBasePath(), trustedAuthorsFile)
}

// LoadTrustStore reads the trusted authors file.
func LoadTrustStore() (*TrustStore, error) {
	data, err := os.ReadFile(trustStorePath())
	if err != nil {
		if os.IsNotExist(err) {
			return &TrustStore{}, nil
		}
		return nil, err
	}
	var ts TrustStore
	if err := json.Unmarshal(data, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

// Save writes the trust store to disk.
func (ts *TrustStore) Save() error {
	if err := os.MkdirAll(SkillsBasePath(), 0755); err != nil {
		return err
	}
	data, _ := json.MarshalIndent(ts, "", "  ")
	return os.WriteFile(trustStorePath(), data, 0644)
}

// IsTrusted checks if an author is trusted.
func (ts *TrustStore) IsTrusted(username string) bool {
	for _, a := range ts.Authors {
		if strings.EqualFold(a.Username, username) {
			return true
		}
	}
	return false
}

// AddAuthor adds an author to the trust store.
func (ts *TrustStore) AddAuthor(username string) {
	if ts.IsTrusted(username) {
		return
	}
	ts.Authors = append(ts.Authors, TrustedAuthor{
		Username:  username,
		TrustedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// RemoveAuthor removes an author from the trust store.
func (ts *TrustStore) RemoveAuthor(username string) bool {
	for i, a := range ts.Authors {
		if strings.EqualFold(a.Username, username) {
			ts.Authors = append(ts.Authors[:i], ts.Authors[i+1:]...)
			return true
		}
	}
	return false
}

// AuthenticatedUser returns the current gh-authenticated username.
// Deprecated: Use Provider.AuthenticatedUser() instead.
func AuthenticatedUser() string {
	return (&GitHubProvider{}).AuthenticatedUser()
}

// IsScriptFile returns true if the filename looks like an executable script.
func IsScriptFile(filename string) bool {
	exts := []string{".sh", ".bash", ".zsh", ".py", ".rb", ".pl", ".js", ".ts"}
	lower := strings.ToLower(filename)
	for _, ext := range exts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// PromptTrust shows the trust gate and returns whether to proceed.
// Returns: "install", "trust-author", or "" (abort).
func PromptTrust(g *Gist, fm *FrontMatter) (string, error) {
	name := fm.Name
	if name == "" {
		name = g.ID
	}

	fmt.Println()
	fmt.Println("╭─────────────────────────────────────────╮")
	fmt.Println("│           ⚠️  Install Skill?             │")
	fmt.Println("╰─────────────────────────────────────────╯")
	fmt.Printf("  Name:   %s\n", name)
	fmt.Printf("  Author: %s\n", g.Owner.Login)
	fmt.Printf("  Gist:   %s\n", g.HTMLURL)
	fmt.Println()

	// File list
	fmt.Println("  Files:")
	var scripts []string
	for filename := range g.Files {
		expanded := ExpandFilename(filename)
		marker := ""
		if IsScriptFile(filename) {
			marker = " ⚡"
			scripts = append(scripts, expanded)
		}
		fmt.Printf("    %s%s\n", expanded, marker)
	}

	if len(scripts) > 0 {
		fmt.Println()
		fmt.Printf("  ⚠️  Contains %d script(s) — review before running\n", len(scripts))
	}

	// SKILL.md preview (first 20 lines after front matter)
	if sf, ok := g.Files["SKILL.md"]; ok {
		fmt.Println()
		fmt.Println("  ── SKILL.md preview ──")
		lines := strings.Split(sf.Content, "\n")
		// Skip front matter
		start := 0
		if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
			for i := 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) == "---" {
					start = i + 1
					break
				}
			}
		}
		count := 0
		for i := start; i < len(lines) && count < 20; i++ {
			fmt.Printf("  │ %s\n", lines[i])
			count++
		}
		if start+20 < len(lines) {
			fmt.Printf("  │ ... (%d more lines)\n", len(lines)-start-20)
		}
	}

	fmt.Println()
	fmt.Println("  [y] Install    [trust-author] Trust all from this author")
	fmt.Println("  [v] View full  [N] Abort")
	fmt.Print("  > ")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return "install", nil
		case "trust-author", "trust":
			return "trust-author", nil
		case "n", "no", "":
			return "", nil
		case "v", "view":
			// Show all file contents
			fmt.Println()
			for filename, file := range g.Files {
				fmt.Printf("  ══ %s ══\n", ExpandFilename(filename))
				for _, line := range strings.Split(file.Content, "\n") {
					fmt.Printf("  │ %s\n", line)
				}
				fmt.Println()
			}
			fmt.Println("  [y] Install    [trust-author] Trust all from this author")
			fmt.Println("  [N] Abort")
			fmt.Print("  > ")
		default:
			fmt.Print("  Invalid choice. [y/trust-author/v/N] > ")
		}
	}
}
