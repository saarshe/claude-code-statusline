package components

import (
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type sessionIDComponent struct{}

func init() { Register(&sessionIDComponent{}) }

func (s *sessionIDComponent) Key() ComponentKey { return "session_id" }

func (s *sessionIDComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.SessionID == "" {
		return ""
	}

	return th.Primary.Render(GetMeta(s.Key()).Prefix(cfg) + data.SessionID)
}
