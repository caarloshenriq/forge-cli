package cmd

import (
	"fmt"
	"os"
	"strings"

	// Removed import of internal package to avoid cycle
	"github.com/caarloshenriq/forge-cli/internal"
	"github.com/caarloshenriq/forge-cli/utils"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type helpModel struct {
	width, height int
	cursor        int
	options       []string
	selected      string
	quitting      bool
	mode          string // "" | "viewing"
}


var (
	titleStyle = gloss.NewStyle().
			Bold(true).
			Foreground(gloss.Color("#00afff")).
			PaddingLeft(2).
			PaddingRight(1).
			Align(gloss.Center)

	cursorStyle = gloss.NewStyle().
			Foreground(gloss.Color("#00ff87")).
			Bold(true)

	selectedStyle = gloss.NewStyle().
			Bold(true).
			Foreground(gloss.Color("#ffffff")).
			Background(gloss.Color("#5c5c5c")).
			PaddingLeft(1).
			PaddingRight(1)

	itemStyle = gloss.NewStyle().
			Foreground(gloss.Color("#cccccc"))
)

func initialHelpModel() helpModel {
	return helpModel{
		options: []string{
			"Changelog Generator",
			"README Generator",
			"How to Contribute",
			"Back to main menu",
		},
	}
}

func (m helpModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m helpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "b":
			if m.mode == "viewing" {
				m.mode = ""
				m.selected = ""
				m.cursor = 0
				return m, nil
			}
		case "enter":
			m.selected = m.options[m.cursor]
			m.mode = "viewing"
			return m, nil
		}
	}
	return m, nil
}

func (m helpModel) View() string {
	if m.mode == "viewing" && m.selected != "" {
		return renderHelpTopic(m.selected) + gloss.NewStyle().Italic(true).PaddingTop(3).Render("Press [b] to go back")
	}

	var lines []string
	lines = append(lines, centerText(m.width, titleStyle.Render("ForgeCLI - Help")))
	lines = append(lines, "")

	for i, opt := range m.options {
		if i == m.cursor {
			lines = append(lines, centerText(m.width, cursorStyle.Render("› ")+selectedStyle.Render(opt)))
		} else {
			lines = append(lines, centerText(m.width, "  "+itemStyle.Render(opt)))
		}
	}
	lines = append(lines, centerText(m.width, gloss.NewStyle().Italic(true).PaddingTop(3).Render("Use ↑ ↓ to navigate and press Enter to view a topic.")))
	return strings.Join(lines, "\n")
}


func centerText(width int, text string) string {
	padding := (width - gloss.Width(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

func renderHelpTopic(topic string) string {
	utils.ClearScreen()
	switch topic {
	case "Changelog Generator":
		return helpChangelog()
	case "README Generator":
		return helpReadme()
	case "How to Contribute":
		return helpContribute()
	case "Back to main menu":
		return internal.RunTUI()
	default:
		return "\nReturning to main menu...\n"
	}
}

func helpChangelog() string {
	return `
📝 Changelog Generator - Help
-----------------------------------------------------
This command generates or updates your CHANGELOG.md file based on your Git commit history.

🔹 Usage:
  forge-cli changelog

🔹 How it works:
  - Only commits with prefixes 'feat:' and 'fix:' are categorized.
  - Others go under 'Other commits'.
  - Already logged commits are skipped.
  - Tracks last commit via: <!-- last-commit: abc123 -->

✅ Tip:
  Commit all changes before generating the changelog.
`
}

func helpReadme() string {
	return `
📄 README Generator - Help
-----------------------------------------------------
This command helps create a structured README.md for your project.

🔹 Usage:
  forge-cli readme

🔹 Options:
  • Generate from template (guided)
  • Write from scratch (multi-line)

💡 Notes:
  - End input with ':done' when writing manually.
  - License defaults to MIT if blank.

✅ Tip:
  Keep README clear and concise.
`
}

func helpContribute() string {
	return `
🤝 How to Contribute - Help
-----------------------------------------------------
We welcome community contributions!

🔹 Start by:
  - Forking on GitHub
  - Creating a new branch
  - Following commit message conventions (feat:, fix:, etc.)
  - Submitting a PR with a clear description

📄 Other ways to contribute:
  - Report issues
  - Improve documentation
  - Review pull requests

✅ Tip:
  Sync your fork regularly with upstream before new changes.
`
}

func RunHelpTUI() {
	p := tea.NewProgram(initialHelpModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running help view: %v\n", err)
		os.Exit(1)
	}
}
