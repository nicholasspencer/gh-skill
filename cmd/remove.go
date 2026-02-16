package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Short:   "Remove an installed skill",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := internal.RemoveSkill(name); err != nil {
			return err
		}
		fmt.Printf("✓ Removed skill %q\n", name)
		return nil
	},
}
