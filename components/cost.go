package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type costComponent struct{}

func init() { Register(&costComponent{}) }

func (c *costComponent) Key() ComponentKey { return "cost" }

func (c *costComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if cfg.Emojis != config.EmojiNone {
		return th.Secondary.Render(fmt.Sprintf("💰 $%.2f", data.Cost.TotalCostUSD))
	}
	return th.Secondary.Render(fmt.Sprintf("$%.2f", data.Cost.TotalCostUSD))
}
