package components

import (
	"os/exec"
	"strings"
	// exec is used by execInDir helper

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func execInDir(dir string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd
}

type gitBranchComponent struct{}

func init() { Register(&gitBranchComponent{}) }

func (g *gitBranchComponent) Key() ComponentKey { return "git_branch" }

func (g *gitBranchComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	branch := gitBranch(data.WorkDir())
	if branch == "" {
		return ""
	}
	return th.Accent.Render(GetMeta(g.Key()).Prefix(cfg) + branch)
}

func gitBranch(dir string) string {
	cmd := execInDir(dir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
