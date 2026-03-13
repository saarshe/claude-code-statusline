package wizard

func runThemeStep(state *WizardState) error {
	selectedTheme, err := runThemeSelector(state.Theme)
	if err != nil {
		return err
	}
	state.Theme = selectedTheme
	return nil
}
