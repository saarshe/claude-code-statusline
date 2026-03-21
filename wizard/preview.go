package wizard

import (
	"github.com/saarshe/claude-code-statusline/render"
	"github.com/saarshe/claude-code-statusline/schema"
)

// MockInput returns a realistic sample input for wizard preview rendering.
// All data is static — no subprocess calls or I/O are performed.
func MockInput() *schema.Input {
	pct := 44.0
	return &schema.Input{
		Model: schema.Model{
			DisplayName: "claude-sonnet-4-6",
		},
		Cwd: "/home/user/project",
		Workspace: schema.Workspace{
			CurrentDir: "/home/user/project",
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
		Git: schema.Git{
			Branch:   "main",
			Staged:   2,
			Modified: 3,
		},
	}
}

// Preview renders the status line using mock data and the given wizard state.
func Preview(state *WizardState) string {
	cfg := state.ToConfig()
	return render.Render(MockInput(), cfg)
}
