// readme_tui.go
// Bubble Tea + Bubbles + Lipgloss README generator (full-screen)
// Modes: Template wizard and Scratch editor
// Features: Snippets palette (Ctrl+Space / Ctrl+@ / F2), Preview (Ctrl+P), Save (Ctrl+S)
// Scratch mode: on save, opens a modal to enter the file name (not forced to README.md)
// Language: English UI

package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/bubbles/help"
	key "github.com/charmbracelet/bubbles/key"
	list "github.com/charmbracelet/bubbles/list"
	textarea "github.com/charmbracelet/bubbles/textarea"
	textinput "github.com/charmbracelet/bubbles/textinput"

	"github.com/charmbracelet/lipgloss"
)

// ───────────────────────── Styles ─────────────────────────
var (
	styleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	styleOK     = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)

	styleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)

	styleBadge = lipgloss.NewStyle().
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1).MarginRight(1)
)

// ─────────────────────── Keymap ────────────────────────
type readmeKeymap struct {
	TogglePreview key.Binding
	Save          key.Binding
	Quit          key.Binding
	Snippets      key.Binding
	Next          key.Binding
	Prev          key.Binding
	Back          key.Binding
}

var keys = readmeKeymap{
	TogglePreview: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "preview")),
	Save:          key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
	Quit:          key.NewBinding(key.WithKeys("ctrl+q", "ctrl+c"), key.WithHelp("ctrl+q", "quit")),
	Snippets: key.NewBinding(
		key.WithKeys("ctrl+space", "ctrl+@", "f2"),
		key.WithHelp("ctrl+space/F2", "snippets"),
	),
	Next: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
	Prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
	Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

// ───────────────────── Snippets ────────────────────────
type snippet struct {
	Title       string
	Description string
	Template    string // cursor marker: "|"
}

func defaultSnippets() []snippet {
	return []snippet{
		{"H1", "Heading level 1", "# |"},
		{"H2", "Heading level 2", "## |"},
		{"Bold", "Strong emphasis", "**|**"},
		{"Italic", "Emphasis", "*|*"},
		{"Strike", "Strikethrough", "~~|~~"},
		{"Link", "[text](url)", "[|](https://example.com)"},
		{"Image", "![alt](url)", "![|](https://example.com/img.png)"},
		{"List", "Bullet list", "- |"},
		{"Numbered", "Ordered list", "1. |"},
		{"Checkbox", "Task item", "- [ ] |"},
		{"Code (inline)", "`code`", "`|`"},
		{"Code block (go)", "```go...```", "```go\n|\n```"},
		{"Quote", "> blockquote", "> |"},
		{"Table", "Simple table", "| A | B |\n|---|---|\n| | |\n| | |"},
	}
}

type snippetItem struct{ snippet }

func (s snippetItem) Title() string       { return s.snippet.Title }
func (s snippetItem) Description() string { return s.snippet.Description }
func (s snippetItem) FilterValue() string { return s.snippet.Title + " " + s.snippet.Description }

// ───────────────────── Modes ───────────────────────────
type mode int

const (
	modeMenu mode = iota
	modeTemplate
	modeScratch
	modeSnippets
)

// ─────────────── Template Data/const ────────────────
type ReadmeTemplateData struct {
	ProjectName  string
	Description  string
	Installation string
	Usage        string
	License      string
}

const readmeTmpl = `# {{.ProjectName}}

{{.Description}}

## Installation

{{.Installation}}

## Usage

{{.Usage}}

## License

{{.License}}
`

// ───────────────────── Model ──────────────────────────
type readmeModel struct {
	width, height int

	mode           mode
	lastMode       mode
	previewEnabled bool

	// menu
	menu list.Model

	// template
	nameInput textinput.Model
	descArea  textarea.Model
	instArea  textarea.Model
	usageArea textarea.Model
	licInput  textinput.Model
	focusIdx  int // 0..4

	// scratch
	editor textarea.Model

	// save modal (scratch)
	saveModalOpen bool
	saveInput     textinput.Model
	saveError     string

	// snippets
	snippetList list.Model

	// feedback
	saveMessage  string
	errorMessage string

	help help.Model
}

