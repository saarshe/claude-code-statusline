package components

import (
	"fmt"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type cacheComponent struct{}

func init() { Register(&cacheComponent{}) }

func (c *cacheComponent) Key() ComponentKey { return "cache" }

// Render shows how many tokens were served from cache ("cached") vs newly
// written into cache ("written"). Use cache_hit for the efficiency percentage.
func (c *cacheComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage

	return th.Secondary.Render(fmt.Sprintf("%s%s reused, %s stored",
		GetMeta(c.Key()).Prefix(cfg),
		HumanizeTokens(u.CacheReadInputTokens),
		HumanizeTokens(u.CacheCreationInputTokens),
	))
}
