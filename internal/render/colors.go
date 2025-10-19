package render

import (
	"boxed/internal/box"

	"github.com/charmbracelet/lipgloss/v2"
)

// getGradientColorAt implements percentage-based sampling from discrete color arrays.
// Uses floor-based indexing rather than interpolation since ANSI 256-color codes are
// discrete values that can't be blended. Clamping prevents panics from caller math errors.
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

// getColorForType uses ANSI 256-color codes rather than RGB hex values because Lip Gloss
// handles terminal capability detection automatically. The library degrades colors based
// on COLORTERM and TERM environment variables, making this more portable than trying to
// implement our own color downsampling.
func (r *LipGlossRenderer) getColorForType(t box.BoxType) string {
	switch t {
	case box.Success:
		return "114"
	case box.Error:
		return "210"
	case box.Info:
		return "111"
	case box.Warning:
		return "179"
	default:
		return "7"
	}
}

// getGradientForType defines color progressions that transition between related hues
// rather than brightness levels. This aesthetic choice emerged from user feedback that
// dark-to-bright gradients felt too dramatic. Gradient arrays are intentionally sparse
// (8 colors) to avoid per-character rendering overhead while maintaining smooth transitions.
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

// getBorderStyle defers border selection to render time rather than parse time to maintain
// separation between the box package (pure data model) and render package (presentation logic).
// This keeps Lip Gloss as a dependency only in the render layer.
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
