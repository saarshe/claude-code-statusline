package wizard

import "github.com/charmbracelet/huh"

func runLayoutStep(state *WizardState) error {
	return runWithPreview(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📐 Layout — how should components fill the lines?").
				Description(
					"Responsive: ignore the line breaks and reflow components to fit\n"+
						"your terminal width, recomputed on each update.\n"+
						"Fixed: keep the line breaks below exactly as shown.\n",
				).
				Options(
					huh.NewOption("Responsive (reflow to width)", "auto"),
					huh.NewOption("Fixed lines", "fixed"),
				).
				Value(&state.Layout),
		),
	), state)
}
