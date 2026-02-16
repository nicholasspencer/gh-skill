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
	Use:     "fork <gist-url-or-id>",
	Aliases: []string{"steal"},
	Short:   "Fork a skill as your own gist",
	Long:    "Downloads a skill and re-publishes it as your own gist. Your copy, your rules.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gistID := internal.ParseGistID(args[0])

		fmt.Printf("Fetching gist %s...\n", gistID)
		gist, err := internal.FetchGist(gistID)
		if err != nil {
			return err
		}

		// Verify it's a skill
		_, _, ok := internal.FindSkillFile(gist.Files)
		if !ok {
			return fmt.Errorf("gist does not contain a *.skill.md file")
		}

		// Re-publish all files as a new gist under your account
		files := make(map[string]string)
		for name, file := range gist.Files {
			files[name] = file.Content
		}

		// Update description
		description := gist.Description
		if !strings.HasPrefix(description, "[gh-skill]") {
			description = "[gh-skill] " + description
		}

		visibility := "secret"
		if forkPublic {
			visibility = "public"
		}
		fmt.Printf("Forking %d files as a %s gist...\n", len(files), visibility)

		newGist, err := internal.CreateGist(description, files, forkPublic)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Forked: %s\n", newGist.HTMLURL)
		fmt.Printf("  Original: %s\n", gist.HTMLURL)
		fmt.Printf("  Install with: gh skill add %s\n", newGist.ID)
		return nil
	},
}

var forkLocalCmd = &cobra.Command{
	Use:     "fork <path>",
	Hidden:  true,
}

// forkFromDir publishes a local skill directory as your own gist (like publish but semantically "stealing")
var stealCmd = &cobra.Command{
	Use:   "steal <gist-url-or-id>",
	Short: "Alias for fork",
	Long:  "Downloads a skill and re-publishes it as your own gist. Same as fork.",
	Args:  cobra.ExactArgs(1),
	RunE:  forkCmd.RunE,
}

func init() {
	forkCmd.Flags().BoolVar(&forkPublic, "public", false, "Create a public gist")

	// Check if first arg is a local directory — if so, behave like publish
	originalRunE := forkCmd.RunE
	forkCmd.RunE = func(cmd *cobra.Command, args []string) error {
		path := args[0]
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			// It's a local directory — read skill, publish as own gist
			skillPath := filepath.Join(path, "SKILL.md")
			if _, err := os.Stat(skillPath); os.IsNotExist(err) {
				return fmt.Errorf("directory must contain a SKILL.md file")
			}

			skillContent, _ := os.ReadFile(skillPath)
			fm, _ := internal.ParseFrontMatter(string(skillContent))
			skillName := fm.Name
			if skillName == "" {
				skillName = filepath.Base(path)
			}
			description := fm.Description
			if description == "" {
				description = skillName
			}
			description = "[gh-skill] " + description

			files := make(map[string]string)
			filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
					return nil
				}
				rel, _ := filepath.Rel(path, p)
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

			fmt.Printf("✓ Published: %s\n", gist.HTMLURL)
			fmt.Printf("  Install with: gh skill add %s\n", gist.ID)
			return nil
		}
		return originalRunE(cmd, args)
	}
}
