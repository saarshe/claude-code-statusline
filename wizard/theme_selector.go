package wizard

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/saarshe/claude-code-statusline/theme"
)

// sortedThemeNames returns theme names sorted alphabetically with "default" first.
func sortedThemeNames() []string {
	names := theme.Names()
	sort.Strings(names)
	for i, n := range names {
		if n == "default" && i != 0 {
			names[0], names[i] = names[i], names[0]
			break
		}
	}
	return names
}

type themeSelectorModel struct {
	names  []string
	cursor int
	done   bool
	goBack bool
	result string
	state  *WizardState
}

func newThemeSelectorModel(state *WizardState) themeSelectorModel {
	names := sortedThemeNames()
	cursor := 0
	for i, n := range names {
		if n == state.Theme {
			cursor = i
			break
		}
	}
	return themeSelectorModel{names: names, cursor: cursor, state: state}
}

func (m themeSelectorModel) Init() tea.Cmd { return nil }

func (m themeSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.names)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.result = m.names[m.cursor]
			m.done = true
			return m, tea.Quit
		case tea.KeyEsc:
			m.goBack = true
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.result = ""
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m themeSelectorModel) View() string {
	if m.done {
		return ""
	}

	// Build a preview state using the hovered theme.
	hovered := *m.state
	hovered.Theme = m.names[m.cursor]

	var b strings.Builder

	b.WriteString(previewBlock(&hovered))
	b.WriteString("\n  " + csTitle.Render("🎨 Choose a color theme") + "\n\n")

	for i, name := range m.names {
		cursor := "  "
		nameStyle := csName
		if i == m.cursor {
			cursor = csCursor.Render("> ")
			nameStyle = csSelected
		}
		b.WriteString(fmt.Sprintf("  %s%s\n", cursor, nameStyle.Render(name)))
	}

	b.WriteString(selectorHints())
	return b.String()
}

// runThemeSelector runs the interactive theme selector and writes the chosen
// theme into state.Theme.
func runThemeSelector(state *WizardState) error {
	m := newThemeSelectorModel(state)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return err
	}
	result := final.(themeSelectorModel)
	if result.goBack {
		return errGoBack
	}
	if result.result == "" {
		// ctrl+c — treat as cancel
		fmt.Println(subtitleStyle.Render("\nSetup cancelled."))
		os.Exit(0)
	}
	state.Theme = result.result
	return nil
}
