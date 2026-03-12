package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type linesSummaryComponent struct{}

func init() { Register(&linesSummaryComponent{}) }

func (l *linesSummaryComponent) Key() ComponentKey { return "lines_summary" }

// Render shows total lines touched (added + removed) as a compact ±N figure.
// Use lines_changed for the full +added -removed breakdown.
func (l *linesSummaryComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	added := data.Cost.TotalLinesAdded
	removed := data.Cost.TotalLinesRemoved
	total := added + removed

	if total == 0 {
		return ""
	}

	return th.Secondary.Render(fmt.Sprintf("%s±%d", GetMeta(l.Key()).Prefix(cfg), total))
}
