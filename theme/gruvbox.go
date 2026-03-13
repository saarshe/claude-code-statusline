package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Gruvbox Dark palette
	// https://github.com/morhetz/gruvbox
	Register(&Theme{
		Name:      "gruvbox",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("#83a598")), // blue
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("#ebdbb2")), // fg
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("#d3869b")), // purple
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("#b8bb26")), // green
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("#fabd2f")), // yellow
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("#fb4934")), // red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#928374")), // gray
	})
}
