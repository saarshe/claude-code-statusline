package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestCost_WithEmoji(t *testing.T) {
	c := Get("cost")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{Cost: schema.Cost{TotalCostUSD: 0.42}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "💰") {
		t.Errorf("expected emoji 💰, got %q", result)
	}
	if !strings.Contains(result, "$0.42") {
		t.Errorf("expected '$0.42', got %q", result)
	}
}

func TestCost_WithoutEmoji(t *testing.T) {
	c := Get("cost")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{Cost: schema.Cost{TotalCostUSD: 0.42}}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "💰") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "$0.42") {
		t.Errorf("expected '$0.42', got %q", result)
	}
}

func TestCost_ZeroCost(t *testing.T) {
	c := Get("cost")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{Cost: schema.Cost{TotalCostUSD: 0}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "$0.00") {
		t.Errorf("expected '$0.00', got %q", result)
	}
}

func TestCost_RoundsToTwoDecimals(t *testing.T) {
	c := Get("cost")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{Cost: schema.Cost{TotalCostUSD: 1.5}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "$1.50") {
		t.Errorf("expected '$1.50', got %q", result)
	}
}
