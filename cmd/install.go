package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholasspencer/gh-skill/internal"
	"github.com/spf13/cobra"
)

var installOutput string

var installCmd = &cobra.Command{
	Use:   "install <gist-url-or-id>",
	Short: "Download skill files to the current directory",
	Long:  "Downloads gist files directly without linking or managing. Use -o to specify an output directory.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gistID := internal.ParseGistID(args[0])

		fmt.Printf("Fetching gist %s...\n", gistID)
		gist, err := internal.FetchGist(gistID)
		if err != nil {
			return err
		}

		// Find skill file to determine name
		skillFileName, skillFile, ok := internal.FindSkillFile(gist.Files)
		if !ok {
			return fmt.Errorf("gist does not contain a *.skill.md file")
		}

		fm, _ := internal.ParseFrontMatter(skillFile.Content)
		name := fm.Name
		if name == "" {
			name = internal.SkillNameFromFile(skillFileName)
		}
		if name == "" {
			name = gist.ID
		}

		// Determine output directory
		outDir := installOutput
		if outDir == "" {
			outDir = "."
		}

		// Create skill subdirectory
		destDir := filepath.Join(outDir, name)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Write all files, expanding paths and renaming skill file
		fileCount := 0
		for filename, file := range gist.Files {
			expanded := internal.ExpandFilename(filename)
			if internal.IsSkillFile(expanded) {
				expanded = "SKILL.md"
			}
			destPath := filepath.Join(destDir, expanded)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", expanded, err)
			}
			if err := os.WriteFile(destPath, []byte(file.Content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", expanded, err)
			}
			fileCount++
		}

		fmt.Printf("✓ Installed %d files to %s/\n", fileCount, destDir)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(&installOutput, "output", "o", "", "Output directory (default: current directory)")
}
