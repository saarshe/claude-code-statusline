package wizard

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func runFeaturesStep(state *WizardState) error {
	if err := runWithPreview(huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("What data do you want to see?").
				Options(featureOptions()...).
				Value(&state.Features),
		),
	), state); err != nil {
		return err
	}

	if len(state.Features) == 0 {
		fmt.Println(subtitleStyle.Render("No features selected — setup cancelled."))
		return errNoFeatures
	}
	return nil
}
