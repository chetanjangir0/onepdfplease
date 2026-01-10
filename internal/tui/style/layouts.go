package style

import (
	"github.com/charmbracelet/lipgloss"
)

func RenderTwoFullCols(termWidth, Height int, style lipgloss.Style, col1View, col2View string) string {
	padding := 0
	height := SplitHeightByPercentage(Height, []float64{1}, padding, 2)
	widths := SplitWidthByPercentage(termWidth, []float64{0.5, 0.5}, padding, 2)

	// truncateView
	col1View = TruncateView(col1View, widths[0])
	col2View = TruncateView(col2View, widths[1])

	// apply style
	col1View = style.
		Padding(padding).
		Width(widths[0]).
		Height(height[0]).
		Render(col1View)
	col2View = style.
		Padding(padding).
		Width(widths[1]).
		Height(height[0]).
		Render(col2View)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		col1View,
		col2View,
	)
}

func RenderTwoFullRows(termWidth, termHeight int, row1Style, row2Style lipgloss.Style, row1View, row2View string) string {
	padding := 0
	heights := SplitHeightByPercentage(termHeight, []float64{0.6, 0.4}, padding, 2)
	width := SplitWidthByPercentage(termWidth, []float64{1}, padding, 2)

	// truncateView
	row1View = TruncateView(row1View, width[0])
	row2View = TruncateView(row2View, width[0])

	// apply style
	row1View = row1Style.
		Padding(padding).
		Width(width[0]).
		Height(heights[0]).
		Render(row1View)
	row2View = row2Style.
		Padding(padding).
		Width(width[0]).
		Height(heights[1]).
		Render(row2View)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		row1View,
		row2View,
	)
}
