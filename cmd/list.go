package cmd

import (
	"fmt"
	"text/tabwriter"
	"os"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		skills, err := internal.ListSkills()
		if err != nil {
			return err
		}
		if len(skills) == 0 {
			fmt.Println("No skills installed. Use `gh skill add <gist>` to install one.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tGIST\tINSTALLED")
		for _, s := range skills {
			version := s.Version
			if version == "" {
				version = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, version, s.GistID, s.InstalledAt[:10])
		}
		return w.Flush()
	},
}
