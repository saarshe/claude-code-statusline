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
	data.ContextWindow.TotalInputTokens = 42000
	data.ContextWindow.ContextWindowSize = 200000
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "42k") {
		t.Errorf("expected current tokens in output, got %q", result)
	}
	if !strings.Contains(result, "200k") {
		t.Errorf("expected max window size in output, got %q", result)
	}
}

func TestContextTokens_ZeroWindowSize(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.TotalInputTokens = 1000
	data.ContextWindow.ContextWindowSize = 0
	result := Get("context_tokens").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when window size is 0, got %q", result)
	}
}

func TestContextTokens_ColorsGreen(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.TotalInputTokens = 10000
	data.ContextWindow.ContextWindowSize = 200000 // 5% — green
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
	data.ContextWindow.TotalInputTokens = 5000
	data.ContextWindow.ContextWindowSize = 100000
	result := Get("context_tokens").Render(data, cfg, theme.Get("default"))
	if strings.Contains(result, "📊") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
