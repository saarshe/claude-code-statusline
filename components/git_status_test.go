package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestGitStatus_InGitRepo(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "."}}

	result := c.Render(input, config.Default(), th)

	if result == "" {
		t.Error("expected non-empty output in git repo")
	}
	if strings.Contains(result, "fatal") || strings.Contains(result, "error") {
		t.Errorf("should not show git error, got %q", result)
	}
}

func TestGitStatus_NotInGitRepo(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "/tmp"}}

	result := c.Render(input, config.Default(), th)

	if strings.Contains(result, "fatal") || strings.Contains(result, "error") {
		t.Errorf("should not show git error output, got %q", result)
	}
}

func TestGitStatus_NoEmoji(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "."}}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "🌿") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
