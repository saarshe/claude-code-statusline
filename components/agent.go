package components

import (
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

type agentComponent struct{}

func init() { Register(&agentComponent{}) }

func (a *agentComponent) Key() ComponentKey { return "agent" }

func (a *agentComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.Agent == nil || data.Agent.Name == "" {
		return ""
	}
	return th.Primary.Render(GetMeta(a.Key()).Prefix(cfg) + data.Agent.Name)
}
