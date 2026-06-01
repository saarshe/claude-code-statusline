package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestEffort_RendersLevel(t *testing.T) {
	for _, lvl := range []string{"low", "medium", "high", "xhigh", "max"} {
		data := &schema.Input{Effort: &schema.Effort{Level: lvl}}
		result := Get("effort").Render(data, config.Default(), theme.Get("default"))
		if !strings.Contains(result, lvl) {
			t.Errorf("expected %q in output, got %q", lvl, result)
		}
	}
}

func TestEffort_EmptyWhenNil(t *testing.T) {
	data := &schema.Input{}
	result := Get("effort").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty when Effort nil, got %q", result)
	}
}

func TestEffort_EmptyWhenLevelEmpty(t *testing.T) {
	data := &schema.Input{Effort: &schema.Effort{Level: ""}}
	result := Get("effort").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty when level empty, got %q", result)
	}
}

func TestEffort_EmojiPrefix(t *testing.T) {
	data := &schema.Input{Effort: &schema.Effort{Level: "high"}}
	cfg := config.Default()
	result := Get("effort").Render(data, cfg, theme.Get("default"))
	if !strings.Contains(result, "🧠") {
		t.Errorf("expected 🧠 emoji, got %q", result)
	}

	cfg.Emojis = config.EmojiNone
	result = Get("effort").Render(data, cfg, theme.Get("default"))
	if strings.Contains(result, "🧠") {
		t.Errorf("expected no emoji when disabled, got %q", result)
	}
	if !strings.Contains(result, "Effort:") {
		t.Errorf("expected text prefix when emojis off, got %q", result)
	}
}
