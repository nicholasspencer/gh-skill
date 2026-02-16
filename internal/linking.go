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

	// Static tool targets
	staticTargets := []struct {
		name string
		dir  string
	}{
		{"claude-code", filepath.Join(home, ".claude", "skills")},
		{"copilot", filepath.Join(home, ".copilot", "skills")},
		{"codex", filepath.Join(home, ".codex", "skills")},
		{"opencode", filepath.Join(home, ".opencode", "skills")},
	}

	for _, t := range staticTargets {
		parent := filepath.Dir(t.dir)
		if _, err := os.Stat(parent); err == nil {
			dirs = append(dirs, t.dir)
		}
	}

	// OpenClaw agent targets (multiple agents, each with their own skills dir)
	for _, t := range openclawAgentTargets(home) {
		parent := filepath.Dir(t.Dir)
		if _, err := os.Stat(parent); err == nil {
			dirs = append(dirs, t.Dir)
		}
	}

	return dirs
}

// KnownTools returns all known tool targets.
func KnownTools() []ToolTarget {
	home, _ := os.UserHomeDir()

	tools := []ToolTarget{
		{"claude-code", filepath.Join(home, ".claude", "skills")},
		{"copilot", filepath.Join(home, ".copilot", "skills")},
		{"cursor", filepath.Join(".cursor", "skills")}, // project-level
		{"codex", filepath.Join(home, ".codex", "skills")},
		{"opencode", filepath.Join(home, ".opencode", "skills")},
	}

	// Append all detected OpenClaw agents
	tools = append(tools, openclawAgentTargets(home)...)

	return tools
}

// openclawConfig represents the relevant parts of openclaw.json
type openclawConfig struct {
	Agents struct {
		Defaults struct {
			Workspace string `json:"workspace"`
		} `json:"defaults"`
		List []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Workspace string `json:"workspace"`
		} `json:"list"`
	} `json:"agents"`
}

// openclawAgentTargets reads ~/.openclaw/openclaw.json and returns a ToolTarget
// for each configured agent's skills directory. Falls back to the legacy
// ~/.chad/skills/ if no config is found.
func openclawAgentTargets(home string) []ToolTarget {
	configPath := filepath.Join(home, ".openclaw", "openclaw.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Fallback: no openclaw config, try legacy path
		return []ToolTarget{
			{"openclaw", filepath.Join(home, ".chad", "skills")},
		}
	}

	var cfg openclawConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return []ToolTarget{
			{"openclaw", filepath.Join(home, ".chad", "skills")},
		}
	}

	defaultWorkspace := cfg.Agents.Defaults.Workspace
	if defaultWorkspace == "" {
		defaultWorkspace = filepath.Join(home, ".chad")
	}

	var targets []ToolTarget

	for _, agent := range cfg.Agents.List {
		workspace := agent.Workspace
		if workspace == "" {
			workspace = defaultWorkspace
		}
		name := agent.Name
		if name == "" {
			name = agent.ID
		}
		targets = append(targets, ToolTarget{
			Name: "openclaw/" + name,
			Dir:  filepath.Join(workspace, "skills"),
		})
	}

	// If no agents configured, fall back to default workspace
	if len(targets) == 0 {
		targets = append(targets, ToolTarget{
			Name: "openclaw",
			Dir:  filepath.Join(defaultWorkspace, "skills"),
		})
	}

	return targets
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
// Supports "openclaw" as a shorthand that matches the first openclaw agent,
// as well as specific "openclaw/<agent>" targets.
func ToolDirByName(name string) (string, error) {
	for _, t := range KnownTools() {
		if t.Name == name {
			return t.Dir, nil
		}
	}
	// Allow bare "openclaw" to match first openclaw agent
	if name == "openclaw" {
		for _, t := range KnownTools() {
			if len(t.Name) > 8 && t.Name[:9] == "openclaw/" {
				return t.Dir, nil
			}
		}
	}
	return "", fmt.Errorf("unknown tool %q (known: claude-code, openclaw, openclaw/<agent>, copilot, cursor, codex, opencode)", name)
}
