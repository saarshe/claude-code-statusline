package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func runFeaturesStep(state *WizardState) error {
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
		return errNoFeatures
	}
	return nil
}
