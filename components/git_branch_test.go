package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestGitBranch_WithBranch(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	input := &schema.Input{
		Git: schema.Git{Branch: "main"},
	}

	result := c.Render(input, config.Default(), th)

	if result == "" {
		t.Error("expected non-empty output when branch is set")
	}
	if !strings.Contains(result, "main") {
		t.Errorf("expected branch name in output, got %q", result)
	}
}

func TestGitBranch_NoBranch(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	input := &schema.Input{}

	result := c.Render(input, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty output when no branch, got %q", result)
	}
}

func TestGitBranch_NoEmoji(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{
		Git: schema.Git{Branch: "main"},
	}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "🌿") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
