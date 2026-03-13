package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Tokyo Night palette
	// https://github.com/enkia/tokyo-night-vscode-theme
	Register(&Theme{
		Name:      "tokyo-night",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("#7aa2f7")), // blue
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("#a9b1d6")), // fg
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("#bb9af7")), // purple
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("#9ece6a")), // green
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("#e0af68")), // yellow
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("#f7768e")), // red
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#565f89")), // comment
	})
}
