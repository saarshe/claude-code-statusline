package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestContextTokensBarPct_ShowsAll(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 200000
	data.ContextWindow.CurrentUsage = &schema.Usage{
		InputTokens: 8000, CacheReadInputTokens: 76000, CacheCreationInputTokens: 4000,
	}
	pct := 44.0
	data.ContextWindow.UsedPercentage = &pct
	result := Get("context_tokens_bar_pct").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "88k") {
		t.Errorf("expected context fill tokens, got %q", result)
	}
	if !strings.Contains(result, "200k") {
		t.Errorf("expected max tokens, got %q", result)
	}
	if !strings.Contains(result, "44%") {
		t.Errorf("expected percentage, got %q", result)
	}
	if !strings.Contains(result, "▓") && !strings.Contains(result, "░") {
		t.Errorf("expected bar characters, got %q", result)
	}
}

func TestContextTokensBarPct_NilCurrentUsage(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 200000
	data.ContextWindow.CurrentUsage = nil
	result := Get("context_tokens_bar_pct").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when current_usage is nil, got %q", result)
	}
}

func TestContextTokensBarPct_ZeroWindowSize(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 0
	data.ContextWindow.CurrentUsage = &schema.Usage{InputTokens: 1000}
	result := Get("context_tokens_bar_pct").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when window size is 0, got %q", result)
	}
}
