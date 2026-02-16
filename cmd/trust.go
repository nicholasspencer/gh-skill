package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var (
	trustList   bool
	trustRemove string
)

var trustCmd = &cobra.Command{
	Use:   "trust [username]",
	Short: "Manage trusted authors",
	Long:  "Add, list, or remove trusted authors. Skills from trusted authors install without a trust prompt.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ts, err := internal.LoadTrustStore()
		if err != nil {
			return err
		}

		// --list
		if trustList {
			if len(ts.Authors) == 0 {
				fmt.Println("No trusted authors.")
				return nil
			}
			for _, a := range ts.Authors {
				fmt.Printf("  %s (trusted %s)\n", a.Username, a.TrustedAt[:10])
			}
			return nil
		}

		// --remove
		if trustRemove != "" {
			if ts.RemoveAuthor(trustRemove) {
				if err := ts.Save(); err != nil {
					return err
				}
				fmt.Printf("✓ Removed %q from trusted authors.\n", trustRemove)
			} else {
				fmt.Printf("Author %q was not trusted.\n", trustRemove)
			}
			return nil
		}

		// Add
		if len(args) == 0 {
			return fmt.Errorf("provide a username, or use --list / --remove")
		}

		ts.AddAuthor(args[0])
		if err := ts.Save(); err != nil {
			return err
		}
		fmt.Printf("✓ Trusted author %q.\n", args[0])
		return nil
	},
}

func init() {
	trustCmd.Flags().BoolVar(&trustList, "list", false, "List trusted authors")
	trustCmd.Flags().StringVar(&trustRemove, "remove", "", "Remove a trusted author")
	rootCmd.AddCommand(trustCmd)
}
