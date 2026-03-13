package theme

import "github.com/charmbracelet/lipgloss"

func init() {
	// Nord palette
	// https://www.nordtheme.com
	Register(&Theme{
		Name:      "nord",
		Primary:   Renderer.NewStyle().Foreground(lipgloss.Color("#88c0d0")), // nord8 (frost)
		Secondary: Renderer.NewStyle().Foreground(lipgloss.Color("#d8dee9")), // nord4 (snow storm)
		Accent:    Renderer.NewStyle().Foreground(lipgloss.Color("#b48ead")), // nord15 (aurora purple)
		Success:   Renderer.NewStyle().Foreground(lipgloss.Color("#a3be8c")), // nord14 (aurora green)
		Warning:   Renderer.NewStyle().Foreground(lipgloss.Color("#ebcb8b")), // nord13 (aurora yellow)
		Danger:    Renderer.NewStyle().Foreground(lipgloss.Color("#bf616a")), // nord11 (aurora red)
		Muted:     Renderer.NewStyle().Foreground(lipgloss.Color("#4c566a")), // nord3 (polar night)
	})
}
