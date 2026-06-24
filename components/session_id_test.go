package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestSessionID_ShowsFullID(t *testing.T) {
	c := Get("session_id")
	th := theme.Get("default")
	id := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	input := &schema.Input{SessionID: id}

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, id) {
		t.Errorf("expected full session id %q, got %q", id, result)
	}
}

func TestSessionID_EmptyWhenNoID(t *testing.T) {
	c := Get("session_id")
	th := theme.Get("default")

	result := c.Render(&schema.Input{}, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string when no session id, got %q", result)
	}
}

func TestSessionID_NoEmoji(t *testing.T) {
	c := Get("session_id")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	id := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	input := &schema.Input{SessionID: id}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "🆔") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "Session: ") {
		t.Errorf("expected text prefix 'Session: ', got %q", result)
	}
	if !strings.Contains(result, id) {
		t.Errorf("expected session id, got %q", result)
	}
}
