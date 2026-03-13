package wizard

// Step represents a single wizard step.
type Step struct {
	// ShouldRun returns whether this step applies given current state.
	// If nil, the step always runs.
	ShouldRun func(*WizardState) bool

	// Run executes the step, mutating state. Returns error on failure.
	Run func(*WizardState) error
}

// Steps defines the wizard steps in order.
// Add, remove, or reorder entries to change the wizard flow.
var Steps = []Step{
	{Run: runThemeStep},
	{Run: runFeaturesStep},
	{ShouldRun: (*WizardState).HasContext, Run: runContextStep},
	{ShouldRun: (*WizardState).HasTokens, Run: runTokensStep},
	{ShouldRun: (*WizardState).HasCache, Run: runCacheStep},
	{ShouldRun: (*WizardState).HasGit, Run: runGitStep},
	{ShouldRun: (*WizardState).HasLines, Run: runLinesStep},
	{Run: runEmojisStep},
}
