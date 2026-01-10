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

var (
	FooterHeight       = 1
	ExpandedHelpHeight = 4 
)

func SplitHeightByPercentage(height int, percentages []float64, padding, borderHeight int) []int {
	numRows := len(percentages)
	if numRows == 0 {
		return []int{}
	}
	var sum float64
	for _, p := range percentages {
		sum += p
	}
	totalBorderHeight := borderHeight * numRows
	totalPadding := padding * numRows
	usableHeight := height - totalPadding - totalBorderHeight
	allocatedHeight := 0
	heights := make([]int, numRows)
	for i, p := range percentages {
		normalizedPercent := p / sum
		heights[i] = int(float64(usableHeight) * normalizedPercent)
		allocatedHeight += heights[i]
	}
	// Distribute remainder to last row
	remainder := usableHeight - allocatedHeight
	heights[len(heights)-1] += remainder
	return heights
}

func SplitWidthByPercentage(width int, percentages []float64, padding, borderWidth int) []int {
	numCols := len(percentages)
	if numCols == 0 {
		return []int{}
	}
	var sum float64
	for _, p := range percentages {
		sum += p
	}
	totalBorderWidth := borderWidth * numCols
	totalPadding := padding * numCols
	usableWidth := width - totalPadding - totalBorderWidth

	allocatedWidth := 0
	widths := make([]int, numCols)
	for i, p := range percentages {
		normalizedPercent := p / sum
		widths[i] = int(float64(usableWidth) * normalizedPercent)
		allocatedWidth += widths[i]
	}

	// Distribute remainder to last column
	remainder := usableWidth - allocatedWidth
	widths[len(widths)-1] += remainder

	return widths
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
