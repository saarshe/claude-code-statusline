package components

import (
	"fmt"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type cacheHitComponent struct{}

func init() { Register(&cacheHitComponent{}) }

func (c *cacheHitComponent) Key() ComponentKey { return "cache_hit" }

// Render shows the cache hit rate: what percentage of input tokens were served
// from cache (cheap) vs sent fresh to the model (full cost). Returns empty
// when there are no cache reads, since 0% hit is not worth showing.
func (c *cacheHitComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.ContextWindow.CurrentUsage == nil {
		return ""
	}
	u := data.ContextWindow.CurrentUsage

	cached := u.CacheReadInputTokens
	if cached == 0 {
		return ""
	}
	total := cached + u.InputTokens
	hitPct := int(float64(cached) / float64(total) * 100)

	prefix := ""
	if cfg.Emojis != config.EmojiNone {
		prefix = "⚡ "
	} else {
		prefix = "Cache: "
	}

	return th.Secondary.Render(fmt.Sprintf("%s%d%% cached", prefix, hitPct))
}
