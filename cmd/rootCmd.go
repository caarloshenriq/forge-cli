package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge-cli",
	Short: "ForgeCLI is a toolkit with useful commands for developers",
	Long:  "ForgeCLI is a toolkit with useful commands for developers. It combines tools like changelog generator, README generator, and more into one CLI application.",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			clearScreen()

			var choice string
			prompt := &survey.Select{
				Message: "What do you want to do? (Press ESC to enable vim mode)",
				Options: []string{
					"Generate Changelog",
					"Generate README",
					"Help",
					"Exit",
				},
			}

			err := survey.AskOne(prompt, &choice)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			switch choice {
			case "Generate Changelog":
				ChangelogCmd.Run(nil, nil)
			case "Generate README":
				ReadmeCmd.Run(nil, nil)
			case "Help":
				HelpInternal.Run(nil, nil)
			case "Exit":
				fmt.Println("Goodbye!")
				os.Exit(0)
			}

			fmt.Println()
		}
	},
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ChangelogCmd)
	rootCmd.AddCommand(ReadmeCmd)
	rootCmd.AddCommand(HelpInternal)
}
