package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func durationInput(ms int64) *schema.Input {
	return &schema.Input{Cost: schema.Cost{TotalDurationMS: ms}}
}

func TestDuration_Seconds(t *testing.T) {
	c := Get("duration")
	th := theme.Get("default")

	result := c.Render(durationInput(45000), config.Default(), th)

	if !strings.Contains(result, "45s") {
		t.Errorf("expected '45s', got %q", result)
	}
}

func TestDuration_MinutesAndSeconds(t *testing.T) {
	c := Get("duration")
	th := theme.Get("default")

	result := c.Render(durationInput(754000), config.Default(), th)

	if !strings.Contains(result, "12m") {
		t.Errorf("expected '12m', got %q", result)
	}
	if !strings.Contains(result, "34s") {
		t.Errorf("expected '34s', got %q", result)
	}
}

func TestDuration_Hours(t *testing.T) {
	c := Get("duration")
	th := theme.Get("default")

	result := c.Render(durationInput(3661000), config.Default(), th)

	if !strings.Contains(result, "1h") {
		t.Errorf("expected '1h', got %q", result)
	}
	if !strings.Contains(result, "1m") {
		t.Errorf("expected '1m', got %q", result)
	}
}

func TestDuration_Zero(t *testing.T) {
	c := Get("duration")
	th := theme.Get("default")

	result := c.Render(durationInput(0), config.Default(), th)

	if !strings.Contains(result, "0s") {
		t.Errorf("expected '0s', got %q", result)
	}
}

func TestDuration_NoEmoji(t *testing.T) {
	c := Get("duration")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(durationInput(45000), cfg, th)

	if strings.Contains(result, "⏱") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "45s") {
		t.Errorf("expected duration, got %q", result)
	}
}
