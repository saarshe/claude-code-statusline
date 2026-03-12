package wizard

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tickMsg is sent on each animation frame.
type tickMsg struct{}

type contextStyleChoice struct {
	label    string
	value    string
	renderEx func(pct float64) string // renders the example at given fill %
}

type contextStyleModel struct {
	choices []contextStyleChoice
	cursor  int
	pct     float64
	done    bool
	result  string
}

var (
	csTitle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	csCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
	csSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	csMuted    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	csName     = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

func newContextStyleModel(current string) contextStyleModel {
	choices := []contextStyleChoice{
		{
			label: "Percentage only",
			value: "pct",
			renderEx: func(pct float64) string {
				return thresholdColor(pct).Render(fmt.Sprintf("%.0f%%", pct))
			},
		},
		{
			label: "Token counts",
			value: "tokens",
			renderEx: func(pct float64) string {
				used := int(pct / 100 * 200)
				return thresholdColor(pct).Render(fmt.Sprintf("%dk / 200k", used))
			},
		},
		{
			label: "Tokens + bar",
			value: "tokens_bar",
			renderEx: func(pct float64) string {
				used := int(pct / 100 * 200)
				return csMuted.Render(fmt.Sprintf("%dk / 200k ", used)) +
					gradientBar(pct) + " " + thresholdColor(pct).Render(fmt.Sprintf("%.0f%%", pct))
			},
		},
		{
			label: "Block bar",
			value: "block",
			renderEx: func(pct float64) string {
				return blockBar(pct) + " " + thresholdColor(pct).Render(fmt.Sprintf("%.0f%%", pct))
			},
		},
		{
			label: "Gradient bar",
			value: "gradient",
			renderEx: func(pct float64) string {
				return gradientBar(pct) + " " + thresholdColor(pct).Render(fmt.Sprintf("%.0f%%", pct)) +
					"  " + csMuted.Render("(green→yellow→red)")
			},
		},
		{
			label: "Solid bar",
			value: "solid",
			renderEx: func(pct float64) string {
				return solidBar(pct) + " " + thresholdColor(pct).Render(fmt.Sprintf("%.0f%%", pct))
			},
		},
		{
			label: "ASCII bar",
			value: "ascii",
			renderEx: func(pct float64) string {
				filled := int(pct / 100 * 10)
				empty := 10 - filled
				bar := "[" + strings.Repeat("=", filled) + strings.Repeat("-", empty) + "]"
				return csMuted.Render(bar) + " " + barPct.Render(fmt.Sprintf("%.0f%%", pct))
			},
		},
	}

	cursor := 0
	for i, c := range choices {
		if c.value == current {
			cursor = i
			break
		}
	}

	return contextStyleModel{choices: choices, cursor: cursor, pct: 44.0}
}

// gradientBar renders a per-character gradient bar (green→yellow→red by position).
func gradientBar(pct float64) string {
	const total = 10
	filled := clamp(int(pct/100*total), 0, total)
	return barPreview(filled, total)
}

// solidBar renders a single-color bar using █, color set by threshold.
func solidBar(pct float64) string {
	const total = 10
	filled := clamp(int(pct/100*total), 0, total)
	empty := total - filled
	col := thresholdColor(pct)
	return col.Render(strings.Repeat("█", filled)) + barDim.Render(strings.Repeat("░", empty))
}

// blockBar renders a single-color bar using ▓, color set by threshold.
func blockBar(pct float64) string {
	const total = 10
	filled := clamp(int(pct/100*total), 0, total)
	empty := total - filled
	col := thresholdColor(pct)
	return col.Render(strings.Repeat("▓", filled)) + barDim.Render(strings.Repeat("░", empty))
}

func thresholdColor(pct float64) lipgloss.Style {
	switch {
	case pct >= 90:
		return barRed
	case pct >= 70:
		return barYellow
	default:
		return barGreen
	}
}

func (m contextStyleModel) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m contextStyleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.pct += 1.5
		if m.pct > 100 {
			m.pct = 0
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.result = m.choices[m.cursor].value
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

func (m contextStyleModel) View() string {
	var b strings.Builder

	b.WriteString("\n  " + csTitle.Render("📊 Context window — how verbose?") + "\n\n")

	for i, c := range m.choices {
		cursor := "  "
		nameStyle := csName
		if i == m.cursor {
			cursor = csCursor.Render("> ")
			nameStyle = csSelected
		}

		name := nameStyle.Render(fmt.Sprintf("%-16s", c.label))
		example := c.renderEx(m.pct)
		b.WriteString(fmt.Sprintf("  %s%s  %s\n", cursor, name, example))
	}

	b.WriteString("\n  " + csMuted.Render("↑/↓ navigate • enter select • ctrl+c cancel") + "\n")
	return b.String()
}

// runContextStyleSelector runs the animated context style selector and returns
// the chosen value. Returns current if the user cancels.
func runContextStyleSelector(current string) (string, error) {
	m := newContextStyleModel(current)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return current, err
	}
	result := final.(contextStyleModel)
	if result.result == "" {
		return current, nil
	}
	return result.result, nil
}
