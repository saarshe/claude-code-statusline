package wizard

import (
	"errors"

	"github.com/charmbracelet/huh"
)

func runFeaturesStep(state *WizardState) error {
	return runWithPreview(huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("What data do you want to see?").
				Options(featureOptions()...).
				Validate(func(selected []string) error {
					if len(selected) == 0 {
						return errors.New("select at least one feature")
					}
					return nil
				}).
				Value(&state.Features),
		),
	), state)
}
