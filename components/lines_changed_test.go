package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func linesInput(added, removed int) *schema.Input {
	return &schema.Input{Cost: schema.Cost{TotalLinesAdded: added, TotalLinesRemoved: removed}}
}

func TestLinesChanged_Both(t *testing.T) {
	c := Get("lines_changed")
	th := theme.Get("default")

	result := c.Render(linesInput(156, 23), config.Default(), th)

	if !strings.Contains(result, "+156") {
		t.Errorf("expected '+156', got %q", result)
	}
	if !strings.Contains(result, "-23") {
		t.Errorf("expected '-23', got %q", result)
	}
}

func TestLinesChanged_BothZero_ReturnsEmpty(t *testing.T) {
	c := Get("lines_changed")
	th := theme.Get("default")

	result := c.Render(linesInput(0, 0), config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for zero lines, got %q", result)
	}
}

func TestLinesChanged_OnlyAdded(t *testing.T) {
	c := Get("lines_changed")
	th := theme.Get("default")

	result := c.Render(linesInput(156, 0), config.Default(), th)

	if !strings.Contains(result, "+156") {
		t.Errorf("expected '+156', got %q", result)
	}
	if strings.Contains(result, "-0") {
		t.Errorf("should not show '-0', got %q", result)
	}
}

func TestLinesChanged_OnlyRemoved(t *testing.T) {
	c := Get("lines_changed")
	th := theme.Get("default")

	result := c.Render(linesInput(0, 23), config.Default(), th)

	if strings.Contains(result, "+0") {
		t.Errorf("should not show '+0', got %q", result)
	}
	if !strings.Contains(result, "-23") {
		t.Errorf("expected '-23', got %q", result)
	}
}

func TestLinesChanged_NoEmoji(t *testing.T) {
	c := Get("lines_changed")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(linesInput(156, 23), cfg, th)

	if strings.Contains(result, "📝") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "+156") {
		t.Errorf("expected '+156', got %q", result)
	}
}
