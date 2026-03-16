package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func cacheInput(read, write int) *schema.Input {
	return &schema.Input{
		ContextWindow: schema.Context{
			CurrentUsage: &schema.Usage{
				InputTokens:              100,
				CacheReadInputTokens:     read,
				CacheCreationInputTokens: write,
			},
		},
	}
}

func TestCache_Normal(t *testing.T) {
	c := Get("cache")
	th := theme.Get("default")

	result := c.Render(cacheInput(5000, 2000), config.Default(), th)

	if !strings.Contains(result, "5.0k") {
		t.Errorf("expected '5.0k' for 5000 cached tokens, got %q", result)
	}
	if !strings.Contains(result, "2.0k") {
		t.Errorf("expected '2.0k' for 2000 written tokens, got %q", result)
	}
	if !strings.Contains(result, "reused") {
		t.Errorf("expected 'reused' label, got %q", result)
	}
	if !strings.Contains(result, "stored") {
		t.Errorf("expected 'stored' label, got %q", result)
	}
}

func TestCache_NilCurrentUsage(t *testing.T) {
	c := Get("cache")
	th := theme.Get("default")

	result := c.Render(&schema.Input{}, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for nil CurrentUsage, got %q", result)
	}
}

func TestCache_ZeroValues(t *testing.T) {
	c := Get("cache")
	th := theme.Get("default")

	result := c.Render(cacheInput(0, 0), config.Default(), th)

	// Should still render with zero values
	if result == "" {
		t.Error("expected non-empty output even for zero values")
	}
}

func TestCache_NoEmoji(t *testing.T) {
	c := Get("cache")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(cacheInput(5000, 2000), cfg, th)

	if strings.Contains(result, "💾") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
