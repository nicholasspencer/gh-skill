package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var searchProvider string

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for skills on GitHub Gists and GitLab Snippets",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := internal.ProviderByName(searchProvider)
		results, err := provider.SearchSnippets(args[0])
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Println("No skills found. Try a different query.")
			return nil
		}
		for _, g := range results {
			fmt.Printf("%-20s %s\n", g.ID, g.Description)
			fmt.Printf("  → gh skill add %s\n\n", g.ID)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchProvider, "provider", "github", "Provider to search (github, gitlab)")
}
