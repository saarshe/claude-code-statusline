package components

import (
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type worktreeComponent struct{}

func init() { Register(&worktreeComponent{}) }

func (w *worktreeComponent) Key() ComponentKey { return "worktree" }

func (w *worktreeComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.Worktree == nil {
		return ""
	}

	label := data.Worktree.Name
	if label == "" {
		label = data.Worktree.Branch
	}
	if label == "" {
		return ""
	}

	return th.Primary.Render(GetMeta(w.Key()).Prefix(cfg) + label)
}
