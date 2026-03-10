package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type contextPctComponent struct{}

func init() { Register(&contextPctComponent{}) }

func (c *contextPctComponent) Key() ComponentKey { return "context_pct" }

func (c *contextPctComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	prefix := ""
	if cfg.Emojis != config.EmojiNone {
		prefix = "📊 "
	}

	pct := data.ContextWindow.UsedPercentage
	if pct == nil {
		return th.Success.Render(prefix + "--")
	}

	style := contextStyle(th, *pct, cfg.ContextBar.Thresholds)
	return style.Render(fmt.Sprintf("%s%.0f%%", prefix, *pct))
}

func contextStyle(th *theme.Theme, pct float64, thresholds []int) lipgloss.Style {
	if len(thresholds) == 2 {
		if pct >= float64(thresholds[1]) {
			return th.Danger
		}
		if pct >= float64(thresholds[0]) {
			return th.Warning
		}
	}
	return th.Success
}
