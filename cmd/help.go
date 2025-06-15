package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/caarloshenriq/forge-cli/utils"
	"github.com/spf13/cobra"
)

var HelpInternal = &cobra.Command{
	Use:   "help",
	Short: "Display help information for ForgeCLI",
	Run: func(cmd *cobra.Command, args []string) {
		utils.ClearScreen()
		showHelpIndex()
	},
}

func showHelpIndex() {
	utils.ClearScreen()
	fmt.Println("ForgeCLI - Help")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("ForgeCLI is an open-source developer toolkit CLI.")
	fmt.Println("It provides useful commands to automate repetitive tasks,")
	fmt.Println("speed up project setup, and improve workflow efficiency.")
	fmt.Println()
	fmt.Println("📦 Current features:")
	fmt.Println("• Changelog Generator")
	fmt.Println("• README Generator")
	fmt.Println()
	fmt.Println("🛠️ Upcoming features:")
	fmt.Println("• .gitignore file generator")
	fmt.Println("• Project templates (Go, Node.js, PHP, and more)")
	fmt.Println()
	fmt.Println("Use the arrow keys in interactive mode to navigate between help topics.")
	fmt.Println("-----------------------------------------------------")

	for {
		var choice string
		prompt := &survey.Select{
			Message: "Which topic would you like help with?",
			Options: []string{
				"Changelog Generator",
				"README Generator",
				"How to Contribute",
				"Back to main menu",
			},
		}

		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		switch choice {
		case "Changelog Generator":
			helpChangelog()
		case "README Generator":
			helpReadme()
		case "How to Contribute":
			helpContribute()
		case "Back to main menu":
			return
		}

		fmt.Println()
	}
}

func helpChangelog() {
	utils.ClearScreen()
	fmt.Println("📝 Changelog Generator - Help")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("This command generates or updates your CHANGELOG.md file based on your Git commit history.")
	fmt.Println()
	fmt.Println("🔹 Usage:")
	fmt.Println("  forge-cli changelog")
	fmt.Println()
	fmt.Println("🔹 How it works:")
	fmt.Println("  - Only commits with prefixes 'feat:' and 'fix:' are categorized as Features and Fixes.")
	fmt.Println("  - Other commits will appear under 'Other commits'.")
	fmt.Println("  - You'll be able to manually select which commits to include in each section.")
	fmt.Println("  - Commits already registered in a previous changelog will be skipped automatically.")
	fmt.Println("  - The latest commit hash is tracked via:")
	fmt.Println("    <!-- last-commit: abc123def456 -->")
	fmt.Println()
	fmt.Println("✅ Tip:")
	fmt.Println("  Make sure all your changes are committed before running this command.")
}

func helpReadme() {
	utils.ClearScreen()
	fmt.Println("📄 README Generator - Help")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("This command helps you create a structured README.md for your project.")
	fmt.Println()
	fmt.Println("🔹 Usage:")
	fmt.Println("  forge-cli readme")
	fmt.Println()
	fmt.Println("🔹 Options:")
	fmt.Println("  • Generate from template: Answer guided questions to auto-generate a standard README.")
	fmt.Println("  • Generate from scratch: Write your own content using a multi-line editor.")
	fmt.Println()
	fmt.Println("💡 Notes:")
	fmt.Println("  - When writing from scratch, end the input with ':done' to finish.")
	fmt.Println("  - For license selection, if left blank, it defaults to 'MIT'.")
	fmt.Println()
	fmt.Println("✅ Tip:")
	fmt.Println("  Keep your README clear and concise. This helps contributors and users understand your project.")
}

func helpContribute() {
	utils.ClearScreen()
	fmt.Println("🤝 How to Contribute - Help")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("We welcome contributions from the community! To get started:")
	fmt.Println()
	fmt.Println("🔹 Fork the repository on GitHub")
	fmt.Println("🔹 Clone your fork and create a new branch")
	fmt.Println("🔹 Make your changes and follow the commit message guidelines (use feat:, fix:, refactor:, etc.)")
	fmt.Println("🔹 Test your code locally")
	fmt.Println("🔹 Submit a Pull Request (PR) with a clear description of what was changed and why")
	fmt.Println()
	fmt.Println("📄 You can also contribute by:")
	fmt.Println("  - Creating issues for bugs or feature suggestions")
	fmt.Println("  - Improving documentation")
	fmt.Println("  - Helping review pull requests")
	fmt.Println()
	fmt.Println("✅ Tip:")
	fmt.Println("  Always sync your fork with the original repo before starting new changes.")
}
