package components

import (
	"fmt"
	"strings"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type linesChangedComponent struct{}

func init() { Register(&linesChangedComponent{}) }

func (l *linesChangedComponent) Key() ComponentKey { return "lines_changed" }

func (l *linesChangedComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	added := data.Cost.TotalLinesAdded
	removed := data.Cost.TotalLinesRemoved

	if added == 0 && removed == 0 {
		return ""
	}

	prefix := EmojiPrefix(cfg, "📝", "")
	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("+%d", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d", removed))
	}

	return th.Secondary.Render(prefix + strings.Join(parts, " "))
}
