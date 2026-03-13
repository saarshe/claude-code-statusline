package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Catppuccin Mocha palette
	// https://github.com/catppuccin/catppuccin
	Register(&Theme{
		Name:      "catppuccin",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("#89b4fa")), // blue
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("#cdd6f4")), // text
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("#cba6f7")), // mauve
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("#a6e3a1")), // green
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("#f9e2af")), // yellow
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("#f38ba8")), // red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#6c7086")), // overlay0
	})
}
