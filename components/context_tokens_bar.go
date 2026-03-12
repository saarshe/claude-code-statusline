package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type contextTokensBarComponent struct{}

func init() { Register(&contextTokensBarComponent{}) }

func (c *contextTokensBarComponent) Key() ComponentKey { return "context_tokens_bar" }

func (c *contextTokensBarComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	cw := data.ContextWindow
	if cw.ContextWindowSize == 0 {
		return ""
	}

	prefix := EmojiPrefix(cfg, "📊", "")
	bar := renderBar(
		func() float64 {
			if cw.UsedPercentage != nil {
				return *cw.UsedPercentage
			}
			return float64(cw.TotalInputTokens) / float64(cw.ContextWindowSize) * 100
		}(),
		cfg.ContextBar.Style,
		cfg.ContextBar.Width,
	)

	used := HumanizeTokens(cw.TotalInputTokens)
	max := HumanizeTokens(cw.ContextWindowSize)

	pct := 0.0
	if cw.UsedPercentage != nil {
		pct = *cw.UsedPercentage
	}
	style := ContextStyle(th, pct, cfg.ContextBar.Thresholds)

	return style.Render(fmt.Sprintf("%s%s / %s", prefix, used, max)) + " " + bar
}
