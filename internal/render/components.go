package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// buildBorderLine accepts color as a parameter rather than using a pre-configured style
// to support vertical gradients where each line needs a different color based on its
// position in the box. The totalWidth calculation accounts for contentPadding on both sides.
func buildBorderLine(border lipgloss.Border, width int, color string, leftCorner, line, rightCorner string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	totalWidth := width + contentPadding*2
	return style.Render(leftCorner + strings.Repeat(line, totalWidth) + rightCorner)
}

// buildSideBorders accepts separate left/right colors to support vertical gradients
// where the border color changes from top to bottom of the box. Currently both sides
// use the same color, but the API supports asymmetric gradients for future enhancements.
func buildSideBorders(border lipgloss.Border, width int, leftColor, rightColor string, content string) string {
	leftStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(leftColor))
	rightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(rightColor))
	return leftStyle.Render(border.Left) + content + rightStyle.Render(border.Right)
}

// buildHeaderLine implements a dual-gradient system: horizontal gradient in the slash
// background and vertical gradient positioning via sideColor. The prefix uses the gradient
// start color to create a seamless transition into the horizontal gradient rather than
// starting with gray (previous design had visual discontinuity).
func buildHeaderLine(border lipgloss.Border, text string, width int, gradient []string, sideColor string, textStyle lipgloss.Style) string {
	if text == "" {
		return ""
	}

	prefix := "╱╱ "
	suffix := " "

	startColor := getGradientColorAt(gradient, 0)
	firstColorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(startColor))
	styledPrefix := firstColorStyle.Render(prefix)
	textWithPadding := styledPrefix + textStyle.Render(text) + suffix

	textWidth := lipgloss.Width(textWithPadding)
	totalWidth := width + contentPadding*2

	if textWidth >= totalWidth {
		truncated := truncateText(text, totalWidth-lipgloss.Width(prefix+suffix))
		textWithPadding = styledPrefix + textStyle.Render(truncated) + suffix
		textWidth = lipgloss.Width(textWithPadding)
	}

	slashCount := totalWidth - textWidth

	var gradientSlashes strings.Builder
	for i := 0; i < slashCount; i++ {
		percentage := float64(i) / float64(slashCount)
		color := getGradientColorAt(gradient, percentage)
		colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		gradientSlashes.WriteString(colorStyle.Render("╱"))
	}

	content := textWithPadding + gradientSlashes.String()

	return buildSideBorders(border, width, sideColor, sideColor, content)
}

// buildFooterLine uses gray (color 240) for all slashes rather than the gradient colors,
// creating visual hierarchy where the header is prominent and the footer is subdued.
// This design choice helps users focus on the header (typically status/title) while
// keeping footer metadata available but not dominant.
func buildFooterLine(border lipgloss.Border, text string, width int, sideColor string) string {
	if text == "" {
		return ""
	}

	grayStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	prefix := "╱╱ "
	suffix := " "
	textWithPadding := prefix + text + suffix

	textWidth := lipgloss.Width(textWithPadding)
	totalWidth := width + contentPadding*2

	if textWidth >= totalWidth {
		truncated := truncateText(text, totalWidth-lipgloss.Width(prefix+suffix))
		textWithPadding = prefix + truncated + suffix
		textWidth = lipgloss.Width(textWithPadding)
	}

	slashCount := totalWidth - textWidth
	slashes := strings.Repeat("╱", slashCount)

	content := grayStyle.Render(textWithPadding + slashes)

	return buildSideBorders(border, width, sideColor, sideColor, content)
}
