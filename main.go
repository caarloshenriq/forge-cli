package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarloshenriq/forge-cli/cmd"
)

func main() {
	fmt.Println("ForgeCLI is a toolkit with useful commands for developers")
	fmt.Println("ForgeCLI combines tools like changelog generator, README generator, and more, into one CLI application.")

	for {
		var choice string
		prompt := &survey.Select{
			Message: "What do you want to do?",
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
			cmd.ChangelogCmd.Run(nil, nil)
		case "Generate README":
			cmd.ReadmeCmd.Run(nil, nil)
		case "Help":
      cmd.HelpInternal.Run(nil, nil)
    case "Exit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		}

		fmt.Println()
	}
}
