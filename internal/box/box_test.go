package box

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoxType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		boxType  BoxType
		expected bool
	}{
		{"success is valid", Success, true},
		{"error is valid", Error, true},
		{"info is valid", Info, true},
		{"warning is valid", Warning, true},
		{"invalid type", BoxType("invalid"), false},
		{"empty type", BoxType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.boxType.IsValid())
		})
	}
}

func TestBoxType_String(t *testing.T) {
	assert.Equal(t, "success", Success.String())
	assert.Equal(t, "error", Error.String())
	assert.Equal(t, "info", Info.String())
	assert.Equal(t, "warning", Warning.String())
}

func TestAllBoxTypes(t *testing.T) {
	types := AllBoxTypes()
	assert.Len(t, types, 4)
	assert.Contains(t, types, Success)
	assert.Contains(t, types, Error)
	assert.Contains(t, types, Info)
	assert.Contains(t, types, Warning)
}

func TestKV_String(t *testing.T) {
	kv := KV{Key: "env", Value: "prod"}
	assert.Equal(t, "env=prod", kv.String())
}

func TestBox_HasContent(t *testing.T) {
	tests := []struct {
		name     string
		box      *Box
		expected bool
	}{
		{
			name:     "empty box has no content",
			box:      &Box{},
			expected: false,
		},
		{
			name:     "box with title has content",
			box:      &Box{Title: "Test"},
			expected: true,
		},
		{
			name:     "box with subtitle has content",
			box:      &Box{Subtitle: "Sub"},
			expected: true,
		},
		{
			name:     "box with KV pairs has content",
			box:      &Box{KVPairs: []KV{{Key: "k", Value: "v"}}},
			expected: true,
		},
		{
			name:     "box with footer has content",
			box:      &Box{Footer: "Footer"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.box.HasContent())
		})
	}
}
