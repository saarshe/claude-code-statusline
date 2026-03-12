package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type contextPctComponent struct{}

func init() { Register(&contextPctComponent{}) }

func (c *contextPctComponent) Key() ComponentKey { return "context_pct" }

func (c *contextPctComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	prefix := EmojiPrefix(cfg, "📊", "")

	pct := data.ContextWindow.UsedPercentage
	if pct == nil {
		return th.Success.Render(prefix + "--")
	}

	return ContextStyle(th, *pct, cfg.ContextBar.Thresholds).Render(fmt.Sprintf("%s%.0f%%", prefix, *pct))
}
