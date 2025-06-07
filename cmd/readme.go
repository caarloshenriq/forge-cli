package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarloshenriq/forge-cli/internal"
	"github.com/spf13/cobra"
)

var ReadmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate a README.md with project info",
	Run: func(cmd *cobra.Command, args []string) {
		var choice string
		prompt := &survey.Select{
			Message: "What do you want to do?",
			Options: []string{
				"Generate from scratch",
				"Generate from template",
				"Back to main menu",
			},
		}

		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		switch choice {
		case "Generate from scratch":
			internal.GenerateReadmeFromScratch()
		case "Generate from template":
			internal.GenerateREADMEFromTemplate()
		case "Exit":
			return
		}
	},
}
