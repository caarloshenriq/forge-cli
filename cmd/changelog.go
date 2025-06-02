package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a CHANGELOG.md based on git commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running changelog generator (coming soon!)")
	},
}
