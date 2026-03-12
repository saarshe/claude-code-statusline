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
	dir := data.WorkDir()
	if dir == "" {
		return ""
	}

	base := filepath.Base(dir)
	if base == "" || base == "." {
		return ""
	}

	return th.Primary.Render(GetMeta(d.Key()).Prefix(cfg) + base)
}
