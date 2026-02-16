package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var linkTarget string

var linkCmd = &cobra.Command{
	Use:   "link <name> --target <tool>",
	Short: "Link a skill to a specific tool's skill directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if linkTarget == "" {
			return fmt.Errorf("--target is required (claude-code, openclaw, openclaw/<agent>, copilot, cursor, codex, opencode)")
		}

		dir, err := internal.ToolDirByName(linkTarget)
		if err != nil {
			return err
		}

		if err := internal.LinkSkill(name, dir); err != nil {
			return err
		}

		fmt.Printf("✓ Linked %q → %s\n", name, dir)
		return nil
	},
}

func init() {
	linkCmd.Flags().StringVar(&linkTarget, "target", "", "Target tool (claude-code, openclaw[/<agent>], copilot, cursor, codex, opencode)")
}
