package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestModel_WithEmoji(t *testing.T) {
	c := Get("model")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{Model: schema.Model{DisplayName: "Opus"}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "🤖") {
		t.Errorf("expected emoji 🤖, got %q", result)
	}
	if !strings.Contains(result, "Opus") {
		t.Errorf("expected model name 'Opus', got %q", result)
	}
}

func TestModel_WithoutEmoji(t *testing.T) {
	c := Get("model")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{Model: schema.Model{DisplayName: "Opus"}}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "🤖") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "Opus") {
		t.Errorf("expected model name 'Opus', got %q", result)
	}
	if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
		t.Errorf("expected brackets around model name, got %q", result)
	}
}

func TestModel_EmptyDisplayName(t *testing.T) {
	c := Get("model")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{Model: schema.Model{DisplayName: ""}}

	result := c.Render(input, cfg, th)

	if result != "" {
		t.Errorf("expected empty string for empty display name, got %q", result)
	}
}
