package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/theme"
)

func TestLinesSummary_ShowsTotal(t *testing.T) {
	c := Get("lines_summary")
	th := theme.Get("default")

	result := c.Render(linesInput(24, 8), config.Default(), th)

	// Total = 24 + 8 = 32
	if !strings.Contains(stripANSI(result), "32") {
		t.Errorf("expected total '32' lines in output, got %q", result)
	}
}

func TestLinesSummary_ZeroLines_ReturnsEmpty(t *testing.T) {
	c := Get("lines_summary")
	th := theme.Get("default")

	result := c.Render(linesInput(0, 0), config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for zero lines, got %q", result)
	}
}

func TestLinesSummary_OnlyAdded(t *testing.T) {
	c := Get("lines_summary")
	th := theme.Get("default")

	result := c.Render(linesInput(42, 0), config.Default(), th)

	if !strings.Contains(stripANSI(result), "42") {
		t.Errorf("expected '42' in output, got %q", result)
	}
}

func TestLinesSummary_NoEmoji(t *testing.T) {
	c := Get("lines_summary")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(linesInput(24, 8), cfg, th)

	if strings.Contains(result, "📝") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(stripANSI(result), "32") {
		t.Errorf("expected total in output, got %q", result)
	}
}
