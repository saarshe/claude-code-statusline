package schema

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"
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
	Git            Git       `json:"-"` // populated at runtime, not from JSON
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

// ContextFillTokens returns the number of tokens currently occupying
// the context window, computed from the most recent API call's usage.
// Returns 0 if CurrentUsage is nil (before the first API call).
func (c *Context) ContextFillTokens() int {
	if c.CurrentUsage == nil {
		return 0
	}
	return c.CurrentUsage.TotalInput()
}

// TotalInput returns the total input tokens for this turn, including
// fresh input, cache reads, and cache creation tokens.
func (u *Usage) TotalInput() int {
	return u.InputTokens + u.CacheReadInputTokens + u.CacheCreationInputTokens
}

// CacheHitPct returns the cache hit percentage (0-100). Returns 0 if there
// are no input tokens.
func (u *Usage) CacheHitPct() int {
	total := u.TotalInput()
	if total == 0 {
		return 0
	}
	return int(float64(u.CacheReadInputTokens) / float64(total) * 100)
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

// Git holds pre-fetched git information so components don't need to exec.
type Git struct {
	Branch   string // current branch name (empty if not a git repo)
	Staged   int    // number of staged files
	Modified int    // number of modified (unstaged) files
}

// WorkDir returns the best available working directory, preferring
// Workspace.CurrentDir and falling back to Cwd.
func (i *Input) WorkDir() string {
	if i.Workspace.CurrentDir != "" {
		return i.Workspace.CurrentDir
	}
	return i.Cwd
}

// PopulateGit fetches git branch and status information once and stores it in
// the Git field. Components read from Git instead of spawning subprocesses.
func (i *Input) PopulateGit() {
	dir := i.WorkDir()
	if dir == "" {
		return
	}
	i.Git.Branch = gitBranch(dir)
	if i.Git.Branch == "" {
		return
	}
	i.Git.Staged, i.Git.Modified = gitCounts(dir)
}

func gitBranch(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitCounts(dir string) (staged, modified int) {
	cmd := exec.Command("git", "status", "--porcelain")
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if len(line) < 2 {
			continue
		}
		index := line[0]
		worktree := line[1]
		if index != ' ' && index != '?' {
			staged++
		}
		if worktree != ' ' && worktree != '?' {
			modified++
		}
	}
	return
}

func Parse(r io.Reader) (*Input, error) {
	var input Input
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return nil, err
	}
	return &input, nil
}
