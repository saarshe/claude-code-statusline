package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Rounded — colored background segments with  rounded separators.
	// Softer take on powerline with pastel-ish backgrounds.
	bg := func(fg, bg string) lipgloss.Style {
		return Renderer.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(bg)).
			Padding(0, 1)
	}

	Register(&Theme{
		Name:      "rounded",
		Primary:   bg("#ffffff", "#5c82db"), // white on soft blue
		Secondary: bg("#ffffff", "#666666"), // white on mid gray
		Accent:    bg("#ffffff", "#9b6fc3"), // white on soft purple
		Success:   bg("#ffffff", "#3a9a5b"), // white on soft green
		Warning:   bg("#000000", "#d4b44e"), // black on soft yellow
		Danger:    bg("#ffffff", "#d14d4d"), // white on soft red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#666666")),
		Separator: "",
	})
}
