package components

import (
	"fmt"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type contextPctComponent struct{}

func init() { Register(&contextPctComponent{}) }

func (c *contextPctComponent) Key() ComponentKey { return "context_pct" }

func (c *contextPctComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	prefix := GetMeta(c.Key()).Prefix(cfg)

	pct := data.ContextWindow.UsedPercentage
	if pct == nil {
		return th.Success.Render(prefix + "--")
	}

	return ContextStyle(th, *pct, cfg.ContextBar.Thresholds).Render(fmt.Sprintf("%s%.0f%%", prefix, *pct))
}
