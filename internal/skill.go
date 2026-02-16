package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const SkillsDir = ".gistskills"

// SkillMeta is the .gistskill.json metadata stored alongside installed skills.
type SkillMeta struct {
	Name        string `json:"name"`
	GistID      string `json:"gist_id"`
	CommitSHA   string `json:"commit_sha"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	GistURL     string `json:"gist_url"`
	InstalledAt string `json:"installed_at"`
	UpdatedAt   string `json:"updated_at"`
}

// FrontMatter represents YAML front matter in SKILL.md.
type FrontMatter struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Tags        []string `yaml:"tags"`
	Tools       []string `yaml:"tools"`
	Author      string   `yaml:"author"`
}

// SkillsBasePath returns the base path for installed skills.
func SkillsBasePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, SkillsDir)
}

// ParseFrontMatter extracts YAML front matter from a SKILL.md content string.
func ParseFrontMatter(content string) (*FrontMatter, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return &FrontMatter{}, nil // no front matter
	}
	var yamlLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		yamlLines = append(yamlLines, line)
	}
	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &fm); err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}
	return &fm, nil
}

// InstallSkill installs a gist as a skill.
func InstallSkill(g *Gist) (*SkillMeta, error) {
	// Check for SKILL.md
	skillFile, ok := g.Files["SKILL.md"]
	if !ok {
		return nil, fmt.Errorf("gist does not contain a SKILL.md file")
	}

	// Parse front matter for name
	fm, err := ParseFrontMatter(skillFile.Content)
	if err != nil {
		return nil, err
	}

	name := fm.Name
	if name == "" {
		// Fallback: use gist ID
		name = g.ID
	}

	// Create skill directory
	skillDir := filepath.Join(SkillsBasePath(), name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Write all files directly (gists are flat, no subdirectory expansion)
	for filename, file := range g.Files {
		destPath := filepath.Join(skillDir, filename)
		if err := os.WriteFile(destPath, []byte(file.Content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	// Build metadata
	commitSHA := ""
	if len(g.History) > 0 {
		commitSHA = g.History[0].Version
	}

	meta := &SkillMeta{
		Name:        name,
		GistID:      g.ID,
		CommitSHA:   commitSHA,
		Description: fm.Description,
		Version:     fm.Version,
		Author:      g.Owner.Login,
		GistURL:     g.HTMLURL,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Write metadata
	metaPath := filepath.Join(skillDir, ".gistskill.json")
	metaData, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}

	return meta, nil
}

// ListSkills lists all installed skills.
func ListSkills() ([]SkillMeta, error) {
	base := SkillsBasePath()
	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var skills []SkillMeta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		metaPath := filepath.Join(base, e.Name(), ".gistskill.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}
		var meta SkillMeta
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		skills = append(skills, meta)
	}
	return skills, nil
}

// GetSkill reads metadata for a single skill.
func GetSkill(name string) (*SkillMeta, error) {
	metaPath := filepath.Join(SkillsBasePath(), name, ".gistskill.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("skill %q not found", name)
	}
	var meta SkillMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// RemoveSkill removes an installed skill and its symlinks.
func RemoveSkill(name string) error {
	skillDir := filepath.Join(SkillsBasePath(), name)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return fmt.Errorf("skill %q not found", name)
	}

	// Remove symlinks from tool directories
	for _, dir := range DetectToolDirs() {
		link := filepath.Join(dir, name)
		target, err := os.Readlink(link)
		if err == nil && strings.HasPrefix(target, SkillsBasePath()) {
			os.Remove(link)
		}
	}

	return os.RemoveAll(skillDir)
}
