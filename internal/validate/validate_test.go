package validate

import (
	"testing"

	"boxed/internal/box"

	"github.com/stretchr/testify/assert"
)

func TestBoxType(t *testing.T) {
	tests := []struct {
		name    string
		boxType string
		wantErr bool
	}{
		{"valid success", "success", false},
		{"valid error", "error", false},
		{"valid info", "info", false},
		{"valid warning", "warning", false},
		{"invalid type", "invalid", true},
		{"empty string", "", true},
		{"wrong case", "Success", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BoxType(tt.boxType)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid box type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestKVPair(t *testing.T) {
	tests := []struct {
		name    string
		kv      string
		wantErr bool
	}{
		{"valid pair", "key=value", false},
		{"valid with spaces", "env name=production", false},
		{"valid empty value", "flag=", false},
		{"valid multiple equals", "url=http://example.com?a=1&b=2", false},
		{"invalid no equals", "keyvalue", true},
		{"invalid empty key", "=value", true},
		{"invalid only equals", "=", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := KVPair(tt.kv)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBorderStyle(t *testing.T) {
	tests := []struct {
		name    string
		style   string
		wantErr bool
	}{
		{"empty (default)", "", false},
		{"normal", "normal", false},
		{"rounded", "rounded", false},
		{"thick", "thick", false},
		{"double", "double", false},
		{"invalid style", "fancy", true},
		{"wrong case", "Rounded", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BorderStyle(tt.style)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid border style")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBox(t *testing.T) {
	tests := []struct {
		name    string
		box     *box.Box
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid box with title",
			box:     &box.Box{Type: box.Success, Title: "Test"},
			wantErr: false,
		},
		{
			name:    "valid box with subtitle",
			box:     &box.Box{Type: box.Success, Subtitle: "Sub"},
			wantErr: false,
		},
		{
			name:    "valid box with KV",
			box:     &box.Box{Type: box.Success, KVPairs: []box.KV{{Key: "k", Value: "v"}}},
			wantErr: false,
		},
		{
			name:    "valid box with footer",
			box:     &box.Box{Type: box.Success, Footer: "Footer"},
			wantErr: false,
		},
		{
			name:    "invalid empty box",
			box:     &box.Box{Type: box.Success},
			wantErr: true,
			errMsg:  "no content",
		},
		{
			name:    "invalid negative width",
			box:     &box.Box{Type: box.Success, Title: "Test", Width: -1},
			wantErr: true,
			errMsg:  "width must be non-negative",
		},
		{
			name:    "valid zero width (auto-size)",
			box:     &box.Box{Type: box.Success, Title: "Test", Width: 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Box(tt.box)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
