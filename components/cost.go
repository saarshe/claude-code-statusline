package components

import (
	"fmt"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type costComponent struct{}

func init() { Register(&costComponent{}) }

func (c *costComponent) Key() ComponentKey { return "cost" }

func (c *costComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	return th.Secondary.Render(fmt.Sprintf("%s$%.2f", GetMeta(c.Key()).Prefix(cfg), data.Cost.TotalCostUSD))
}
