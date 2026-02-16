package cmd

import (
	"fmt"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var updateAll bool

var updateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a skill to the latest gist revision",
	RunE: func(cmd *cobra.Command, args []string) error {
		if updateAll {
			skills, err := internal.ListSkills()
			if err != nil {
				return err
			}
			if len(skills) == 0 {
				fmt.Println("No skills installed.")
				return nil
			}
			for _, s := range skills {
				if err := updateSkill(s.Name, s.GistID); err != nil {
					fmt.Printf("✗ Failed to update %s: %v\n", s.Name, err)
				}
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("provide a skill name or use --all")
		}

		meta, err := internal.GetSkill(args[0])
		if err != nil {
			return err
		}
		return updateSkill(meta.Name, meta.GistID)
	},
}

func updateSkill(name, gistID string) error {
	gist, err := internal.FetchGist(gistID)
	if err != nil {
		return err
	}
	meta, err := internal.InstallSkill(gist)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Updated %q to v%s\n", meta.Name, meta.Version)
	return nil
}

func init() {
	updateCmd.Flags().BoolVar(&updateAll, "all", false, "Update all installed skills")
}
