package internal

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("#00afff")).
		Padding(1, 0).
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

	helpStyle = gloss.NewStyle().
		Foreground(gloss.Color("#888888")).
		Italic(true).
		MarginTop(1)
)

type model struct {
	width, height   int
	selected, cursor int
	choices         []string
	quitting        bool
	selectedValue   string
}

func initialModel() model {
	return model{
		choices: []string{
			"Generate Changelog",
			"Generate README",
			"Help",
		},
		selected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selectedValue = m.choices[m.cursor]
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		if m.selectedValue != "" {
			return fmt.Sprintf("\n  You selected: %s\n\n", m.selectedValue)
		}
		return "\n  Goodbye!\n\n"
	}

	var lines []string

	// Título
	lines = append(lines, centerText(m.width, titleStyle.Render("Forge CLI")))
	lines = append(lines, "")

	// Menu
	for i, choice := range m.choices {
		var line string
		if i == m.cursor {
			line = cursorStyle.Render("› ") + selectedStyle.Render(choice)
		} else {
			line = "  " + itemStyle.Render(choice)
		}
		lines = append(lines, centerText(m.width, line))
	}

	// Ajuda
	lines = append(lines, "")
	lines = append(lines, centerText(m.width, helpStyle.Render("↑/↓: navigate • enter: select • q: quit")))

	return strings.Join(lines, "\n")
}

func centerText(width int, text string) string {
	padding := (width - gloss.Width(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

func RunTUI() string {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running TUI: %v\n", err)
		os.Exit(1)
	}

	m, ok := finalModel.(model)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected model type\n")
		os.Exit(1)
	}

	return m.selectedValue
}
