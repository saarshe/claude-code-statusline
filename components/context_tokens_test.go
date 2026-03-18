package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestContextTokens_NormalUsage(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 200000
	data.ContextWindow.CurrentUsage = &schema.Usage{
		InputTokens: 2000, CacheReadInputTokens: 38000, CacheCreationInputTokens: 2000,
	}
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "42k") {
		t.Errorf("expected context fill tokens in output, got %q", result)
	}
	if !strings.Contains(result, "200k") {
		t.Errorf("expected max window size in output, got %q", result)
	}
}

func TestContextTokens_ZeroWindowSize(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 0
	data.ContextWindow.CurrentUsage = &schema.Usage{InputTokens: 1000}
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when window size is 0, got %q", result)
	}
}

func TestContextTokens_NilCurrentUsage(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.TotalInputTokens = 5000 // should be ignored
	data.ContextWindow.ContextWindowSize = 200000
	data.ContextWindow.CurrentUsage = nil
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when current_usage is nil, got %q", result)
	}
}

func TestContextTokens_ColorsGreen(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 200000
	data.ContextWindow.CurrentUsage = &schema.Usage{InputTokens: 10000}
	pct := 5.0
	data.ContextWindow.UsedPercentage = &pct
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestContextTokens_NoEmoji(t *testing.T) {
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	data := &schema.Input{}
	data.ContextWindow.ContextWindowSize = 100000
	data.ContextWindow.CurrentUsage = &schema.Usage{InputTokens: 5000}
	result := Get("context_tokens").Render(data, cfg, theme.Get("default"))
	if strings.Contains(result, "📊") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
