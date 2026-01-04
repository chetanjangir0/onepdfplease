package style

import "github.com/charmbracelet/lipgloss"


type Style struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style

}

var DefaultStyle = &Style{
	FocusedBorder : lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")), // Bright pink/magenta

	BlurredBorder : lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")), // Dim gray
}
