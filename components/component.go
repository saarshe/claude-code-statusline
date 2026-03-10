package components

import (
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

// ComponentKey is the string identifier used in config (e.g. "model", "cost").
type ComponentKey string

// Component renders a single piece of status line output.
type Component interface {
	Key() ComponentKey
	// Render returns the formatted string. Returns empty string if nothing to show.
	Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string
}

var registry = map[ComponentKey]Component{}

func Register(c Component) {
	registry[c.Key()] = c
}

func Get(key string) Component {
	return registry[ComponentKey(key)]
}

func registeredKeys() []string {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, string(k))
	}
	return keys
}
