package components

import (
	"fmt"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type contextTokensBarPctComponent struct{}

func init() { Register(&contextTokensBarPctComponent{}) }

func (c *contextTokensBarPctComponent) Key() ComponentKey { return "context_tokens_bar_pct" }

func (c *contextTokensBarPctComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	cw := data.ContextWindow
	if cw.ContextWindowSize == 0 || cw.CurrentUsage == nil {
		return ""
	}

	prefix := GetMeta(c.Key()).Prefix(cfg)

	pct := 0.0
	if cw.UsedPercentage != nil {
		pct = *cw.UsedPercentage
	} else {
		pct = float64(cw.ContextFillTokens()) / float64(cw.ContextWindowSize) * 100
	}

	bar := renderBar(pct, cfg.ContextBar.Style, cfg.ContextBar.Width)
	used := HumanizeTokens(cw.ContextFillTokens())
	max := HumanizeTokens(cw.ContextWindowSize)
	style := ContextStyle(th, pct, cfg.ContextBar.Thresholds)

	tokens := style.Render(fmt.Sprintf("%s%s / %s", prefix, used, max))
	pctStr := style.Render(fmt.Sprintf("%.0f%%", pct))

	if cfg.ContextBar.Style == config.BarGradient {
		return tokens + " " + bar + " " + pctStr
	}
	return tokens + " " + style.Render(bar) + " " + pctStr
}
