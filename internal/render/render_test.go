package render

import (
	"testing"

	"boxed/internal/box"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRenderer struct {
	rendered []string
}

func (m *mockRenderer) RenderBox(b *box.Box) string {
	output := b.Type.String() + ": " + b.Title
	m.rendered = append(m.rendered, output)
	return output
}

func TestMockRenderer(t *testing.T) {
	mock := &mockRenderer{}
	b := &box.Box{
		Type:  box.Success,
		Title: "Test",
	}

	result := mock.RenderBox(b)

	assert.Equal(t, "success: Test", result)
	assert.Len(t, mock.rendered, 1)
}

func TestNewLipGlossRenderer(t *testing.T) {
	renderer := NewLipGlossRenderer()

	require.NotNil(t, renderer)
}

func TestLipGlossRenderer_RenderBox(t *testing.T) {
	tests := []struct {
		name     string
		box      *box.Box
		contains []string
	}{
		{
			name: "success box with title",
			box: &box.Box{
				Type:  box.Success,
				Title: "Deploy Complete",
			},
			contains: []string{"Deploy Complete"},
		},
		{
			name: "error box with all fields",
			box: &box.Box{
				Type:     box.Error,
				Title:    "Build Failed",
				Subtitle: "v1.2.3",
				KVPairs: []box.KV{
					{Key: "exit", Value: "1"},
					{Key: "duration", Value: "5m"},
				},
				Footer: "Check logs",
			},
			contains: []string{
				"Build Failed",
				"v1.2.3",
				"exit",
				": 1",
				"duration",
				": 5m",
				"Check logs",
			},
		},
		{
			name: "info box with KV pairs",
			box: &box.Box{
				Type:  box.Info,
				Title: "Status",
				KVPairs: []box.KV{
					{Key: "env", Value: "prod"},
					{Key: "region", Value: "us-east-1"},
				},
			},
			contains: []string{
				"Status",
				"env",
				": prod",
				"region",
				": us-east-1",
			},
		},
		{
			name: "warning box with custom width",
			box: &box.Box{
				Type:  box.Warning,
				Title: "Deprecated",
				Width: 50,
			},
			contains: []string{"Deprecated"},
		},
		{
			name: "box with different border styles",
			box: &box.Box{
				Type:        box.Success,
				Title:       "Test",
				BorderStyle: "thick",
			},
			contains: []string{"Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewLipGlossRenderer()

			result := renderer.RenderBox(tt.box)

			require.NotEmpty(t, result)
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected,
					"rendered output should contain %q", expected)
			}
		})
	}
}

func TestLipGlossRenderer_MultiCellCharacters(t *testing.T) {
	tests := []struct {
		name string
		box  *box.Box
	}{
		{
			name: "emoji in title",
			box: &box.Box{
				Type:  box.Success,
				Title: "Deploy Complete ✅",
			},
		},
		{
			name: "CJK characters",
			box: &box.Box{
				Type:  box.Info,
				Title: "部署完了",
				KVPairs: []box.KV{
					{Key: "環境", Value: "本番"},
				},
			},
		},
		{
			name: "mixed unicode",
			box: &box.Box{
				Type:     box.Warning,
				Title:    "Déploiement 🚀",
				Subtitle: "Versión 2.0",
				Footer:   "Готово ✓",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewLipGlossRenderer()

			result := renderer.RenderBox(tt.box)

			require.NotEmpty(t, result)
			assert.Contains(t, result, tt.box.Title)
		})
	}
}
