package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var (
	forkPublic bool
)

var forkCmd = &cobra.Command{
	Use:     "fork <skill-name>",
	Aliases: []string{"steal"},
	Short:   "Fork a skill as your own gist",
	Long: `Finds the skill by name (searches installed skills, then nearby skill directories),
and re-publishes it as your own gist. Your copy, your rules.

Search order:
  1. ~/.gistskills/<name>/
  2. ./skills/<name>/  (current directory)
  3. Walk up to find nearest skills/<name>/
  4. ./<name>/SKILL.md (bare directory)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Find the skill directory
		skillDir, source, err := findSkillDir(name)
		if err != nil {
			return err
		}

		fmt.Printf("Found skill %q at %s (%s)\n", name, skillDir, source)

		// Read SKILL.md
		skillPath := filepath.Join(skillDir, "SKILL.md")
		skillContent, err := os.ReadFile(skillPath)
		if err != nil {
			return fmt.Errorf("could not read SKILL.md in %s: %w", skillDir, err)
		}

		fm, _ := internal.ParseFrontMatter(string(skillContent))
		skillName := fm.Name
		if skillName == "" {
			skillName = name
		}
		description := fm.Description
		if description == "" {
			description = skillName
		}
		description = "[gh-skill] " + description

		// Collect all files
		files := make(map[string]string)
		filepath.Walk(skillDir, func(p string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
				return nil
			}
			rel, _ := filepath.Rel(skillDir, p)
			content, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			gistName := internal.FlattenFilename(rel)
			if gistName == "SKILL.md" {
				gistName = internal.SkillFileName(skillName)
			}
			files[gistName] = string(content)
			return nil
		})

		visibility := "secret"
		if forkPublic {
			visibility = "public"
		}
		fmt.Printf("Publishing %d files as a %s gist...\n", len(files), visibility)

		gist, err := internal.CreateGist(description, files, forkPublic)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Forked: %s\n", gist.HTMLURL)
		fmt.Printf("  Install with: gh skill add %s\n", gist.ID)
		return nil
	},
}

// findSkillDir searches for a skill by name in standard locations.
// Returns (path, source_description, error).
func findSkillDir(name string) (string, string, error) {
	// 1. Check ~/.gistskills/<name>/
	managed := filepath.Join(internal.SkillsBasePath(), name)
	if hasSkillMD(managed) {
		return managed, "installed", nil
	}

	// 2. Check ./skills/<name>/ (current directory)
	cwd, _ := os.Getwd()
	local := filepath.Join(cwd, "skills", name)
	if hasSkillMD(local) {
		return local, "local", nil
	}

	// 3. Walk up directories looking for skills/<name>/
	dir := cwd
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break // hit root
		}
		dir = parent
		candidate := filepath.Join(dir, "skills", name)
		if hasSkillMD(candidate) {
			return candidate, "ancestor", nil
		}
	}

	// 4. Check ./<name>/SKILL.md (bare directory in CWD)
	bare := filepath.Join(cwd, name)
	if hasSkillMD(bare) {
		return bare, "directory", nil
	}

	return "", "", fmt.Errorf("skill %q not found\n\nSearched:\n  ~/.gistskills/%s/\n  ./skills/%s/\n  (ancestor dirs)/skills/%s/\n  ./%s/", name, name, name, name, name)
}

func hasSkillMD(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "SKILL.md"))
	return err == nil
}

func init() {
	forkCmd.Flags().BoolVar(&forkPublic, "public", false, "Create a public gist")
}
