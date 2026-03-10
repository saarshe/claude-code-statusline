package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	Register(&Theme{
		Name:      "default",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("6")),  // cyan
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("7")),  // white
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("5")),  // magenta
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("2")),  // green
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("3")),  // yellow
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("1")),  // red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("8")),  // gray
	})
}
