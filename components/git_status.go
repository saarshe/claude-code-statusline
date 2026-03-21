package components

import (
	"fmt"
	"strings"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type gitStatusComponent struct{}

func init() { Register(&gitStatusComponent{}) }

func (g *gitStatusComponent) Key() ComponentKey { return "git_status" }

func (g *gitStatusComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	branch := data.Git.Branch
	if branch == "" {
		return ""
	}

	parts := []string{branch}
	if data.Git.Staged > 0 {
		parts = append(parts, fmt.Sprintf("+%d", data.Git.Staged))
	}
	if data.Git.Modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d", data.Git.Modified))
	}

	return th.Accent.Render(GetMeta(g.Key()).Prefix(cfg) + strings.Join(parts, " "))
}
