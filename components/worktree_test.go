package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestWorktree_ShowsNameWhenActive(t *testing.T) {
	data := &schema.Input{Worktree: &schema.Worktree{Name: "my-feature"}}
	result := Get("worktree").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "my-feature") {
		t.Errorf("expected worktree name in output, got %q", result)
	}
}

func TestWorktree_EmptyWhenNil(t *testing.T) {
	data := &schema.Input{}
	result := Get("worktree").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when no worktree, got %q", result)
	}
}

func TestWorktree_FallsBackToBranchWhenNameEmpty(t *testing.T) {
	data := &schema.Input{Worktree: &schema.Worktree{Name: "", Branch: "feature/foo"}}
	result := Get("worktree").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "feature/foo") {
		t.Errorf("expected branch as fallback, got %q", result)
	}
}

func TestWorktree_EmptyWhenBothNameAndBranchEmpty(t *testing.T) {
	data := &schema.Input{Worktree: &schema.Worktree{Name: "", Branch: ""}}
	result := Get("worktree").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output, got %q", result)
	}
}
