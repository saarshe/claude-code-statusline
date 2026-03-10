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
	if cfg.Emojis != config.EmojiNone {
		return th.Primary.Render("🤖 " + name)
	}
	return th.Primary.Render("[" + name + "]")
}
