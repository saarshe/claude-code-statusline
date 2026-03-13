package theme

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Renderer with forced ANSI256 for piped stdout.
// See render/render.go for explanation of WithUnsafe().
var Renderer = lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.ANSI256), termenv.WithUnsafe())

type Theme struct {
	Name      string
	Primary   lipgloss.Style
	Secondary lipgloss.Style
	Accent    lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Danger    lipgloss.Style
	Muted     lipgloss.Style

	// Separator is the string placed between components (including any spacing).
	// If empty, the renderer falls back to the config separator with Muted styling.
	// Examples: " | ", " · ", "", "".
	Separator string
}

var registry = map[string]*Theme{}

func Register(t *Theme) {
	registry[t.Name] = t
}

func Get(name string) *Theme {
	if t, ok := registry[name]; ok {
		return t
	}
	return registry["default"]
}

func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
