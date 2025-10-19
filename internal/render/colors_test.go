package render

import (
	"testing"

	"boxed/internal/box"

	"github.com/stretchr/testify/assert"
)

func TestLipGlossRenderer_GetColorForType(t *testing.T) {
	renderer := NewLipGlossRenderer()

	tests := []struct {
		boxType      box.BoxType
		expectedCode string
	}{
		{box.Success, "114"},
		{box.Error, "210"},
		{box.Info, "111"},
		{box.Warning, "179"},
	}

	for _, tt := range tests {
		t.Run(tt.boxType.String(), func(t *testing.T) {
			color := renderer.getColorForType(tt.boxType)
			assert.Equal(t, tt.expectedCode, color)
		})
	}
}

func TestLipGlossRenderer_GetBorderStyle(t *testing.T) {
	renderer := NewLipGlossRenderer()

	tests := []string{
		"",
		"normal",
		"rounded",
		"thick",
		"double",
	}

	for _, style := range tests {
		t.Run(style, func(t *testing.T) {
			border := renderer.getBorderStyle(style)
			assert.NotEmpty(t, border.Top)
		})
	}
}
