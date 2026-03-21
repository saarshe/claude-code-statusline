package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestGitStatus_WithBranch(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	input := &schema.Input{
		Git: schema.Git{Branch: "main", Staged: 2, Modified: 3},
	}

	result := c.Render(input, config.Default(), th)

	if result == "" {
		t.Error("expected non-empty output")
	}
	if !strings.Contains(result, "main") {
		t.Errorf("expected branch name, got %q", result)
	}
	if !strings.Contains(result, "+2") {
		t.Errorf("expected staged count, got %q", result)
	}
	if !strings.Contains(result, "~3") {
		t.Errorf("expected modified count, got %q", result)
	}
}

func TestGitStatus_NoBranch(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	input := &schema.Input{}

	result := c.Render(input, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty output when no branch, got %q", result)
	}
}

func TestGitStatus_CleanRepo(t *testing.T) {
	c := Get("git_status")
	th := theme.Get("default")
	input := &schema.Input{
		Git: schema.Git{Branch: "main"},
	}

	result := c.Render(input, config.Default(), th)

	if !strings.Contains(result, "main") {
		t.Errorf("expected branch name, got %q", result)
	}
	if strings.Contains(result, "+") || strings.Contains(result, "~") {
		t.Errorf("expected no counts for clean repo, got %q", result)
	}
}

func TestGitStatus_NoEmoji(t *testing.T) {
	c := Get("git_status")
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
