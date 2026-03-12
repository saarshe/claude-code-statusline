package components

import (
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
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

	return th.Primary.Render(EmojiPrefix(cfg, "🌿", "") + label)
}
