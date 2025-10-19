package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// wrapText implements word-boundary wrapping with special handling for long unbreakable
// strings (URLs, file paths). The design assumes terminal output where horizontal scrolling
// isn't available, so we break at special characters or force-break rather than overflowing.
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

// breakLongWord breaks at special characters (/, _, -, .) commonly found in paths and URLs
// rather than arbitrary character positions. This preserves readability by keeping logical
// segments together (e.g., "/usr/local" breaks after slashes, not mid-word).
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

// truncateText handles header/footer text that can't wrap to multiple lines because they're
// embedded in border components. Uses rune iteration rather than byte slicing to correctly
// handle multi-byte Unicode characters (emoji, CJK text).
func truncateText(text string, maxWidth int) string {
	if maxWidth <= 3 {
		return "..."
	}

	textWidth := lipgloss.Width(text)
	if textWidth <= maxWidth {
		return text
	}

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
