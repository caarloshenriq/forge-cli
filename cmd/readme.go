package cmd

import (
	"fmt"
	"os"

	"github.com/caarloshenriq/forge-cli/internal"
	"github.com/caarloshenriq/forge-cli/utils"
	"github.com/spf13/cobra"
)

func runReadmeTUI() error {
	utils.ClearScreen()
	return internal.RunReadmeTUI()
}

var ReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate a README.md using the interactive TUI (template or scratch, with snippets)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runReadmeTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "README TUI error: %v\n", err)
			return err
		}
		return nil
	},
}
