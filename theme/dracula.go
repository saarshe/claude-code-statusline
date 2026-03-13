package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Dracula palette
	// https://draculatheme.com
	Register(&Theme{
		Name:      "dracula",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("#8be9fd")), // cyan
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("#f8f8f2")), // foreground
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("#bd93f9")), // purple
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("#50fa7b")), // green
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("#f1fa8c")), // yellow
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("#ff5555")), // red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#6272a4")), // comment
	})
}
