package wizard

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/render"
	"github.com/saars/claude-code-statusline/theme"
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

// themePreviewCfg is a compact config used to render theme previews.
var themePreviewCfg = &config.Config{
	Emojis: config.EmojiAll,
	ContextBar: config.ContextBarConfig{
		Style:      config.BarSolid,
		Width:      10,
		Thresholds: []int{70, 90},
	},
	Separator: config.SeparatorConfig{Character: "|"},
	Lines: []config.LineConfig{
		{Components: []string{"model", "context_bar", "cost"}},
	},
}

// themePreview renders a sample status line using the given theme.
func themePreview(name string) string {
	th := theme.Get(name)
	cfg := *themePreviewCfg
	cfg.Theme = name
	return render.RenderWithTheme(MockInput(), &cfg, th)
}

type themeSelectorModel struct {
	names  []string
	cursor int
	done   bool
	result string
}

func newThemeSelectorModel(current string) themeSelectorModel {
	names := sortedThemeNames()
	cursor := 0
	for i, n := range names {
		if n == current {
			cursor = i
			break
		}
	}
	return themeSelectorModel{names: names, cursor: cursor}
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

	var b strings.Builder

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

	// Live preview of the hovered theme.
	b.WriteString("\n  " + csMuted.Render("preview: ") + themePreview(m.names[m.cursor]) + "\n")
	b.WriteString("\n  " + csMuted.Render("↑/↓ navigate • enter select • ctrl+c cancel") + "\n")

	return b.String()
}

// runThemeSelector runs the interactive theme selector and returns the chosen
// theme name. Returns current if the user cancels.
func runThemeSelector(current string) (string, error) {
	m := newThemeSelectorModel(current)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return current, err
	}
	result := final.(themeSelectorModel)
	if result.result == "" {
		return current, nil
	}
	return result.result, nil
}
