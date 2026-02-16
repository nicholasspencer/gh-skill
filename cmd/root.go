package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage AI agent skills stored as GitHub Gists and GitLab Snippets",
	Long:  "gh skill — install, publish, and manage AI agent skills backed by GitHub Gists and GitLab Snippets.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(forkCmd)
	rootCmd.AddCommand(initCmd)
}
