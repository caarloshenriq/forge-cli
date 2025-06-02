package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
  "github.com/caarloshenriq/forge-cli/cmd"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "forge-cli",
		Short: "ForgeCLI is a toolkit with useful commands for developers",
		Long:  "ForgeCLI combines tools like changelog generator, README generator, and more, into one CLI application.",
	}

	// Add commands here
	rootCmd.AddCommand(cmd.ChangelogCmd)
	rootCmd.AddCommand(cmd.ReadmeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
