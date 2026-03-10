package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestCacheHit_ShowsPercentage(t *testing.T) {
	c := Get("cache_hit")
	th := theme.Get("default")
	// 5000 cached out of 5000+8500 = 13500 total input → 37%
	input := cacheInput(5000, 2000)
	input.ContextWindow.CurrentUsage.InputTokens = 8500

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, "%") {
		t.Errorf("expected percentage in output, got %q", result)
	}
}

func TestCacheHit_100Percent(t *testing.T) {
	c := Get("cache_hit")
	th := theme.Get("default")
	input := cacheInput(5000, 0)
	input.ContextWindow.CurrentUsage.InputTokens = 0

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, "100%") {
		t.Errorf("expected '100%%', got %q", result)
	}
}

func TestCacheHit_ZeroCache_ReturnsEmpty(t *testing.T) {
	c := Get("cache_hit")
	th := theme.Get("default")
	// No cache reads → nothing interesting to show
	input := cacheInput(0, 0)
	input.ContextWindow.CurrentUsage.InputTokens = 8500

	result := c.Render(input, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string when no cache reads, got %q", result)
	}
}

func TestCacheHit_NilCurrentUsage_ReturnsEmpty(t *testing.T) {
	c := Get("cache_hit")
	th := theme.Get("default")
	result := c.Render(&schema.Input{}, config.Default(), th)
	if result != "" {
		t.Errorf("expected empty string for nil usage, got %q", result)
	}
}

func TestCacheHit_NoEmoji(t *testing.T) {
	c := Get("cache_hit")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	input := cacheInput(5000, 2000)
	input.ContextWindow.CurrentUsage.InputTokens = 8500
	result := c.Render(input, cfg, th)

	if strings.Contains(result, "⚡") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "%") {
		t.Errorf("expected percentage in output, got %q", result)
	}
}
