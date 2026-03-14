package wizard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/saars/claude-code-statusline/components"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/settings"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	sectionStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	optNameStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	optDescStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	keyStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
	actionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	hintSep       = actionStyle.Render(" · ")

	// bar preview colors — use RGB so they're vivid regardless of terminal theme
	barGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#22DD55"))
	barYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDCC00"))
	barRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#DD3333"))
	barDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	barPct    = lipgloss.NewStyle().Foreground(lipgloss.Color("#22DD55"))

	// errNoFeatures is a sentinel returned by the features step when the user
	// selects nothing; Run treats this as a graceful cancellation.
	errNoFeatures = errors.New("no features selected")

	// errGoBack is a sentinel returned by a step when the user presses Escape
	// to navigate back to the previous step.
	errGoBack = errors.New("go back")
)

// hint renders a single "key → action" pair with distinct colors.
func hint(k, action string) string {
	return keyStyle.Render(k) + actionStyle.Render(" "+action)
}

// Run launches the interactive setup wizard. Pass empty strings for cfgPath
// and settingsPath to use the default locations.
func Run(cfgPath, settingsPath string) error {
	if cfgPath == "" {
		cfgPath = config.ConfigPath()
	}
	if settingsPath == "" {
		home, _ := os.UserHomeDir()
		settingsPath = filepath.Join(home, ".claude", "settings.json")
	}

	state := DefaultState()

	fmt.Println(headerStyle.Render("claude-code-statusline setup"))
	fmt.Println(hint("x/space", "toggle") + hintSep + hint("enter", "submit (not select!)") + hintSep + hint("esc", "back") + hintSep + hint("ctrl+c", "cancel"))
	fmt.Println()

	// ── Interactive steps ────────────────────────────────────────────────────

	for i := 0; i < len(Steps); {
		step := Steps[i]
		if step.ShouldRun != nil && !step.ShouldRun(state) {
			i++
			continue
		}
		if err := step.Run(state); err != nil {
			if errors.Is(err, errGoBack) {
				i = prevRunnableStep(Steps, i, state)
				continue
			}
			if errors.Is(err, errNoFeatures) {
				return nil
			}
			return err
		}
		i++
	}

	// ── Confirm ──────────────────────────────────────────────────────────────

	confirm := true
	if err := runWithPreview(huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("💾 Save this configuration?").
				Value(&confirm),
		),
	), state); err != nil {
		return err
	}
	if !confirm {
		fmt.Println(subtitleStyle.Render("Cancelled — no changes made."))
		return nil
	}

	// ── Write config ─────────────────────────────────────────────────────────

	tomlStr, err := state.ToTOML()
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	if err := os.WriteFile(cfgPath, []byte(tomlStr), 0o644); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	fmt.Printf("Config written to %s\n", cfgPath)

	// ── settings.json ────────────────────────────────────────────────────────

	binaryPath, _ := os.Executable()

	existing, _ := settings.Read(settingsPath)
	if existing != nil && existing.StatusLine != nil && existing.StatusLine.Command == binaryPath {
		fmt.Println(subtitleStyle.Render("settings.json already configured correctly — no changes needed."))
	} else {
		updateSettings := true
		if err := run(huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("⚙️  Update ~/.claude/settings.json automatically?").
					Description(fmt.Sprintf("Will set statusLine.command to %q", binaryPath)).
					Value(&updateSettings),
			),
		)); err != nil {
			return err
		}

		if updateSettings {
			if err := settings.WriteStatusLine(settingsPath, binaryPath); err != nil {
				return fmt.Errorf("could not update settings.json: %w", err)
			}
			fmt.Printf("settings.json updated at %s\n", settingsPath)
		} else {
			fmt.Printf("\nAdd this to ~/.claude/settings.json manually:\n")
			fmt.Printf("  \"statusLine\": {\n    \"type\": \"command\",\n    \"command\": %q\n  }\n\n", binaryPath)
		}
	}

	fmt.Println()
	fmt.Println(headerStyle.Render("Done!") + " " + subtitleStyle.Render("Restart Claude Code to see your status line."))
	return nil
}

// run executes a huh form with the Charm theme and converts ErrUserAborted
// (Ctrl+C) to a clean cancellation instead of an error.
func run(form *huh.Form) error {
	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Toggle = key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("x/space", "toggle"),
	)
	err := form.WithTheme(huh.ThemeCharm()).WithKeyMap(km).Run()
	if errors.Is(err, huh.ErrUserAborted) {
		fmt.Println(subtitleStyle.Render("\nSetup cancelled."))
		os.Exit(0)
	}
	return err
}

