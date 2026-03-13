package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func tokensInput(input, output, cacheRead, cacheCreation int) *schema.Input {
	return &schema.Input{
		ContextWindow: schema.Context{
			TotalInputTokens:  150000,
			TotalOutputTokens: 35000,
			CurrentUsage: &schema.Usage{
				InputTokens:              input,
				OutputTokens:             output,
				CacheReadInputTokens:     cacheRead,
				CacheCreationInputTokens: cacheCreation,
			},
		},
	}
}

// --- tokens (per-turn totals) ---

func TestTokens_ShowsTotalInput(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")

	// input=3, cacheRead=112000, cacheCreation=500 → total=112503 → "112k"
	result := c.Render(tokensInput(3, 514, 112000, 500), config.Default(), th)

	if !strings.Contains(result, "112k") {
		t.Errorf("expected total input '112k', got %q", result)
	}
	if !strings.Contains(result, "514") {
		t.Errorf("expected output '514', got %q", result)
	}
}

func TestTokens_NilCurrentUsage(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")
	result := c.Render(&schema.Input{}, config.Default(), th)
	if result != "" {
		t.Errorf("expected empty for nil CurrentUsage, got %q", result)
	}
}

func TestTokens_NoEmoji(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(tokensInput(3, 514, 112000, 500), cfg, th)
	if strings.Contains(result, "🎟") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "Tok:") {
		t.Errorf("expected text prefix, got %q", result)
	}
}

// --- tokens_cache (per-turn + cache hit) ---

func TestTokensCache_ShowsCacheHit(t *testing.T) {
	c := Get("tokens_cache")
	th := theme.Get("default")

	result := c.Render(tokensInput(3, 514, 112000, 500), config.Default(), th)

	if !strings.Contains(result, "99% cached") {
		t.Errorf("expected cache hit pct, got %q", result)
	}
	if !strings.Contains(result, "112k") {
		t.Errorf("expected total input, got %q", result)
	}
}

func TestTokensCache_NilCurrentUsage(t *testing.T) {
	c := Get("tokens_cache")
	th := theme.Get("default")
	result := c.Render(&schema.Input{}, config.Default(), th)
	if result != "" {
		t.Errorf("expected empty for nil CurrentUsage, got %q", result)
	}
}

// --- tokens_session (session cumulative) ---

func TestTokensSession_ShowsCumulativeOutput(t *testing.T) {
	c := Get("tokens_session")
	th := theme.Get("default")

	result := c.Render(tokensInput(3, 514, 112000, 500), config.Default(), th)

	if !strings.Contains(result, "35k") {
		t.Errorf("expected cumulative output '35k', got %q", result)
	}
	if !strings.Contains(result, "out") {
		t.Errorf("expected 'out' label, got %q", result)
	}
}

func TestTokensSession_EmptyWhenZero(t *testing.T) {
	c := Get("tokens_session")
	th := theme.Get("default")

	input := &schema.Input{
		ContextWindow: schema.Context{
			TotalInputTokens:  0,
			TotalOutputTokens: 0,
		},
	}
	result := c.Render(input, config.Default(), th)
	if result != "" {
		t.Errorf("expected empty for zero totals, got %q", result)
	}
}

// --- tokens_full (full breakdown) ---

func TestTokensFull_ShowsEverything(t *testing.T) {
	c := Get("tokens_full")
	th := theme.Get("default")

	result := c.Render(tokensInput(3, 514, 112000, 500), config.Default(), th)

	if !strings.Contains(result, "112k") {
		t.Errorf("expected total input, got %q", result)
	}
	if !strings.Contains(result, "99% cached") {
		t.Errorf("expected cache hit pct, got %q", result)
	}
	if !strings.Contains(result, "514") {
		t.Errorf("expected turn output, got %q", result)
	}
	if !strings.Contains(result, "35k") {
		t.Errorf("expected cumulative output, got %q", result)
	}
}

func TestTokensFull_NilCurrentUsage(t *testing.T) {
	c := Get("tokens_full")
	th := theme.Get("default")
	result := c.Render(&schema.Input{}, config.Default(), th)
	if result != "" {
		t.Errorf("expected empty for nil CurrentUsage, got %q", result)
	}
}
