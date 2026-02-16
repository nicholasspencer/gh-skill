package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ToolTarget represents a known AI tool's skill directory.
type ToolTarget struct {
	Name string
	Dir  string
}

// DetectToolDirs returns paths to skill directories for detected AI tools.
func DetectToolDirs() []string {
	home, _ := os.UserHomeDir()
	var dirs []string

	targets := []struct {
		name string
		dir  string
	}{
		{"claude-code", filepath.Join(home, ".claude", "skills")},
		{"openclaw", openclawSkillsDir(home)},
		{"codex", filepath.Join(home, ".codex", "skills")},
		{"opencode", filepath.Join(home, ".opencode", "skills")},
	}

	for _, t := range targets {
		// Check if parent tool dir exists (e.g., ~/.claude/)
		parent := filepath.Dir(t.dir)
		if _, err := os.Stat(parent); err == nil {
			dirs = append(dirs, t.dir)
		}
	}
	return dirs
}

// KnownTools returns all known tool targets.
func KnownTools() []ToolTarget {
	home, _ := os.UserHomeDir()
	return []ToolTarget{
		{"claude-code", filepath.Join(home, ".claude", "skills")},
		{"openclaw", openclawSkillsDir(home)},
		{"cursor", filepath.Join(".cursor", "skills")}, // project-level
		{"codex", filepath.Join(home, ".codex", "skills")},
		{"opencode", filepath.Join(home, ".opencode", "skills")},
	}
}

func openclawSkillsDir(home string) string {
	// Try to read openclaw config for skills_dir
	configPath := filepath.Join(home, ".chad", "openclaw.json")
	data, err := os.ReadFile(configPath)
	if err == nil {
		var cfg map[string]interface{}
		if json.Unmarshal(data, &cfg) == nil {
			if dir, ok := cfg["skills_dir"].(string); ok && dir != "" {
				return dir
			}
		}
	}
	return filepath.Join(home, ".chad", "skills")
}

// LinkSkill creates a symlink for a skill in the given tool directory.
func LinkSkill(skillName, toolDir string) error {
	skillDir := filepath.Join(SkillsBasePath(), skillName)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return fmt.Errorf("skill %q not found", skillName)
	}

	if err := os.MkdirAll(toolDir, 0755); err != nil {
		return fmt.Errorf("failed to create tool directory %s: %w", toolDir, err)
	}

	linkPath := filepath.Join(toolDir, skillName)
	// Remove existing symlink if present
	os.Remove(linkPath)

	if err := os.Symlink(skillDir, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	return nil
}

// AutoLink links a skill to all detected tool directories.
func AutoLink(skillName string) []string {
	var linked []string
	for _, dir := range DetectToolDirs() {
		if err := LinkSkill(skillName, dir); err == nil {
			linked = append(linked, dir)
		}
	}
	return linked
}

// ToolDirByName returns the skill directory for a named tool.
func ToolDirByName(name string) (string, error) {
	for _, t := range KnownTools() {
		if t.Name == name {
			return t.Dir, nil
		}
	}
	return "", fmt.Errorf("unknown tool %q (known: claude-code, openclaw, cursor, codex, opencode)", name)
}
