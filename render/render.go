package render

import (
	"strings"

	"github.com/saars/claude-code-statusline/components"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func Render(input *schema.Input, cfg *config.Config) string {
	th := theme.Get(cfg.Theme)
	return RenderWithTheme(input, cfg, th)
}

func RenderWithTheme(input *schema.Input, cfg *config.Config, th *theme.Theme) string {
	var sep string
	if th.Separator != "" {
		sep = th.Muted.Render(th.Separator)
	} else {
		sep = th.Muted.Render(" " + cfg.Separator.Character + " ")
	}

	lineOutputs := []string{}
	for _, line := range cfg.Lines {
		parts := []string{}
		for _, key := range line.Components {
			c := components.Get(key)
			if c == nil {
				continue
			}
			result := c.Render(input, cfg, th)
			if result != "" {
				parts = append(parts, result)
			}
		}
		if len(parts) > 0 {
			lineOutputs = append(lineOutputs, strings.Join(parts, sep))
		}
	}

	return strings.Join(lineOutputs, "\n")
}
