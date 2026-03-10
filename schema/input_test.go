package schema

import (
	"strings"
	"testing"
)

func ptr[T any](v T) *T {
	return &v
}

func TestParse_FullInput(t *testing.T) {
	json := `{
		"cwd": "/home/user/project",
		"session_id": "abc123",
		"transcript_path": "/tmp/transcript.json",
		"version": "1.0.0",
		"model": {
			"id": "claude-opus-4-6",
			"display_name": "Opus"
		},
		"workspace": {
			"current_dir": "/home/user/project",
			"project_dir": "/home/user/project"
		},
		"cost": {
			"total_cost_usd": 0.42,
			"total_duration_ms": 754000,
			"total_api_duration_ms": 500000,
			"total_lines_added": 156,
			"total_lines_removed": 23
		},
		"context_window": {
			"total_input_tokens": 50000,
			"total_output_tokens": 5000,
			"context_window_size": 200000,
			"used_percentage": 28.0,
			"remaining_percentage": 72.0,
			"current_usage": {
				"input_tokens": 8500,
				"output_tokens": 1200,
				"cache_creation_input_tokens": 2000,
				"cache_read_input_tokens": 5000
			}
		},
		"exceeds_200k_tokens": false,
		"output_style": {"name": "default"},
		"vim": {"mode": "normal"},
		"agent": {"name": "code-reviewer"},
		"worktree": {
			"name": "feature-branch",
			"path": "/tmp/worktree",
			"branch": "feature-branch",
			"original_cwd": "/home/user/project",
			"original_branch": "main"
		}
	}`

	input, err := Parse(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if input.Cwd != "/home/user/project" {
		t.Errorf("Cwd = %q, want %q", input.Cwd, "/home/user/project")
	}
	if input.SessionID != "abc123" {
		t.Errorf("SessionID = %q, want %q", input.SessionID, "abc123")
	}
	if input.TranscriptPath != "/tmp/transcript.json" {
		t.Errorf("TranscriptPath = %q, want %q", input.TranscriptPath, "/tmp/transcript.json")
	}
	if input.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", input.Version, "1.0.0")
	}
	if input.Model.ID != "claude-opus-4-6" {
		t.Errorf("Model.ID = %q, want %q", input.Model.ID, "claude-opus-4-6")
	}
	if input.Model.DisplayName != "Opus" {
		t.Errorf("Model.DisplayName = %q, want %q", input.Model.DisplayName, "Opus")
	}
	if input.Workspace.CurrentDir != "/home/user/project" {
		t.Errorf("Workspace.CurrentDir = %q, want %q", input.Workspace.CurrentDir, "/home/user/project")
	}
	if input.Cost.TotalCostUSD != 0.42 {
		t.Errorf("Cost.TotalCostUSD = %f, want %f", input.Cost.TotalCostUSD, 0.42)
	}
	if input.Cost.TotalDurationMS != 754000 {
		t.Errorf("Cost.TotalDurationMS = %d, want %d", input.Cost.TotalDurationMS, int64(754000))
	}
	if input.Cost.TotalLinesAdded != 156 {
		t.Errorf("Cost.TotalLinesAdded = %d, want %d", input.Cost.TotalLinesAdded, 156)
	}
	if input.Cost.TotalLinesRemoved != 23 {
		t.Errorf("Cost.TotalLinesRemoved = %d, want %d", input.Cost.TotalLinesRemoved, 23)
	}
	if input.ContextWindow.TotalInputTokens != 50000 {
		t.Errorf("ContextWindow.TotalInputTokens = %d, want %d", input.ContextWindow.TotalInputTokens, 50000)
	}
	if input.ContextWindow.ContextWindowSize != 200000 {
		t.Errorf("ContextWindow.ContextWindowSize = %d, want %d", input.ContextWindow.ContextWindowSize, 200000)
	}
	if input.ContextWindow.UsedPercentage == nil || *input.ContextWindow.UsedPercentage != 28.0 {
		t.Errorf("ContextWindow.UsedPercentage = %v, want 28.0", input.ContextWindow.UsedPercentage)
	}
	if input.ContextWindow.RemainingPct == nil || *input.ContextWindow.RemainingPct != 72.0 {
		t.Errorf("ContextWindow.RemainingPct = %v, want 72.0", input.ContextWindow.RemainingPct)
	}
	if input.ContextWindow.CurrentUsage == nil {
		t.Fatal("ContextWindow.CurrentUsage is nil, want non-nil")
	}
	if input.ContextWindow.CurrentUsage.InputTokens != 8500 {
		t.Errorf("CurrentUsage.InputTokens = %d, want %d", input.ContextWindow.CurrentUsage.InputTokens, 8500)
	}
	if input.ContextWindow.CurrentUsage.CacheReadInputTokens != 5000 {
		t.Errorf("CurrentUsage.CacheReadInputTokens = %d, want %d", input.ContextWindow.CurrentUsage.CacheReadInputTokens, 5000)
	}
	if input.Exceeds200k {
		t.Error("Exceeds200k = true, want false")
	}
	if input.OutputStyle == nil || input.OutputStyle.Name != "default" {
		t.Errorf("OutputStyle = %v, want name=default", input.OutputStyle)
	}
	if input.Vim == nil || input.Vim.Mode != "normal" {
		t.Errorf("Vim = %v, want mode=normal", input.Vim)
	}
	if input.Agent == nil || input.Agent.Name != "code-reviewer" {
		t.Errorf("Agent = %v, want name=code-reviewer", input.Agent)
	}
	if input.Worktree == nil || input.Worktree.Name != "feature-branch" {
		t.Errorf("Worktree = %v, want name=feature-branch", input.Worktree)
	}
	if input.Worktree.OriginalBranch != "main" {
		t.Errorf("Worktree.OriginalBranch = %q, want %q", input.Worktree.OriginalBranch, "main")
	}
}

