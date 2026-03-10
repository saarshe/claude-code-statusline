package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func tokensInput(input, output int) *schema.Input {
	return &schema.Input{
		ContextWindow: schema.Context{
			CurrentUsage: &schema.Usage{
				InputTokens:  input,
				OutputTokens: output,
			},
		},
	}
}

func TestTokens_Normal(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")

	result := c.Render(tokensInput(8500, 1200), config.Default(), th)

	if !strings.Contains(result, "8.5k") {
		t.Errorf("expected '8.5k' for 8500 input tokens, got %q", result)
	}
	if !strings.Contains(result, "1.2k") {
		t.Errorf("expected '1.2k' for 1200 output tokens, got %q", result)
	}
}

func TestTokens_Zero(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")

	result := c.Render(tokensInput(0, 0), config.Default(), th)

	if !strings.Contains(result, "0") {
		t.Errorf("expected '0' for zero tokens, got %q", result)
	}
}

func TestTokens_NilCurrentUsage(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")

	result := c.Render(&schema.Input{}, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for nil CurrentUsage, got %q", result)
	}
}

func TestTokens_NoEmoji(t *testing.T) {
	c := Get("tokens")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(tokensInput(8500, 1200), cfg, th)

	if strings.Contains(result, "🎟") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "8.5k") {
		t.Errorf("expected token counts, got %q", result)
	}
}
