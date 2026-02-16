package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <gist-url-or-id>",
	Short: "Install a skill from a GitHub Gist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gistID := internal.ParseGistID(args[0])
		fmt.Printf("Fetching gist %s...\n", gistID)

		gist, err := internal.FetchGist(gistID)
		if err != nil {
			return err
		}

		meta, err := internal.InstallSkill(gist)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Installed skill %q (v%s)\n", meta.Name, meta.Version)

		// Auto-link to detected tools
		linked := internal.AutoLink(meta.Name)
		for _, dir := range linked {
			fmt.Printf("  → Linked to %s\n", dir)
		}

		return nil
	},
}
