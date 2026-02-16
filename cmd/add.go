package cmd

import (
	"fmt"
	"strings"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var (
	addYes   bool
	addIdgaf bool
)

var addCmd = &cobra.Command{
	Use:   "add <gist-url-or-id>",
	Short: "Install a skill from a GitHub Gist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, snippetID := internal.DetectProvider(args[0])
		fmt.Printf("Fetching %s snippet %s...\n", provider.Name(), snippetID)

		gist, err := provider.FetchSnippet(snippetID)
		if err != nil {
			return err
		}

		// Find skill file (*.skill.md or legacy SKILL.md)
		_, skillFile, ok := internal.FindSkillFile(gist.Files)
		if !ok {
			return fmt.Errorf("gist does not contain a *.skill.md file")
		}

		fm, err := internal.ParseFrontMatter(skillFile.Content)
		if err != nil {
			return err
		}

		// Trust gate
		skipPrompt := addYes || addIdgaf
		if !skipPrompt {
			// Own gists/snippets are implicitly trusted
			if authUser := provider.AuthenticatedUser(); authUser != "" && strings.EqualFold(authUser, gist.Owner.Login) {
				skipPrompt = true
			}
		}
		if !skipPrompt {
			ts, err := internal.LoadTrustStore()
			if err != nil {
				return err
			}
			if ts.IsTrusted(gist.Owner.Login) {
				skipPrompt = true
				fmt.Printf("Author %q is trusted.\n", gist.Owner.Login)
			}
		}

		if !skipPrompt {
			decision, err := internal.PromptTrust(gist, fm)
			if err != nil {
				return err
			}
			switch decision {
			case "":
				fmt.Println("Aborted.")
				return nil
			case "trust-author":
				ts, _ := internal.LoadTrustStore()
				ts.AddAuthor(gist.Owner.Login)
				if err := ts.Save(); err != nil {
					return fmt.Errorf("failed to save trust store: %w", err)
				}
				fmt.Printf("✓ Trusted author %q for future installs.\n", gist.Owner.Login)
			}
		}

		meta, err := internal.InstallSkill(gist, provider.Name())
		if err != nil {
			return err
		}

		fmt.Printf("✓ Installed skill %q (v%s)\n", meta.Name, meta.Version)

		// Auto-link to detected tools
		linked := internal.AutoLink(meta.Name)
		for _, dir := range linked {
			fmt.Printf("  → Linked to %s\n", dir)
		}

		// Lazy init: install meta skill if not present
		ensureMetaSkill(linked)

		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addYes, "yes", "y", false, "Skip trust prompt")
	addCmd.Flags().BoolVar(&addIdgaf, "idgaf", false, "Skip trust prompt (alias)")
}
