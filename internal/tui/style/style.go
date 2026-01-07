package style

import (
	"path/filepath"
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

func RenderColumnLayout(termWidth, height int, style lipgloss.Style, columnViews ...string) string {
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
		colWidth := columnContentWidth
		if i == numColumns-1 {
			colWidth += remainder
		}
		colContent := TruncateView(columnViews[i], colWidth)

		columns[i] = style.
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
		lines[i] = TruncateFilenameMiddle(line, length, 4)
	}
	return strings.Join(lines, "\n")
}

func TruncateFilenameMiddle(name string, maxWidth int, keep int) string {
	const tail = "..."

	// If it already fits, return as-is
	if ansi.StringWidth(name) <= maxWidth {
		return name
	}

	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	// if no extension or too small to bother
	if ext == "" || maxWidth <= ansi.StringWidth(ext)+len(tail)+keep {
		return ansi.Truncate(name, maxWidth, tail)
	}

	// Take last `keep` chars from base
	runes := []rune(base)
	end := string(runes[max(len(runes)-keep, 0):])

	suffix := end + ext

	// Remaining space for prefix
	remaining := maxWidth - ansi.StringWidth(suffix) - ansi.StringWidth(tail)

	if remaining <= 0 {
		return ansi.Truncate(name, maxWidth, tail)
	}

	prefix := ansi.Truncate(base, remaining, "")
	return prefix + tail + suffix
}
