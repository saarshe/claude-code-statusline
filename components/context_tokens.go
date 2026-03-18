package components

import (
	"fmt"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type contextTokensComponent struct{}

func init() { Register(&contextTokensComponent{}) }

func (c *contextTokensComponent) Key() ComponentKey { return "context_tokens" }

func (c *contextTokensComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	cw := data.ContextWindow
	if cw.ContextWindowSize == 0 || cw.CurrentUsage == nil {
		return ""
	}

	prefix := GetMeta(c.Key()).Prefix(cfg)
	used := HumanizeTokens(cw.ContextFillTokens())
	max := HumanizeTokens(cw.ContextWindowSize)

	style := th.Success
	if cw.UsedPercentage != nil {
		style = ContextStyle(th, *cw.UsedPercentage, cfg.ContextBar.Thresholds)
	}

	return style.Render(fmt.Sprintf("%s%s / %s", prefix, used, max))
}
