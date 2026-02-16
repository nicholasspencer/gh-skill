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
	publishPublic bool
	publishSecret bool
)

var publishCmd = &cobra.Command{
	Use:   "publish <path>",
	Short: "Publish a local skill folder as a GitHub Gist",
	Long:  "Creates a secret (unlisted) gist by default. Use --public to make it discoverable.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("%s is not a directory", dir)
		}

		// Check for SKILL.md
		skillPath := filepath.Join(dir, "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			return fmt.Errorf("directory must contain a SKILL.md file")
		}

		// Read SKILL.md for description
		skillContent, _ := os.ReadFile(skillPath)
		fm, _ := internal.ParseFrontMatter(string(skillContent))
		description := fm.Description
		if description == "" {
			description = fm.Name
		}
		description += " #gistskill"

		// Determine visibility: secret by default, --public overrides, --secret is explicit
		isPublic := false
		if publishPublic {
			isPublic = true
		}

		// Collect top-level files only (gists are flat)
		files := make(map[string]string)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				continue
			}
			files[e.Name()] = string(content)
		}

		visibility := "secret"
		if isPublic {
			visibility = "public"
		}
		fmt.Printf("Publishing %d files as a %s gist...\n", len(files), visibility)

		gist, err := internal.CreateGist(description, files, isPublic)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Published: %s\n", gist.HTMLURL)
		fmt.Printf("  Install with: gh skill add %s\n", gist.ID)
		return nil
	},
}

func init() {
	publishCmd.Flags().BoolVar(&publishPublic, "public", false, "Create a public gist")
	publishCmd.Flags().BoolVar(&publishSecret, "secret", false, "Create a secret (unlisted) gist (default)")
}
