package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type Style struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style
}

var DefaultStyle = &Style{
	FocusedBorder: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")), // Bright pink/magenta

	BlurredBorder: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")), // Dim gray
}

func RenderColumnLayout(termWidth, height int, columnViews ...string) string {
	numColumns := len(columnViews)
	spacing := 0 // Space between columns

	totalSpacing := spacing * (numColumns - 1)

	borderWidthPerColumn := 2
	totalBorderWidth := borderWidthPerColumn * numColumns

	paddingPerColumn := 0
	totalPadding := paddingPerColumn * numColumns

	usableWidth := termWidth - totalBorderWidth - totalPadding - totalSpacing

	// Each column gets equal percentage
	columnContentWidth := usableWidth / numColumns

	// Adjust for any remainder pixels
	remainder := usableWidth % numColumns

	columnHeight := height - paddingPerColumn - borderWidthPerColumn

	columns := make([]string, numColumns)
	for i := range numColumns {
		colContent := lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render(columnViews[i]))

		colWidth := columnContentWidth
		if i == numColumns-1 {
			colWidth += remainder
		}
		columns[i] = DefaultStyle.FocusedBorder.
			Padding(paddingPerColumn).
			Width(colWidth).
			Height(columnHeight).
			Render(colContent)
	}

	// Join columns horizontally with spacing
	spacer := strings.Repeat(" ", spacing)
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		addSpacerInBetween(columns, spacer)...,
	)

	return row
}

func addSpacerInBetween(cols []string, spacer string) []string {
	if len(cols) == 0 {
		return nil
	}

	out := make([]string, 0, len(cols)*2-1)
	for i, c := range cols {
		if i > 0 {
			out = append(out, spacer)
		}
		out = append(out, c)
	}
	return out
}

func TruncateView(view string, length int) string {

	lines := strings.Split(view, "\n")
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, length, "...")
	}
	return strings.Join(lines, "\n")
}
