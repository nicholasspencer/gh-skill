package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

//go:embed init_skill/*
var initSkillFS embed.FS

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Install the gh-skill meta skill into detected tool directories",
	Long:  "Copies the gh-skill meta skill (which teaches agents how to search and install skills) into all detected AI tool skill directories.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirs := internal.DetectToolDirs()
		if len(dirs) == 0 {
			fmt.Println("No AI tool skill directories detected.")
			fmt.Println("Supported tools: claude-code, openclaw (multi-agent), copilot, codex, opencode")
			return nil
		}

		bundled := []struct {
			name string
			file string
		}{
			{"gh-skill", "init_skill/SKILL.md"},
			{"skill-creator", "init_skill/skill-creator.SKILL.md"},
		}

		installed := 0
		for _, dir := range dirs {
			for _, skill := range bundled {
				content, err := initSkillFS.ReadFile(skill.file)
				if err != nil {
					fmt.Printf("⚠️  Failed to read embedded %s skill: %v\n", skill.name, err)
					continue
				}
				skillDir := filepath.Join(dir, skill.name)
				if err := os.MkdirAll(skillDir, 0755); err != nil {
					fmt.Printf("⚠️  Failed to create %s: %v\n", skillDir, err)
					continue
				}
				destPath := filepath.Join(skillDir, "SKILL.md")
				if err := os.WriteFile(destPath, content, 0644); err != nil {
					fmt.Printf("⚠️  Failed to write %s: %v\n", destPath, err)
					continue
				}
				fmt.Printf("✓ Installed %s skill to %s\n", skill.name, skillDir)
				installed++
			}
		}

		if installed > 0 {
			fmt.Printf("\n%d skill(s) installed across detected tools.\n", installed)
		}
		return nil
	},
}

// ensureMetaSkill installs bundled skills into detected tool directories
// if not already present. Called lazily from `gh skill add`.
func ensureMetaSkill(_ []string) {
	bundled := []struct {
		name string
		file string
	}{
		{"gh-skill", "init_skill/SKILL.md"},
		{"skill-creator", "init_skill/skill-creator.SKILL.md"},
	}

	for _, dir := range internal.DetectToolDirs() {
		for _, skill := range bundled {
			destDir := filepath.Join(dir, skill.name)
			destPath := filepath.Join(destDir, "SKILL.md")
			if _, err := os.Stat(destPath); err == nil {
				continue
			}
			content, err := initSkillFS.ReadFile(skill.file)
			if err != nil {
				continue
			}
			if err := os.MkdirAll(destDir, 0755); err != nil {
				continue
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				continue
			}
			fmt.Printf("  → Installed %s skill to %s\n", skill.name, destDir)
		}
	}
}

func init() {
	// registered in root.go
}