// Preview helper (fills the template for the side preview)
func renderTemplatePreview(d ReadmeTemplateData) string {
	var buf bytes.Buffer
	tmpl, err := template.New("preview").Parse(readmeTmpl)
	if err != nil {
		return fmt.Sprintf("template parse error: %v", err)
	}
	if err := tmpl.Execute(&buf, d); err != nil {
		return fmt.Sprintf("template exec error: %v", err)
	}
	return buf.String()
}

func ReadmeInitialModel() readmeModel {
	// menu
	items := []list.Item{
		snippetItem{snippet{Title: "Generate from Template", Description: "Fill in fields and generate README.md"}},
		snippetItem{snippet{Title: "Write from Scratch", Description: "Markdown editor with snippets"}},
	}
	menu := list.New(items, list.NewDefaultDelegate(), 0, 8)
	menu.Title = "README Generator"
	menu.SetShowPagination(false)
	menu.SetShowStatusBar(false)
	menu.SetShowHelp(false)

	// template widgets
	name := textinput.New()
	name.Placeholder = "Project name"
	name.Focus()

	desc := textarea.New()
	desc.Placeholder = "Project description…"
	desc.SetHeight(6)

	inst := textarea.New()
	inst.Placeholder = "Installation steps…"
	inst.SetHeight(6)

	usage := textarea.New()
	usage.Placeholder = "Usage examples…"
	usage.SetHeight(6)

	lic := textinput.New()
	lic.Placeholder = "MIT"

	// scratch editor (full-screen)
	ed := textarea.New()
	ed.Placeholder = "Start typing your Markdown… (Snippets: Ctrl+Space/F2)"
	ed.ShowLineNumbers = true

	// save modal input
	saveInput := textinput.New()
	saveInput.Placeholder = "filename (e.g., README.md)"

	// snippets palette
	var sitems []list.Item
	for _, s := range defaultSnippets() {
		sitems = append(sitems, snippetItem{s})
	}
	slist := list.New(sitems, list.NewDefaultDelegate(), 0, 12)
	slist.Title = "Snippets"
	slist.SetFilteringEnabled(true)
	slist.SetShowPagination(false)
	slist.SetShowStatusBar(false)
	slist.SetShowHelp(false)

	return readmeModel{
		mode:        modeMenu,
		menu:        menu,
		nameInput:   name,
		descArea:    desc,
		instArea:    inst,
		usageArea:   usage,
		licInput:    lic,
		editor:      ed,
		saveInput:   saveInput,
		snippetList: slist,
		help:        help.New(),
	}
}

// ───────────────────── Init/Update ───────────────────
func (m readmeModel) Init() tea.Cmd { return nil }

