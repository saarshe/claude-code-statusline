package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type effortComponent struct{}

func init() { Register(&effortComponent{}) }

func (e *effortComponent) Key() ComponentKey { return "effort" }

func (e *effortComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.Effort == nil || data.Effort.Level == "" {
		return ""
	}
	level := data.Effort.Level
	return effortStyle(th, level).Render(GetMeta(e.Key()).Prefix(cfg) + level)
}

// effortStyle returns a color-ramped style for the level.
// Higher effort surfaces louder colors so it stays visible at a glance.
func effortStyle(th *theme.Theme, level string) lipgloss.Style {
	switch level {
	case "low":
		return th.Muted
	case "medium":
		return th.Secondary
	case "high":
		return th.Warning
	case "xhigh", "max":
		return th.Danger
	default:
		return th.Primary
	}
}
