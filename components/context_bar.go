package components

import (
	"fmt"
	"strings"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type contextBarComponent struct{}

func init() { Register(&contextBarComponent{}) }

func (c *contextBarComponent) Key() ComponentKey { return "context_bar" }

func (c *contextBarComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	pct := data.ContextWindow.UsedPercentage
	if pct == nil {
		return ""
	}

	prefix := ""
	if cfg.Emojis != config.EmojiNone {
		prefix = "📊 "
	}

	bar := renderBar(*pct, cfg.ContextBar.Style, cfg.ContextBar.Width)
	text := fmt.Sprintf("%s%s %.0f%%", prefix, bar, *pct)
	if cfg.ContextBar.Style == config.BarPercent {
		text = fmt.Sprintf("%s%.0f%%", prefix, *pct)
	}
	return contextStyle(th, *pct, cfg.ContextBar.Thresholds).Render(text)
}

func renderBar(pct float64, style config.BarStyle, width int) string {
	if width <= 0 {
		width = 10
	}
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	switch style {
	case config.BarSolid:
		return strings.Repeat("█", filled) + strings.Repeat("░", empty)
	case config.BarASCII:
		return "[" + strings.Repeat("=", filled) + strings.Repeat("-", empty) + "]"
	case config.BarPercent:
		return ""
	default: // BarBlock
		return strings.Repeat("▓", filled) + strings.Repeat("░", empty)
	}
}
