package cmd

import (
	"github.com/caarloshenriq/forge-cli/internal"
	"github.com/spf13/cobra"
)

var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a CHANGELOG.md based on git commits",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GenerateChangelog()
	},
}
