package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Powerline — colored background segments with  arrow separators.
	// Inspired by classic powerline / powerlevel10k "Classic" style.
	bg := func(fg, bg string) lipgloss.Style {
		return Renderer.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(bg)).
			Padding(0, 1)
	}

	Register(&Theme{
		Name:      "powerline",
		Primary:   bg("#ffffff", "#4e6bbd"), // white on blue
		Secondary: bg("#ffffff", "#555555"), // white on gray
		Accent:    bg("#ffffff", "#8a5bbd"), // white on purple
		Success:   bg("#ffffff", "#2d8b46"), // white on green
		Warning:   bg("#000000", "#c4a53e"), // black on yellow
		Danger:    bg("#ffffff", "#c23b3b"), // white on red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#555555")),
		Separator: "",
	})
}
