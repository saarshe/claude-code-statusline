package wizard

import (
	"time"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/render"
	"github.com/saarshe/claude-code-statusline/schema"
)

// MockInput returns a realistic sample input for wizard preview rendering.
// All data is static — no subprocess calls or I/O are performed.
func MockInput() *schema.Input {
	pct := 44.0
	// Reset times are computed from now so the rate_limits_reset preview
	// shows a sensible countdown ("in 2h14m", "in 3d") whenever the wizard runs.
	now := time.Now().Unix()
	return &schema.Input{
		Model: schema.Model{
			DisplayName: "claude-sonnet-4-6",
		},
		SessionID: "7f3a9c20-1b8e-4d6a-9f12-0a4b7c8e5d31",
		Cwd:       "/home/user/project",
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
		Effort:   &schema.Effort{Level: "high"},
		Plan:     schema.Plan{OrgType: "claude_max", RateLimitTier: "default_claude_max_20x"},
		RateLimits: &schema.RateLimits{
			FiveHour: &schema.RateLimitWindow{UsedPercentage: 23, ResetsAt: now + 2*3600 + 14*60},
			SevenDay: &schema.RateLimitWindow{UsedPercentage: 41, ResetsAt: now + 3*86400},
		},
		PR: &schema.PR{
			Number:      1234,
			URL:         "https://github.com/saarshe/claude-code-statusline/pull/1234",
			ReviewState: "pending",
		},
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
//
// It always renders the wizard-width-inferred rows verbatim, even in auto mode:
// at the wizard's current width that is exactly what the live renderer would
// produce, and it avoids depending on COLUMNS (which may be unset in the wizard
// process).
func Preview(state *WizardState) string {
	cfg := state.toConfigWithLayout(state.InferLayout())
	cfg.Layout = config.LayoutFixed
	return render.Render(MockInput(), cfg)
}
