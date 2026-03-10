package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestDirectory_ShowsBasename(t *testing.T) {
	c := Get("directory")
	th := theme.Get("default")
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "/home/user/my-project"}}

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, "my-project") {
		t.Errorf("expected 'my-project', got %q", result)
	}
	if strings.Contains(result, "/home/user/") {
		t.Errorf("should not show full path, got %q", result)
	}
}

func TestDirectory_FallsBackToCwd(t *testing.T) {
	c := Get("directory")
	th := theme.Get("default")
	input := &schema.Input{
		Cwd:       "/home/user/fallback-dir",
		Workspace: schema.Workspace{CurrentDir: ""},
	}

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, "fallback-dir") {
		t.Errorf("expected fallback to Cwd 'fallback-dir', got %q", result)
	}
}

func TestDirectory_EmptyDir(t *testing.T) {
	c := Get("directory")
	th := theme.Get("default")

	result := c.Render(&schema.Input{}, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for empty dir, got %q", result)
	}
}

func TestDirectory_NoEmoji(t *testing.T) {
	c := Get("directory")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "/home/user/my-project"}}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "📁") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "my-project") {
		t.Errorf("expected dirname, got %q", result)
	}
}
