package components

import (
	"path/filepath"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type directoryComponent struct{}

func init() { Register(&directoryComponent{}) }

func (d *directoryComponent) Key() ComponentKey { return "directory" }

func (d *directoryComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	dir := data.Workspace.CurrentDir
	if dir == "" {
		dir = data.Cwd
	}
	if dir == "" {
		return ""
	}

	base := filepath.Base(dir)
	if base == "" || base == "." {
		return ""
	}

	if cfg.Emojis != config.EmojiNone {
		return th.Primary.Render("📁 " + base)
	}
	return th.Primary.Render(base)
}
