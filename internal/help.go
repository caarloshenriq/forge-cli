package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelpInternal = &cobra.Command{
	Use:   "Help",
	Short: "Display help information for ForgeCLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running help (coming soon!)")
	},
}