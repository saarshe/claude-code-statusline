package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type prComponent struct{}

func init() { Register(&prComponent{}) }

func (p *prComponent) Key() ComponentKey { return "pr" }

func (p *prComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.PR == nil || data.PR.Number == 0 {
		return ""
	}

	body := fmt.Sprintf("%s#%d", GetMeta(p.Key()).Prefix(cfg), data.PR.Number)
	if mark := reviewStateMark(data.PR.ReviewState); mark != "" {
		body += " " + mark
	}

	styled := reviewStateStyle(th, data.PR.ReviewState).Render(body)

	if data.PR.URL != "" {
		return wrapHyperlink(data.PR.URL, styled)
	}
	return styled
}

// reviewStateMark returns a short glyph for a PR review state.
func reviewStateMark(state string) string {
	switch state {
	case "approved":
		return "✓"
	case "changes_requested":
		return "✗"
	case "pending":
		return "⏳"
	case "draft":
		return "·"
	default:
		return ""
	}
}

func reviewStateStyle(th *theme.Theme, state string) lipgloss.Style {
	switch state {
	case "approved":
		return th.Success
	case "changes_requested":
		return th.Danger
	case "pending":
		return th.Warning
	case "draft":
		return th.Muted
	default:
		return th.Primary
	}
}

// wrapHyperlink wraps text in an OSC 8 escape sequence so terminals that
// support hyperlinks make it clickable. Uses BEL terminator for broad
// compatibility (matches the docs example).
func wrapHyperlink(url, text string) string {
	return "\x1b]8;;" + url + "\x07" + text + "\x1b]8;;\x07"
}
