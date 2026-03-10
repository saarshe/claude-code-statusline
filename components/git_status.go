package components

import (
	"fmt"
	"strings"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type gitStatusComponent struct{}

func init() { Register(&gitStatusComponent{}) }

func (g *gitStatusComponent) Key() ComponentKey { return "git_status" }

func (g *gitStatusComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	dir := data.Workspace.CurrentDir
	if dir == "" {
		dir = data.Cwd
	}

	branch := gitBranch(dir)
	if branch == "" {
		return ""
	}

	staged, modified := gitCounts(dir)

	prefix := ""
	if cfg.Emojis != config.EmojiNone {
		prefix = "🌿 "
	}

	parts := []string{branch}
	if staged > 0 {
		parts = append(parts, fmt.Sprintf("+%d", staged))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d", modified))
	}

	return th.Accent.Render(prefix + strings.Join(parts, " "))
}

func gitCounts(dir string) (staged, modified int) {
	cmd := execInDir(dir, "git", "status", "--porcelain")
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
