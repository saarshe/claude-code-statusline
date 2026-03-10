package schema

import (
	"encoding/json"
	"io"
)

type Input struct {
	Cwd            string    `json:"cwd"`
	SessionID      string    `json:"session_id"`
	TranscriptPath string    `json:"transcript_path"`
	Version        string    `json:"version"`
	Model          Model     `json:"model"`
	Workspace      Workspace `json:"workspace"`
	Cost           Cost      `json:"cost"`
	ContextWindow  Context   `json:"context_window"`
	Exceeds200k    bool      `json:"exceeds_200k_tokens"`
	OutputStyle    *Style    `json:"output_style,omitempty"`
	Vim            *Vim      `json:"vim,omitempty"`
	Agent          *Agent    `json:"agent,omitempty"`
	Worktree       *Worktree `json:"worktree,omitempty"`
}

type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

type Cost struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMS    int64   `json:"total_duration_ms"`
	TotalAPIDurationMS int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

type Context struct {
	TotalInputTokens  int      `json:"total_input_tokens"`
	TotalOutputTokens int      `json:"total_output_tokens"`
	ContextWindowSize int      `json:"context_window_size"`
	UsedPercentage    *float64 `json:"used_percentage"`
	RemainingPct      *float64 `json:"remaining_percentage"`
	CurrentUsage      *Usage   `json:"current_usage"`
}

type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

type Style struct {
	Name string `json:"name"`
}

type Vim struct {
	Mode string `json:"mode"`
}

type Agent struct {
	Name string `json:"name"`
}

type Worktree struct {
	Name           string `json:"name"`
	Path           string `json:"path"`
	Branch         string `json:"branch,omitempty"`
	OriginalCwd    string `json:"original_cwd"`
	OriginalBranch string `json:"original_branch,omitempty"`
}

func Parse(r io.Reader) (*Input, error) {
	var input Input
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return nil, err
	}
	return &input, nil
}
