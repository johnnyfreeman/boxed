package render

import (
	"strings"

	"boxed/internal/box"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	// maxLineWidth is the maximum width for text content before wrapping occurs.
	// This prevents extremely wide boxes and ensures readability in standard terminals.
	maxLineWidth = 100

	// contentPadding is the number of spaces on each side of content lines (left and right).
	// Total horizontal space added to content width = contentPadding * 2.
	contentPadding = 3
)

// Renderer abstracts the box rendering logic, enabling dependency injection
// for testing. Production code uses LipGlossRenderer, while tests can inject
// mock renderers to verify behavior without terminal dependencies.
type Renderer interface {
	RenderBox(b *box.Box) string
}

// LipGlossRenderer implements Renderer using Lip Gloss v2 for terminal styling.
// Unlike v1, Lip Gloss v2 uses a simpler API without explicit renderer objects,
// relying on package-level functions and the global Writer for output detection.
type LipGlossRenderer struct{}

// NewLipGlossRenderer creates a renderer that uses Lip Gloss v2's global
// configuration. Color detection and terminal capability detection happen
// automatically via the lipgloss.Writer which reads from os.Stdout and os.Environ.
func NewLipGlossRenderer() *LipGlossRenderer {
	return &LipGlossRenderer{}
}

// RenderBox converts a Box model into a styled terminal string using Lip Gloss.
// The layout embeds title/subtitle in the top border and footer in the bottom border,
// with KV pairs in the content area. Keys are dimmed while values remain at full brightness.
// Auto-sizing calculates the minimum width needed to display all content without wrapping.
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

// processKVPairs formats KV pairs with wrapping and indentation, returning the
// formatted lines and the maximum width encountered. Keys are rendered with the
// provided keyStyle, and values are wrapped at word boundaries to fit within maxLineWidth.
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

// buildHeaderText constructs the header text from title and subtitle with appropriate styling.
// The title is bold and colored, the subtitle is italic and faint. If only one is provided,
// that one is returned. If both are provided, they're combined with a space separator.
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

// calculateBoxWidth determines the final box width based on content, header, footer,
// and user-specified width. The box is sized to fit the widest element, with a maximum
// cap on footer width to prevent extremely wide boxes. If the user specifies a width
// that's larger than the minimum needed, that width is used.
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

func getGradientColorAt(gradient []string, percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}

	colorIndex := int(percentage * float64(len(gradient)-1))
	if colorIndex >= len(gradient) {
		colorIndex = len(gradient) - 1
	}
	return gradient[colorIndex]
}

func buildBorderLine(border lipgloss.Border, width int, color string, leftCorner, line, rightCorner string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	totalWidth := width + contentPadding*2
	return style.Render(leftCorner + strings.Repeat(line, totalWidth) + rightCorner)
}

func buildSideBorders(border lipgloss.Border, width int, leftColor, rightColor string, content string) string {
	leftStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(leftColor))
	rightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(rightColor))
	return leftStyle.Render(border.Left) + content + rightStyle.Render(border.Right)
}

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

// getColorForType maps box types to ANSI color codes appropriate for
// each semantic meaning. Uses 8-bit ANSI codes for broad terminal compatibility
// while still providing clear visual distinction. The color profile detection
// in Lip Gloss will automatically degrade these to available colors.
func (r *LipGlossRenderer) getColorForType(t box.BoxType) string {
	switch t {
	case box.Success:
		return "114" // Tokyo Night: bright lime green (#9ece6a)
	case box.Error:
		return "210" // Tokyo Night: soft pink-red (#f7768e)
	case box.Info:
		return "111" // Tokyo Night: bright sky blue (#7aa2f7)
	case box.Warning:
		return "179" // Tokyo Night: warm golden orange (#e0af68)
	default:
		return "7"
	}
}

func (r *LipGlossRenderer) getGradientForType(t box.BoxType) []string {
	switch t {
	case box.Success:
		return []string{"114", "120", "156", "157", "158", "122", "86", "50"}
	case box.Error:
		return []string{"210", "211", "217", "218", "219", "225", "224", "223"}
	case box.Info:
		return []string{"111", "117", "153", "189", "225", "219", "213", "177"}
	case box.Warning:
		return []string{"179", "215", "221", "227", "228", "229", "223", "217"}
	default:
		return []string{"238", "240", "242", "244", "246", "248", "250"}
	}
}

// getBorderStyle maps user-friendly border style names to Lip Gloss Border structs.
// Empty string defaults to rounded borders as a sensible aesthetic choice for
// modern terminals. The style selection happens at render time rather than parse time
// to keep the box package free of Lip Gloss dependencies.
func (r *LipGlossRenderer) getBorderStyle(style string) lipgloss.Border {
	switch style {
	case "normal":
		return lipgloss.NormalBorder()
	case "rounded":
		return lipgloss.RoundedBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "double":
		return lipgloss.DoubleBorder()
	default:
		return lipgloss.RoundedBorder()
	}
}

// wrapText wraps text at word boundaries to fit within the specified width.
// For words longer than maxWidth (like file paths or URLs), it breaks them at
// special characters (/, _, -, .) or by character count as a last resort.
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for i, word := range words {
		wordWidth := lipgloss.Width(word)

		if wordWidth > maxWidth {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}

			wrappedWord := breakLongWord(word, maxWidth)
			wordLines := strings.Split(wrappedWord, "\n")
			for j, wl := range wordLines {
				if j == len(wordLines)-1 && i < len(words)-1 {
					currentLine.WriteString(wl)
				} else {
					lines = append(lines, wl)
				}
			}
			continue
		}

		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else {
			lineWithWord := currentLine.String() + " " + word
			if lipgloss.Width(lineWithWord) <= maxWidth {
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
			} else {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
			}
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

// breakLongWord breaks a single long word into multiple lines at maxWidth.
// Prefers breaking at special characters (/, _, -, .) for better readability.
func breakLongWord(word string, maxWidth int) string {
	if lipgloss.Width(word) <= maxWidth {
		return word
	}

	var lines []string
	var currentChunk strings.Builder
	runes := []rune(word)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		currentChunk.WriteRune(r)

		if lipgloss.Width(currentChunk.String()) >= maxWidth {
			chunkStr := currentChunk.String()

			lastBreak := strings.LastIndexAny(chunkStr, "/_-.")
			if lastBreak > 0 && lastBreak < len(chunkStr)-1 {
				lines = append(lines, chunkStr[:lastBreak+1])
				currentChunk.Reset()
				currentChunk.WriteString(chunkStr[lastBreak+1:])
			} else {
				lines = append(lines, chunkStr)
				currentChunk.Reset()
			}
		}
	}

	if currentChunk.Len() > 0 {
		lines = append(lines, currentChunk.String())
	}

	return strings.Join(lines, "\n")
}

// truncateText truncates text to fit within maxWidth, adding an ellipsis if truncated.
// This is used for text embedded in borders (title/subtitle/footer) that can't wrap to multiple lines.
func truncateText(text string, maxWidth int) string {
	if maxWidth <= 3 {
		return "..."
	}

	textWidth := lipgloss.Width(text)
	if textWidth <= maxWidth {
		return text
	}

	// Reserve 3 characters for ellipsis
	targetWidth := maxWidth - 3
	runes := []rune(text)
	var currentWidth int
	var cutoff int

	for i, r := range runes {
		runeWidth := lipgloss.Width(string(r))
		if currentWidth+runeWidth > targetWidth {
			cutoff = i
			break
		}
		currentWidth += runeWidth
	}

	if cutoff == 0 {
		return "..."
	}

	return string(runes[:cutoff]) + "..."
}
