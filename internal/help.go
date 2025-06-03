package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HelpInternal = &cobra.Command{
	Use:   "help",
	Short: "Display help information for ForgeCLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ForgeCLI - Help")
		fmt.Println("-----------------------------------------------------")
		fmt.Println("ForgeCLI is an open-source developer toolkit CLI.")
		fmt.Println("It provides useful commands to automate repetitive tasks,")
		fmt.Println("speed up project setup, and improve workflow efficiency.")
		fmt.Println()
		fmt.Println("Current features include:")
		fmt.Println("- Changelog Generator: Automatically generate a CHANGELOG.md")
		fmt.Println("  file based on your git commit history.")
		fmt.Println("- README Generator: Create a custom README.md file by answering")
		fmt.Println("  interactive prompts about your project.")
		fmt.Println()
		fmt.Println("Upcoming features will include:")
		fmt.Println("- .gitignore file generator")
		// fmt.Println("- GitHub/GitLab integration for release creation")
		fmt.Println("- Project templates for Go, Node.js, PHP, and more")
		// fmt.Println("- Automated license and contributing file generation")
		fmt.Println()
		fmt.Println("Use the arrow keys in interactive mode to navigate between options.")
		fmt.Println("To run a command directly, use: forge-cli [command]")
		fmt.Println("-----------------------------------------------------")
	},
}
