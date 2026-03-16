package components

import (
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type separatorComponent struct{}

func init() { Register(&separatorComponent{}) }

func (s *separatorComponent) Key() ComponentKey { return "separator" }

func (s *separatorComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	return th.Muted.Render(cfg.Separator.Character)
}
