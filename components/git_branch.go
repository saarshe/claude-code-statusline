package components

import (
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type gitBranchComponent struct{}

func init() { Register(&gitBranchComponent{}) }

func (g *gitBranchComponent) Key() ComponentKey { return "git_branch" }

func (g *gitBranchComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	branch := data.Git.Branch
	if branch == "" {
		return ""
	}
	return th.Accent.Render(GetMeta(g.Key()).Prefix(cfg) + branch)
}
