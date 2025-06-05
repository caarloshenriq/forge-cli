package cmd

import (
	"github.com/caarloshenriq/forge-cli/internal"
	"github.com/spf13/cobra"
)

var ReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate a README.md with project info",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GenerateREADME()
	},
}
