package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var publishPublic bool

var publishCmd = &cobra.Command{
	Use:   "publish <path>",
	Short: "Publish a local skill folder as a GitHub Gist",
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

		// Collect all files, flattening subdirectories
		files := make(map[string]string)
		filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() || fi.Name() == ".gistskill.json" || strings.HasPrefix(fi.Name(), ".") {
				return nil
			}
			rel, _ := filepath.Rel(dir, path)
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			gistName := internal.FlattenFilename(rel)
			files[gistName] = string(content)
			return nil
		})

		fmt.Printf("Publishing %d files as a %s gist...\n", len(files), map[bool]string{true: "public", false: "private"}[publishPublic])

		gist, err := internal.CreateGist(description, files, publishPublic)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Published: %s\n", gist.HTMLURL)
		fmt.Printf("  Install with: gh skill add %s\n", gist.ID)
		return nil
	},
}

func init() {
	publishCmd.Flags().BoolVar(&publishPublic, "public", true, "Create a public gist (default true)")
}
