package wizard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/settings"
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	sectionStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	optNameStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	optDescStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	keyStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))  // yellow
	actionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))             // gray
	hintSep       = actionStyle.Render(" · ")
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

	// ── Step 1: What data do you want to see? ─────────────────────────────────

	fmt.Println(headerStyle.Render("claude-code-statusline setup"))
	fmt.Println(hint("x/space", "toggle") + hintSep + hint("enter", "submit (not select!)") + hintSep + hint("ctrl+c", "cancel"))
	fmt.Println()

	selected := state.Features
	if err := run(huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("What data do you want to see?").
				Options(featureOptions()...).
				Value(&selected),
		),
	)); err != nil {
		return err
	}
	state.Features = selected

	if len(state.Features) == 0 {
		fmt.Println(subtitleStyle.Render("No features selected — setup cancelled."))
		return nil
	}

	// ── Step 2: Context window style (conditional) ────────────────────────────

	if state.HasContext() {
		if err := run(huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("📊 Context window — how verbose?").
					Options(
						huh.NewOption(opt("Percentage only", "44%"), "pct"),
						huh.NewOption(opt("Token counts", "88k / 200k"), "tokens"),
						huh.NewOption(opt("Tokens + bar", "88k / 200k ▓▓▓▓░░░░░░"), "tokens_bar"),
						huh.NewOption(opt("Block bar", "▓▓▓▓░░░░░░ 44%"), "block"),
						huh.NewOption(opt("Gradient bar", "▓▓▓▓▓▓░░░░ 44%  (green→yellow→red zones)"), "gradient"),
						huh.NewOption(opt("Solid bar", "████░░░░░░ 44%"), "solid"),
						huh.NewOption(opt("ASCII bar", "[====------] 44%"), "ascii"),
					).
					Value(&state.ContextStyle),
			),
		)); err != nil {
			return err
		}
	}

	// ── Step 3: Cache style (conditional) ─────────────────────────────────────

	if state.HasCache() {
		if err := run(huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("💾 Cache — how verbose?").
					Description(
						"Each turn, Claude reuses previously processed context from cache (cheap)\n" +
							"and stores new context for the next turn to reuse.\n",
					).
					Options(
						huh.NewOption(opt("Efficiency", "⚡ 37% cached"), "hit"),
						huh.NewOption(opt("Counts", "💾 5.0k reused, 2.0k stored"), "counts"),
					).
					Value(&state.CacheStyle),
			),
		)); err != nil {
			return err
		}
	}

	// ── Step 4: Git style (conditional) ──────────────────────────────────────

	if state.HasGit() {
		if err := run(huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("🌿 Git — how verbose?").
					Options(
						huh.NewOption(opt("Branch only", "🌿 main"), "branch"),
						huh.NewOption(opt("Branch + changes", "🌿 main +1 ~9"), "status"),
					).
					Value(&state.GitStyle),
			),
		)); err != nil {
			return err
		}
	}

	// ── Step 5: Lines changed style (conditional) ────────────────────────────

	if state.HasLines() {
		if err := run(huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("📝 Lines changed — how verbose?").
					Options(
						huh.NewOption(opt("Summary", "📝 ±32"), "summary"),
						huh.NewOption(opt("Detail", "📝 +24 -8"), "detail"),
					).
					Value(&state.LinesStyle),
			),
		)); err != nil {
			return err
		}
	}

	// ── Step 5: Emojis ────────────────────────────────────────────────────────

	if err := run(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("✨ Show emojis?").
				Options(
					huh.NewOption(opt("Yes", "🤖 claude-sonnet-4-6 | 📊 44% | 💰 $2.57"), "all"),
					huh.NewOption(opt("No", "claude-sonnet-4-6 | 44% | $2.57"), "none"),
				).
				Value(&state.Emojis),
		),
	)); err != nil {
		return err
	}

	// ── Step 6: Preview + confirm ─────────────────────────────────────────────

	fmt.Println()
	fmt.Println(sectionStyle.Render("Preview"))
	fmt.Println(subtitleStyle.Render(strings.Repeat("─", 50)))
	fmt.Println(Preview(state))
	fmt.Println(subtitleStyle.Render(strings.Repeat("─", 50)))
	fmt.Println()

	confirm := true
	if err := run(huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("💾 Save this configuration?").
				Value(&confirm),
		),
	)); err != nil {
		return err
	}
	if !confirm {
		fmt.Println(subtitleStyle.Render("Cancelled — no changes made."))
		return nil
	}

	// ── Step 7: Write config ──────────────────────────────────────────────────

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

	// ── Step 8: settings.json ─────────────────────────────────────────────────

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

// opt renders a two-column option label: name in cyan, example in gray.
func opt(name, example string) string {
	return optNameStyle.Render(fmt.Sprintf("%-16s", name)) + optDescStyle.Render(example)
}

func featureOptions() []huh.Option[string] {
	type feature struct {
		key   string
		emoji string
		name  string
		desc  string
	}
	features := []feature{
		{"model", "🤖", "Model name", "which Claude model is active"},
		{"context", "📊", "Context window", "how full the context is"},
		{"tokens", "🎟️ ", "Token counts", "input / output tokens this turn"},
		{"cache", "💾", "Cache", "how much context is reused vs freshly processed"},
		{"cost", "💰", "Cost", "total session cost in USD"},
		{"duration", "⏱️ ", "Duration", "total session time"},
		{"git", "🌿", "Git", "branch name and file change counts"},
		{"lines_changed", "📝", "Lines changed", "total lines added / removed"},
		{"directory", "📁", "Directory", "current working directory"},
		{"agent", "🤖", "Agent name", "shown when running as a sub-agent"},
		{"worktree", "🌿", "Worktree", "shown when working in a git worktree"},
	}
	opts := make([]huh.Option[string], len(features))
	for i, f := range features {
		label := f.emoji + " " + optNameStyle.Render(fmt.Sprintf("%-16s", f.name)) + optDescStyle.Render(f.desc)
		opts[i] = huh.NewOption(label, f.key)
	}
	return opts
}
