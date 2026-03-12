package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	case config.BarGradient:
		return renderGradientBar(filled, empty, width)
	default: // BarBlock
		return strings.Repeat("▓", filled) + strings.Repeat("░", empty)
	}
}

// renderGradientBar colors each filled character by its position's zone:
// green zone → yellow zone → red zone, so danger areas are always visible.
func renderGradientBar(filled, empty, width int) string {
	// zone boundaries as character positions
	greenEnd := int(0.70 * float64(width))
	yellowEnd := int(0.90 * float64(width))

	green  := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	red    := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	dim    := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	var b strings.Builder
	for i := range filled {
		ch := "▓"
		switch {
		case i < greenEnd:
			b.WriteString(green.Render(ch))
		case i < yellowEnd:
			b.WriteString(yellow.Render(ch))
		default:
			b.WriteString(red.Render(ch))
		}
	}
	b.WriteString(dim.Render(strings.Repeat("░", empty)))
	return b.String()
}
