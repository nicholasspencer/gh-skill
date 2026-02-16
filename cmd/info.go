package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show details about an installed skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		meta, err := internal.GetSkill(args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Name:        %s\n", meta.Name)
		fmt.Printf("Version:     %s\n", meta.Version)
		fmt.Printf("Description: %s\n", meta.Description)
		fmt.Printf("Author:      %s\n", meta.Author)
		fmt.Printf("Gist:        %s\n", meta.GistURL)
		fmt.Printf("Provider:    %s\n", meta.EffectiveProvider())
		fmt.Printf("Gist ID:     %s\n", meta.GistID)
		fmt.Printf("Commit:      %s\n", meta.CommitSHA)
		fmt.Printf("Installed:   %s\n", meta.InstalledAt)
		fmt.Printf("Updated:     %s\n", meta.UpdatedAt)

		// List files
		skillDir := filepath.Join(internal.SkillsBasePath(), meta.Name)
		fmt.Println("\nFiles:")
		filepath.Walk(skillDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || info.Name() == ".gistskill.json" {
				return nil
			}
			rel, _ := filepath.Rel(skillDir, path)
			fmt.Printf("  %s\n", rel)
			return nil
		})

		return nil
	},
}