func (m readmeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Full-screen sizing
		pad := 4
		contentW := m.width - pad
		contentH := m.height - 6 // header + footer room
		if contentW < 40 {
			contentW = m.width
		}
		if contentH < 10 {
			contentH = m.height - 2
		}
		m.menu.SetSize(contentW-4, contentH)
		m.snippetList.SetSize(contentW-4, contentH)

		m.nameInput.Width = contentW - 6
		m.licInput.Width = contentW - 6
		for _, ta := range []*textarea.Model{&m.descArea, &m.instArea, &m.usageArea, &m.editor} {
			ta.SetWidth(contentW - 6)
		}
		// editor height uses most of screen
		m.editor.SetHeight(contentH)
		return m, nil

	case tea.KeyMsg:
		// global keys (except when modal is open)
		if m.saveModalOpen {
			// handle modal keys only
			switch msg.Type {
			case tea.KeyEnter:
				name := strings.TrimSpace(m.saveInput.Value())
				if name == "" {
					m.saveError = "Filename is required"
					return m, nil
				}
				if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil && filepath.Dir(name) != "." {
					m.saveError = err.Error()
					return m, nil
				}
				if err := os.WriteFile(name, []byte(m.editor.Value()), 0o644); err != nil {
					m.saveError = err.Error()
					return m, nil
				}
				m.saveMessage = fmt.Sprintf("Saved %s", name)
				m.saveError = ""
				m.saveModalOpen = false
				return m, nil
			case tea.KeyEsc:
				m.saveModalOpen = false
				m.saveError = ""
				return m, nil
			}
			// propagate modal input updates
			var cmd tea.Cmd
			m.saveInput, cmd = m.saveInput.Update(msg)
			return m, cmd
		}

		// not in modal
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.TogglePreview):
			m.previewEnabled = !m.previewEnabled
			return m, nil
		case key.Matches(msg, keys.Snippets):
			if m.mode == modeScratch || m.mode == modeTemplate {
				m.lastMode = m.mode
				m.mode = modeSnippets
				m.snippetList.ResetFilter()
				return m, nil
			}
		}
	}

	// state-specific
	switch m.mode {
	case modeMenu:
		var cmd tea.Cmd
		m.menu, cmd = m.menu.Update(msg)
		if km, ok := msg.(tea.KeyMsg); ok && km.Type == tea.KeyEnter {
			if m.menu.Index() == 0 {
				m.mode = modeTemplate
				m.setTemplateFocus(0)
			} else {
				m.mode = modeScratch
				m.editor.Focus()
			}
		}
		return m, cmd

	case modeSnippets:
		var cmd tea.Cmd
		m.snippetList, cmd = m.snippetList.Update(msg)
		if km, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(km, keys.Back):
				m.mode = m.lastMode
			case km.Type == tea.KeyEnter:
				if it, ok := m.snippetList.SelectedItem().(snippetItem); ok {
					m.insertSnippet(it.snippet)
				}
				m.mode = m.lastMode
			}
		}
		return m, cmd

	case modeTemplate:
		if k, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(k, keys.Next):
				m.setTemplateFocus(m.focusIdx + 1)
				return m, nil
			case key.Matches(k, keys.Prev):
				m.setTemplateFocus(m.focusIdx - 1)
				return m, nil
			case key.Matches(k, keys.Save):
				if err := m.saveTemplate(); err != nil {
					m.errorMessage = err.Error()
					m.saveMessage = ""
				} else {
					m.errorMessage = ""
					m.saveMessage = "Saved README.md"
				}
				return m, nil
			}
		}
		m.nameInput, _ = m.nameInput.Update(msg)
		m.licInput, _ = m.licInput.Update(msg)
		m.descArea, _ = m.descArea.Update(msg)
		m.instArea, _ = m.instArea.Update(msg)
		m.usageArea, _ = m.usageArea.Update(msg)
		return m, nil

	case modeScratch:
		if k, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(k, keys.Save):
				// open modal
				m.saveModalOpen = true
				m.saveError = ""
				m.saveInput.SetValue("README.md")
				m.saveInput.Focus()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	}

	return m, nil
}

// ───────────────────── View ───────────────────────────
func (m readmeModel) View() string {
	title := styleTitle.Render("README Generator") + "\n"
	help := styleSubtle.Render("↑/↓ select • enter confirm • ctrl+q quit • ctrl+space/F2 snippets • ctrl+s save")

	switch m.mode {
	case modeMenu:
		content := styleCard.Width(m.width - 4).Render(m.menu.View())
		return title + content + "\n" + help

	case modeSnippets:
		// opcional: título próprio para o palette
		snipsTitle := styleTitle.Render("Snippets") + "\n"
		content := styleCard.Width(m.width - 4).Render(m.snippetList.View())
		return snipsTitle + content + "\n" + styleSubtle.Render("enter insert • esc back • / filter")

	case modeTemplate:
		left := strings.Join([]string{
			styleBadge.Render("Name") + " " + m.nameInput.View(),
			styleBadge.Render("Description") + "\n" + m.descArea.View(),
			styleBadge.Render("Installation") + "\n" + m.instArea.View(),
			styleBadge.Render("Usage") + "\n" + m.usageArea.View(),
			styleBadge.Render("License") + " " + m.licInput.View(),
		}, "\n\n")

		content := styleCard.Width(m.width - 4).Render(left)
		if m.previewEnabled {
			preview := styleCard.
				Width(m.width/2 - 4).
				Render(renderTemplatePreview(ReadmeTemplateData{
					ProjectName:  m.nameInput.Value(),
					Description:  m.descArea.Value(),
					Installation: m.instArea.Value(),
					Usage:        m.usageArea.Value(),
					License:      defaultLicense(m.licInput.Value()),
				}))
			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				styleCard.Width(m.width/2).Render(left),
				preview,
			)
		}

		status := ""
		if m.errorMessage != "" {
			status = "\n" + styleErr.Render(m.errorMessage)
		} else if m.saveMessage != "" {
			status = "\n" + styleOK.Render(m.saveMessage)
		}

		// título no topo; help no rodapé
		return title + content + status + "\n" + styleSubtle.Render("[tab/shift+tab] move • [ctrl+space/F2] snippets • [ctrl+s] save • [ctrl+p] preview")

	case modeScratch:
		ed := styleCard.Width(m.width - 4).Height(m.height - 8).Render(m.editor.View())

		status := ""
		if m.errorMessage != "" {
			status = "\n" + styleErr.Render(m.errorMessage)
		} else if m.saveMessage != "" {
			status = "\n" + styleOK.Render(m.saveMessage)
		}

		view := title + ed + status + "\n" + help

		// quando o modal está aberto, mostramos só o modal (overlay simples)
		if m.saveModalOpen {
			view = m.renderSaveModal(view)
		}
		return view
	}
	return ""
}