// prevRunnableStep returns the index of the previous step that would run
// given the current state. If there is no such step, it returns the current
// index (i.e. stay on the first step).
func prevRunnableStep(steps []Step, cur int, state *WizardState) int {
	for j := cur - 1; j >= 0; j-- {
		if steps[j].ShouldRun == nil || steps[j].ShouldRun(state) {
			return j
		}
	}
	return cur
}

// previewBlock returns the preview header rendered as a string for embedding
// in Bubble Tea views.
func previewBlock(state *WizardState) string {
	return "\n" + sectionStyle.Render("Preview") + "\n" +
		subtitleStyle.Render(strings.Repeat("─", 50)) + "\n" +
		Preview(state) + "\n" +
		subtitleStyle.Render(strings.Repeat("─", 50)) + "\n"
}

// previewFormModel wraps a huh.Form in a Bubble Tea model that renders a live
// status-line preview above the form. The preview reflects state on every frame.
type previewFormModel struct {
	form   *huh.Form
	state  *WizardState
	goBack bool
}

func (m previewFormModel) Init() tea.Cmd { return m.form.Init() }

func (m previewFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Let the form process the message first so that escape sequences
	// (arrow keys, colors, etc.) are consumed properly by Bubble Tea.
	f, cmd := m.form.Update(msg)
	m.form = f.(*huh.Form)
	if m.form.State != huh.StateNormal {
		return m, tea.Quit
	}
	// After the form has processed, treat a bare Escape as "go back".
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEsc {
		m.goBack = true
		return m, tea.Quit
	}
	return m, cmd
}

func (m previewFormModel) View() string {
	if m.goBack || m.form.State != huh.StateNormal {
		return ""
	}
	return previewBlock(m.state) + m.form.View()
}

// runWithPreview runs a huh form wrapped in a Bubble Tea program that shows a
// live status-line preview above the form on every render frame.
func runWithPreview(form *huh.Form, state *WizardState) error {
	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Toggle = key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("x/space", "toggle"),
	)
	form = form.WithTheme(huh.ThemeCharm()).WithKeyMap(km)
	p := tea.NewProgram(previewFormModel{form: form, state: state})
	final, err := p.Run()
	if err != nil {
		return err
	}
	fm := final.(previewFormModel)
	if fm.goBack {
		return errGoBack
	}
	if fm.form.State == huh.StateAborted {
		fmt.Println(subtitleStyle.Render("\nSetup cancelled."))
		os.Exit(0)
	}
	return nil
}

// opt renders a two-column option label: name in cyan, example in gray.
func opt(name, example string) string {
	return optNameStyle.Render(fmt.Sprintf("%-20s", name)) + optDescStyle.Render(example)
}

// barPreview renders a colored gradient bar for wizard option labels.
// Always uses █/░ since shade block chars (▓) ignore ANSI colors in many terminals.
func barPreview(filled, total int) string {
	greenEnd := int(0.70 * float64(total))
	yellowEnd := int(0.90 * float64(total))
	empty := total - filled

	gFill := components.Clamp(filled, 0, greenEnd)
	yFill := components.Clamp(filled-greenEnd, 0, yellowEnd-greenEnd)
	rFill := components.Clamp(filled-yellowEnd, 0, total-yellowEnd)

	return barGreen.Render(strings.Repeat("█", gFill)) +
		barYellow.Render(strings.Repeat("█", yFill)) +
		barRed.Render(strings.Repeat("█", rFill)) +
		barDim.Render(strings.Repeat("░", empty))
}

// styleOptions builds huh options for a feature's style selector by looking up
// component names from FeatureStyles/GetMeta. Examples are passed in because
// they're wizard-specific presentation (mock rendered output).
func styleOptions(feature string, examples map[string]string) []huh.Option[string] {
	styles := components.FeatureStyles[feature]
	opts := make([]huh.Option[string], len(styles))
	for i, s := range styles {
		m := components.GetMeta(s.ComponentKey)
		opts[i] = huh.NewOption(opt(m.Name, examples[s.Value]), s.Value)
	}
	return opts
}

func featureOptions() []huh.Option[string] {
	opts := make([]huh.Option[string], len(components.FeatureMeta))
	for i, f := range components.FeatureMeta {
		label := f.Meta.Emoji + " " + optNameStyle.Render(fmt.Sprintf("%-20s", f.Meta.Name)) + optDescStyle.Render(f.Meta.Desc)
		opts[i] = huh.NewOption(label, f.Key)
	}
	return opts
}
