package wizard

import "github.com/charmbracelet/huh"

func runContextStep(state *WizardState) error {
	style, err := runContextStyleSelector(state.ContextStyle)
	if err != nil {
		return err
	}
	state.ContextStyle = style
	return nil
}

func runTokensStep(state *WizardState) error {
	tokenExamples := map[string]string{
		"turn":       "🎟️ In: 112k Out: 514",
		"turn_cache": "🎟️ In: 112k (99% cached) Out: 514",
		"session":    "🎟️ 35k out",
		"full":       "🎟️ In: 112k (99% cached) Out: 514 · Session: 35k out",
	}
	return run(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🎟️  Tokens — how verbose?").
				Description(
					"Input = total tokens sent to the model (including cached).\n"+
						"Output = tokens the model generated (the main cost driver — 5x input price).\n"+
						"Cache hit % = how much input was reused cheaply from the previous turn.\n",
				).
				Options(styleOptions("tokens", tokenExamples)...).
				Value(&state.TokenStyle),
		),
	))
}

func runCacheStep(state *WizardState) error {
	cacheExamples := map[string]string{
		"hit":    "⚡ 37% cached",
		"counts": "💾 5.0k reused, 2.0k stored",
	}
	return run(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("💾 Cache — how verbose?").
				Description(
					"Each turn, Claude reuses previously processed context from cache (cheap)\n"+
						"and stores new context for the next turn to reuse.\n",
				).
				Options(styleOptions("cache", cacheExamples)...).
				Value(&state.CacheStyle),
		),
	))
}

func runGitStep(state *WizardState) error {
	gitExamples := map[string]string{
		"branch": "🌿 main",
		"status": "🌿 main +1 ~9",
	}
	return run(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🌿 Git — how verbose?").
				Options(styleOptions("git", gitExamples)...).
				Value(&state.GitStyle),
		),
	))
}

func runLinesStep(state *WizardState) error {
	linesExamples := map[string]string{
		"summary": "📝 ±32",
		"detail":  "📝 +24 -8",
	}
	return run(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📝 Lines changed — how verbose?").
				Options(styleOptions("lines_changed", linesExamples)...).
				Value(&state.LinesStyle),
		),
	))
}
