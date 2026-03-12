package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type tokensComponent struct{}

func init() { Register(&tokensComponent{}) }

func (t *tokensComponent) Key() ComponentKey { return "tokens" }

func (t *tokensComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage

	return th.Secondary.Render(fmt.Sprintf("%sIn: %s Out: %s",
		GetMeta(t.Key()).Prefix(cfg),
		HumanizeTokens(u.InputTokens),
		HumanizeTokens(u.OutputTokens),
	))
}
