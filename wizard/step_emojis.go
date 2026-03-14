package wizard

import "github.com/charmbracelet/huh"

func runEmojisStep(state *WizardState) error {
	return runWithPreview(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("✨ Show emojis?").
				Options(
					huh.NewOption(opt("Yes", "🤖 claude-sonnet-4-6 | 📊 44% | 💰 $2.57"), "all"),
					huh.NewOption(opt("No", "claude-sonnet-4-6 | 44% | $2.57"), "none"),
				).
				Value(&state.Emojis),
		),
	), state)
}
