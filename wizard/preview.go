package wizard

import (
	"os"

	"github.com/saars/claude-code-statusline/render"
	"github.com/saars/claude-code-statusline/schema"
)

// MockInput returns a realistic sample input for wizard preview rendering.
// It uses the actual working directory so that git components (branch, status)
// produce real output.
func MockInput() *schema.Input {
	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = "."
	}
	pct := 44.0
	return &schema.Input{
		Model: schema.Model{
			DisplayName: "claude-sonnet-4-6",
		},
		Cwd: cwd,
		Workspace: schema.Workspace{
			CurrentDir: cwd,
		},
		ContextWindow: schema.Context{
			UsedPercentage:    &pct,
			TotalInputTokens:  88000,
			ContextWindowSize: 200000,
			CurrentUsage: &schema.Usage{
				InputTokens:              8500,
				OutputTokens:             1200,
				CacheReadInputTokens:     5000,
				CacheCreationInputTokens: 2000,
			},
		},
		Agent:    &schema.Agent{Name: "subagent"},
		Worktree: &schema.Worktree{Name: "feature-branch"},
		Cost: schema.Cost{
			TotalCostUSD:      2.57,
			TotalDurationMS:   83000,
			TotalLinesAdded:   24,
			TotalLinesRemoved: 8,
		},
	}
}

// Preview renders the status line using mock data and the given wizard state.
func Preview(state *WizardState) string {
	cfg := state.ToConfig()
	return render.Render(MockInput(), cfg)
}
