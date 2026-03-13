package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

// tokens — per-turn totals: "In: 112k Out: 514"
type tokensComponent struct{}

func init() { Register(&tokensComponent{}) }

func (t *tokensComponent) Key() ComponentKey { return "tokens" }

func (t *tokensComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage

	return th.Secondary.Render(fmt.Sprintf("%sIn: %s Out: %s",
		GetMeta(t.Key()).Prefix(cfg),
		HumanizeTokens(u.TotalInput()),
		HumanizeTokens(u.OutputTokens),
	))
}

// tokens_cache — per-turn with cache: "In: 112k (99% cached) Out: 514"
type tokensCacheComponent struct{}

func init() { Register(&tokensCacheComponent{}) }

func (t *tokensCacheComponent) Key() ComponentKey { return "tokens_cache" }

func (t *tokensCacheComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage

	return th.Secondary.Render(fmt.Sprintf("%sIn: %s (%d%% cached) Out: %s",
		GetMeta(t.Key()).Prefix(cfg),
		HumanizeTokens(u.TotalInput()),
		u.CacheHitPct(),
		HumanizeTokens(u.OutputTokens),
	))
}

// tokens_session — session cumulative output: "↓35k out"
// Only shows output tokens because total_input_tokens from Claude Code
// excludes cached tokens, making it misleadingly low.
type tokensSessionComponent struct{}

func init() { Register(&tokensSessionComponent{}) }

func (t *tokensSessionComponent) Key() ComponentKey { return "tokens_session" }

func (t *tokensSessionComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	cw := data.ContextWindow
	if cw.TotalOutputTokens == 0 {
		return ""
	}

	return th.Secondary.Render(fmt.Sprintf("%s%s out",
		GetMeta(t.Key()).Prefix(cfg),
		HumanizeTokens(cw.TotalOutputTokens),
	))
}

// tokens_full — full breakdown: "In: 112k (99% cached) Out: 514 · Total ↓35k"
type tokensFullComponent struct{}

func init() { Register(&tokensFullComponent{}) }

func (t *tokensFullComponent) Key() ComponentKey { return "tokens_full" }

func (t *tokensFullComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage
	cw := data.ContextWindow

	return th.Secondary.Render(fmt.Sprintf("%sIn: %s (%d%% cached) Out: %s · Session: %s out",
		GetMeta(t.Key()).Prefix(cfg),
		HumanizeTokens(u.TotalInput()),
		u.CacheHitPct(),
		HumanizeTokens(u.OutputTokens),
		HumanizeTokens(cw.TotalOutputTokens),
	))
}