func (m readmeModel) renderSaveModal(background string) string {
	boxW := m.width / 2
	if boxW < 40 {
		boxW = 40
	}
	if boxW > m.width-4 {
		boxW = m.width - 4
	}
	lines := []string{
		styleTitle.Render("Save As"),
		"Enter the filename to save your document:",
		"",
		m.saveInput.View(),
	}
	if m.saveError != "" {
		lines = append(lines, "", styleErr.Render(m.saveError))
	}
	lines = append(lines, "", styleSubtle.Render("Enter to save • Esc to cancel"))
	modal := styleCard.Width(boxW).Render(strings.Join(lines, "\n"))
	return modal
}

// ───────────── Focus Helpers ─────────────
func (m *readmeModel) setTemplateFocus(idx int) {
	idx = (idx%5 + 5) % 5
	m.focusIdx = idx
	m.nameInput.Blur()
	m.descArea.Blur()
	m.instArea.Blur()
	m.usageArea.Blur()
	m.licInput.Blur()
	switch idx {
	case 0:
		m.nameInput.Focus()
	case 1:
		m.descArea.Focus()
	case 2:
		m.instArea.Focus()
	case 3:
		m.usageArea.Focus()
	case 4:
		m.licInput.Focus()
	}
}

// ───────────── Snippet insertion ─────────────
func moveToMarkerAfterInsert(area *textarea.Model, tmpl string) {
	idx := strings.IndexRune(tmpl, '|')
	if idx < 0 {
		return
	}

	runes := []rune(tmpl)
	line, col, total := 0, 0, 1
	for i, r := range runes {
		if r == '\n' {
			total++
		}
		if i < idx {
			if r == '\n' {
				line++
				col = 0
			} else {
				col++
			}
		}
	}
	linesAfter := total - 1 - line
	for i := 0; i < linesAfter; i++ {
		area.CursorUp()
	}
	li := area.LineInfo()
	area.SetCursor(li.StartColumn + col)
}

func (m *readmeModel) insertSnippet(s snippet) {
	body := strings.ReplaceAll(s.Template, "|", "")
	switch {
	case m.descArea.Focused():
		m.descArea.InsertString(body)
		moveToMarkerAfterInsert(&m.descArea, s.Template)
	case m.instArea.Focused():
		m.instArea.InsertString(body)
		moveToMarkerAfterInsert(&m.instArea, s.Template)
	case m.usageArea.Focused():
		m.usageArea.InsertString(body)
		moveToMarkerAfterInsert(&m.usageArea, s.Template)
	default:
		m.editor.InsertString(body)
		moveToMarkerAfterInsert(&m.editor, s.Template)
	}
}

// ───────────── Saving ─────────────
func (m *readmeModel) saveTemplate() error {
	data := ReadmeTemplateData{
		ProjectName:  strings.TrimSpace(m.nameInput.Value()),
		Description:  strings.TrimSpace(m.descArea.Value()),
		Installation: strings.TrimSpace(m.instArea.Value()),
		Usage:        strings.TrimSpace(m.usageArea.Value()),
		License:      strings.TrimSpace(m.licInput.Value()),
	}
	if data.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if data.License == "" {
		data.License = "MIT"
	}

	tmpl, err := template.New("readme").Parse(readmeTmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	return os.WriteFile("README.md", buf.Bytes(), 0o644)
}

// ───────────── Entry ─────────────
func RunReadmeTUI() error {
	m := ReadmeInitialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func defaultLicense(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "MIT"
	}
	return v
}
