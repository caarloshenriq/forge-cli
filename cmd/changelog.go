package cmd

import (
	"github.com/caarloshenriq/forge-cli/internal/changelog"
	"github.com/spf13/cobra"
)

var fromTag string
var toTag string

var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a CHANGELOG.md based on git commits",
	Run: func(cmd *cobra.Command, args []string) {
		changelog.GenerateChangelog(fromTag, toTag)
	},
}

func init() {
	ChangelogCmd.Flags().StringVar(&fromTag, "from", "", "Tag to start changelog from (optional)")
	ChangelogCmd.Flags().StringVar(&toTag, "to", "", "Tag to end changelog at (defaults to HEAD)")
}