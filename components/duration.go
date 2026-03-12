package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type durationComponent struct{}

func init() { Register(&durationComponent{}) }

func (d *durationComponent) Key() ComponentKey { return "duration" }

func (d *durationComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	ms := data.Cost.TotalDurationMS

	return th.Secondary.Render(GetMeta(d.Key()).Prefix(cfg) + formatDuration(ms))
}

func formatDuration(ms int64) string {
	totalSecs := ms / 1000
	hours := totalSecs / 3600
	mins := (totalSecs % 3600) / 60
	secs := totalSecs % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
