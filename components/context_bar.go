package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type contextBarComponent struct{}

func init() { Register(&contextBarComponent{}) }

func (c *contextBarComponent) Key() ComponentKey { return "context_bar" }

func (c *contextBarComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	pct := data.ContextWindow.UsedPercentage
	if pct == nil {
		return ""
	}

	prefix := GetMeta(c.Key()).Prefix(cfg)
	bar := renderBar(*pct, cfg.ContextBar.Style, cfg.ContextBar.Width, th)

	// Gradient bar has per-character colors; wrapping with an outer Render would
	// override them. Render the percentage separately and concatenate.
	if cfg.ContextBar.Style == config.BarGradient {
		pctStr := ContextStyle(th, *pct, cfg.ContextBar.Thresholds).Render(fmt.Sprintf("%.0f%%", *pct))
		return prefix + bar + " " + pctStr
	}

	text := fmt.Sprintf("%s%s %.0f%%", prefix, bar, *pct)
	if cfg.ContextBar.Style == config.BarPercent {
		text = fmt.Sprintf("%s%.0f%%", prefix, *pct)
	}
	return ContextStyle(th, *pct, cfg.ContextBar.Thresholds).Render(text)
}

func renderBar(pct float64, style config.BarStyle, width int, th *theme.Theme) string {
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
		return renderGradientBar(filled, empty, width, th)
	default: // BarBlock
		return strings.Repeat("█", filled) + strings.Repeat("░", empty)
	}
}

var gradDim = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

// renderGradientBar batches consecutive characters in the same color zone
// so themes with background colors render a continuous bar, not per-char boxes.
func renderGradientBar(filled, empty, width int, th *theme.Theme) string {
	greenEnd := Clamp(int(0.70*float64(width)), 0, filled)
	yellowEnd := Clamp(int(0.90*float64(width)), 0, filled)

	var b strings.Builder
	if greenEnd > 0 {
		b.WriteString(th.Success.Render(strings.Repeat("█", greenEnd)))
	}
	if yellowEnd > greenEnd {
		b.WriteString(th.Warning.Render(strings.Repeat("█", yellowEnd-greenEnd)))
	}
	if filled > yellowEnd {
		b.WriteString(th.Danger.Render(strings.Repeat("█", filled-yellowEnd)))
	}
	if empty > 0 {
		b.WriteString(gradDim.Render(strings.Repeat("░", empty)))
	}
	return b.String()
}
