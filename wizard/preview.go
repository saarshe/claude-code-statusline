package wizard

import (
	"github.com/saars/claude-code-statusline/render"
	"github.com/saars/claude-code-statusline/schema"
)

// MockInput returns a realistic sample input for wizard preview rendering.
func MockInput() *schema.Input {
	pct := 44.0
	return &schema.Input{
		Model: schema.Model{
			DisplayName: "claude-sonnet-4-6",
		},
		Cwd: "~/dev/my-project",
		Workspace: schema.Workspace{
			CurrentDir: "~/dev/my-project",
		},
		ContextWindow: schema.Context{
			UsedPercentage: &pct,
			CurrentUsage: &schema.Usage{
				InputTokens:              8500,
				OutputTokens:             1200,
				CacheReadInputTokens:     5000,
				CacheCreationInputTokens: 2000,
			},
		},
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
