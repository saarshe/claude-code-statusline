package components

import (
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type modelComponent struct{}

func init() { Register(&modelComponent{}) }

func (m *modelComponent) Key() ComponentKey { return "model" }

func (m *modelComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	name := data.Model.DisplayName
	if name == "" {
		return ""
	}
	meta := GetMeta(m.Key())
	return th.Primary.Render(meta.Prefix(cfg) + name + meta.Suffix(cfg))
}
