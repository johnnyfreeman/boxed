package render

import (
	"strings"

	"boxed/internal/box"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	maxLineWidth   = 100
	contentPadding = 3
)

type Renderer interface {
	RenderBox(b *box.Box) string
}

type LipGlossRenderer struct{}

func NewLipGlossRenderer() *LipGlossRenderer {
	return &LipGlossRenderer{}
}

// RenderBox implements a two-pass layout algorithm: first pass measures all content to
// determine minimum box width, second pass renders each line with gradient colors based on
// vertical position. This avoids re-rendering when the box size changes and separates
// measurement concerns from styling concerns.
func (r *LipGlossRenderer) RenderBox(b *box.Box) string {
	borderColor := r.getColorForType(b.Type)
	border := r.getBorderStyle(b.BorderStyle)
	gradient := r.getGradientForType(b.Type)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(borderColor))
	subtitleStyle := lipgloss.NewStyle().Italic(true).Faint(true)
	keyStyle := lipgloss.NewStyle().Faint(true)

	contentLines, maxContentWidth := processKVPairs(b.KVPairs, keyStyle)
	headerText := buildHeaderText(b.Title, b.Subtitle, titleStyle, subtitleStyle)

	headerWidth := lipgloss.Width(headerText)
	footerWidth := lipgloss.Width(b.Footer)
	contentWidth := calculateBoxWidth(maxContentWidth, headerWidth, footerWidth, b.Width)

	totalLines := 1
	if headerText != "" {
		totalLines++
	}
	if len(contentLines) == 0 {
		totalLines++
	} else {
		totalLines += 2 + len(contentLines)
	}
	if b.Footer != "" {
		totalLines++
	}
	totalLines++

	var lines []string
	lineIndex := 0

	borderColor = getGradientColorAt(gradient, float64(lineIndex)/float64(totalLines-1))
	lines = append(lines, buildBorderLine(border, contentWidth, borderColor, border.TopLeft, border.Top, border.TopRight))
	lineIndex++

	if headerText != "" {
		headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
		percentage := float64(lineIndex) / float64(totalLines-1)
		sideColor := getGradientColorAt(gradient, percentage)
		lines = append(lines, buildHeaderLine(border, headerText, contentWidth, gradient, sideColor, headerStyle))
		lineIndex++
	}

	percentage := float64(lineIndex) / float64(totalLines-1)
	sideColor := getGradientColorAt(gradient, percentage)
	emptyLine := buildSideBorders(border, contentWidth, sideColor, sideColor, strings.Repeat(" ", contentWidth+contentPadding*2))

	if len(contentLines) == 0 {
		lines = append(lines, emptyLine)
		lineIndex++
	} else {
		lines = append(lines, emptyLine)
		lineIndex++

		for _, line := range contentLines {
			padding := contentWidth - lipgloss.Width(line)
			leftPad := strings.Repeat(" ", contentPadding)
			rightPad := strings.Repeat(" ", contentPadding)
			paddedLine := leftPad + line + strings.Repeat(" ", padding) + rightPad
			percentage := float64(lineIndex) / float64(totalLines-1)
			sideColor := getGradientColorAt(gradient, percentage)
			lines = append(lines, buildSideBorders(border, contentWidth, sideColor, sideColor, paddedLine))
			lineIndex++
		}

		percentage = float64(lineIndex) / float64(totalLines-1)
		sideColor = getGradientColorAt(gradient, percentage)
		emptyLine = buildSideBorders(border, contentWidth, sideColor, sideColor, strings.Repeat(" ", contentWidth+contentPadding*2))
		lines = append(lines, emptyLine)
		lineIndex++
	}

	if b.Footer != "" {
		percentage := float64(lineIndex) / float64(totalLines-1)
		sideColor := getGradientColorAt(gradient, percentage)
		lines = append(lines, buildFooterLine(border, b.Footer, contentWidth, sideColor))
		lineIndex++
	}

	borderColor = getGradientColorAt(gradient, float64(lineIndex)/float64(totalLines-1))
	lines = append(lines, buildBorderLine(border, contentWidth, borderColor, border.BottomLeft, border.Bottom, border.BottomRight))

	return strings.Join(lines, "\n")
}

func processKVPairs(kvPairs []box.KV, keyStyle lipgloss.Style) (lines []string, maxWidth int) {
	for _, kv := range kvPairs {
		key := keyStyle.Render(kv.Key)
		keyPrefix := key + ": "
		keyPrefixWidth := lipgloss.Width(keyPrefix)

		wrappedValue := wrapText(kv.Value, maxLineWidth-keyPrefixWidth)
		valueLines := strings.Split(wrappedValue, "\n")

		for i, valueLine := range valueLines {
			var line string
			if i == 0 {
				line = keyPrefix + valueLine
			} else {
				indent := strings.Repeat(" ", keyPrefixWidth)
				line = indent + valueLine
			}
			lines = append(lines, line)
			lineWidth := lipgloss.Width(line)
			if lineWidth > maxWidth {
				maxWidth = lineWidth
			}
		}
	}
	return lines, maxWidth
}

func buildHeaderText(title, subtitle string, titleStyle, subtitleStyle lipgloss.Style) string {
	if title != "" {
		headerText := titleStyle.Render(title)
		if subtitle != "" {
			headerText += " " + subtitleStyle.Render(subtitle)
		}
		return headerText
	} else if subtitle != "" {
		return subtitleStyle.Render(subtitle)
	}
	return ""
}

// calculateBoxWidth caps footer width at maxLineWidth to prevent pathologically wide boxes
// when footers contain long timestamps or file paths. This emerged from testing where
// ISO 8601 timestamps with timezone info created boxes too wide for standard 80-column terminals.
func calculateBoxWidth(contentWidth, headerWidth, footerWidth, requestedWidth int) int {
	// Cap footer width at maxLineWidth to prevent extremely wide boxes
	if footerWidth > maxLineWidth {
		footerWidth = maxLineWidth
	}

	minWidth := contentWidth
	if headerWidth > minWidth {
		minWidth = headerWidth
	}
	if footerWidth > minWidth {
		minWidth = footerWidth
	}

	if requestedWidth > 0 && requestedWidth > minWidth {
		return requestedWidth
	}

	return minWidth
}
