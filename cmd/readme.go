package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate a README.md with project info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running README generator (coming soon!)")
	},
}
