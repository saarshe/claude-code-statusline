package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestContextTokensBar_ShowsTokensAndBar(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.TotalInputTokens = 88000
	data.ContextWindow.ContextWindowSize = 200000
	pct := 44.0
	data.ContextWindow.UsedPercentage = &pct
	result := Get("context_tokens_bar").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "88k") {
		t.Errorf("expected current tokens in output, got %q", result)
	}
	if !strings.Contains(result, "200k") {
		t.Errorf("expected max tokens in output, got %q", result)
	}
	// bar characters should be present
	if !strings.Contains(result, "▓") && !strings.Contains(result, "░") {
		t.Errorf("expected bar characters in output, got %q", result)
	}
}

func TestContextTokensBar_ZeroWindowSize(t *testing.T) {
	data := &schema.Input{}
	data.ContextWindow.TotalInputTokens = 1000
	data.ContextWindow.ContextWindowSize = 0
	result := Get("context_tokens_bar").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when window size is 0, got %q", result)
	}
}