func TestParse_NullOptionalFields(t *testing.T) {
	json := `{
		"cwd": "/home/user/project",
		"session_id": "abc123",
		"transcript_path": "/tmp/transcript.json",
		"version": "1.0.0",
		"model": {"id": "claude-opus-4-6", "display_name": "Opus"},
		"workspace": {"current_dir": "/home/user/project", "project_dir": "/home/user/project"},
		"cost": {"total_cost_usd": 0, "total_duration_ms": 0, "total_api_duration_ms": 0, "total_lines_added": 0, "total_lines_removed": 0},
		"context_window": {
			"total_input_tokens": 0,
			"total_output_tokens": 0,
			"context_window_size": 200000,
			"used_percentage": null,
			"remaining_percentage": null,
			"current_usage": null
		},
		"exceeds_200k_tokens": false
	}`

	input, err := Parse(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if input.ContextWindow.UsedPercentage != nil {
		t.Errorf("UsedPercentage = %v, want nil", input.ContextWindow.UsedPercentage)
	}
	if input.ContextWindow.RemainingPct != nil {
		t.Errorf("RemainingPct = %v, want nil", input.ContextWindow.RemainingPct)
	}
	if input.ContextWindow.CurrentUsage != nil {
		t.Errorf("CurrentUsage = %v, want nil", input.ContextWindow.CurrentUsage)
	}
	if input.OutputStyle != nil {
		t.Errorf("OutputStyle = %v, want nil", input.OutputStyle)
	}
	if input.Vim != nil {
		t.Errorf("Vim = %v, want nil", input.Vim)
	}
	if input.Agent != nil {
		t.Errorf("Agent = %v, want nil", input.Agent)
	}
	if input.Worktree != nil {
		t.Errorf("Worktree = %v, want nil", input.Worktree)
	}
}

func TestParse_MinimalJSON(t *testing.T) {
	input, err := Parse(strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if input.Cwd != "" {
		t.Errorf("Cwd = %q, want empty", input.Cwd)
	}
	if input.Model.DisplayName != "" {
		t.Errorf("Model.DisplayName = %q, want empty", input.Model.DisplayName)
	}
	if input.Cost.TotalCostUSD != 0 {
		t.Errorf("Cost.TotalCostUSD = %f, want 0", input.Cost.TotalCostUSD)
	}
	if input.ContextWindow.CurrentUsage != nil {
		t.Errorf("CurrentUsage = %v, want nil", input.ContextWindow.CurrentUsage)
	}
}

func TestParse_MalformedJSON(t *testing.T) {
	_, err := Parse(strings.NewReader(`{not json`))
	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	_, err := Parse(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty input, got nil")
	}
}

func TestParse_UnknownFields(t *testing.T) {
	json := `{
		"cwd": "/test",
		"some_future_field": "hello",
		"another_unknown": 42,
		"model": {"id": "test", "display_name": "Test", "unknown_nested": true}
	}`

	input, err := Parse(strings.NewReader(json))
	if err != nil {
		t.Fatalf("unexpected error for unknown fields: %v", err)
	}

	if input.Cwd != "/test" {
		t.Errorf("Cwd = %q, want %q", input.Cwd, "/test")
	}
	if input.Model.DisplayName != "Test" {
		t.Errorf("Model.DisplayName = %q, want %q", input.Model.DisplayName, "Test")
	}
}
