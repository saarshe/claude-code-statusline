package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestGitBranch_InGitRepo(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "."}}

	result := c.Render(input, config.Default(), th)

	if result == "" {
		t.Error("expected non-empty branch name in git repo")
	}
	if strings.Contains(result, "fatal") {
		t.Errorf("unexpected git error in output: %q", result)
	}
}

func TestGitBranch_NotInGitRepo(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "/tmp"}}

	result := c.Render(input, config.Default(), th)

	if strings.Contains(result, "fatal") || strings.Contains(result, "error") {
		t.Errorf("should not show git error output, got %q", result)
	}
}

func TestGitBranch_NoEmoji(t *testing.T) {
	c := Get("git_branch")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone
	input := &schema.Input{Workspace: schema.Workspace{CurrentDir: "."}}

	result := c.Render(input, cfg, th)

	if strings.Contains(result, "🌿") {
		t.Errorf("expected no emoji, got %q", result)
	}
}
