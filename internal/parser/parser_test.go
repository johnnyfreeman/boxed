package parser

import (
	"testing"

	"boxed/internal/box"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBox(t *testing.T) {
	tests := []struct {
		name    string
		boxType string
		opts    Options
		want    *box.Box
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid success box with title",
			boxType: "success",
			opts:    Options{Title: "Deploy Complete"},
			want: &box.Box{
				Type:  box.Success,
				Title: "Deploy Complete",
			},
			wantErr: false,
		},
		{
			name:    "valid error box with all fields",
			boxType: "error",
			opts: Options{
				Title:       "Build Failed",
				Subtitle:    "v1.2.3",
				KVFlags:     []string{"exit=1", "duration=5m"},
				Footer:      "Check logs",
				Width:       80,
				BorderStyle: "rounded",
			},
			want: &box.Box{
				Type:     box.Error,
				Title:    "Build Failed",
				Subtitle: "v1.2.3",
				KVPairs: []box.KV{
					{Key: "exit", Value: "1"},
					{Key: "duration", Value: "5m"},
				},
				Footer:      "Check logs",
				Width:       80,
				BorderStyle: "rounded",
			},
			wantErr: false,
		},
		{
			name:    "valid box with complex KV values",
			boxType: "info",
			opts: Options{
				Title:   "Test",
				KVFlags: []string{"url=http://example.com?a=1&b=2"},
			},
			want: &box.Box{
				Type:  box.Info,
				Title: "Test",
				KVPairs: []box.KV{
					{Key: "url", Value: "http://example.com?a=1&b=2"},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid box type",
			boxType: "invalid",
			opts:    Options{Title: "Test"},
			wantErr: true,
			errMsg:  "invalid box type",
		},
		{
			name:    "invalid border style",
			boxType: "success",
			opts: Options{
				Title:       "Test",
				BorderStyle: "fancy",
			},
			wantErr: true,
			errMsg:  "invalid border style",
		},
		{
			name:    "invalid KV pair",
			boxType: "success",
			opts: Options{
				Title:   "Test",
				KVFlags: []string{"invalid"},
			},
			wantErr: true,
			errMsg:  "invalid key-value pair",
		},
		{
			name:    "empty box (no content)",
			boxType: "success",
			opts:    Options{},
			wantErr: true,
			errMsg:  "no content",
		},
		{
			name:    "negative width",
			boxType: "success",
			opts: Options{
				Title: "Test",
				Width: -1,
			},
			wantErr: true,
			errMsg:  "width must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBox(tt.boxType, tt.opts)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.want.Type, got.Type)
				assert.Equal(t, tt.want.Title, got.Title)
				assert.Equal(t, tt.want.Subtitle, got.Subtitle)
				assert.Equal(t, tt.want.KVPairs, got.KVPairs)
				assert.Equal(t, tt.want.Footer, got.Footer)
				assert.Equal(t, tt.want.Width, got.Width)
				assert.Equal(t, tt.want.BorderStyle, got.BorderStyle)
			}
		})
	}
}

func TestParseKVPairs(t *testing.T) {
	tests := []struct {
		name    string
		kvFlags []string
		want    []box.KV
		wantErr bool
	}{
		{
			name:    "empty flags",
			kvFlags: nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "single KV pair",
			kvFlags: []string{"key=value"},
			want:    []box.KV{{Key: "key", Value: "value"}},
			wantErr: false,
		},
		{
			name:    "multiple KV pairs",
			kvFlags: []string{"env=prod", "region=us-east-1"},
			want: []box.KV{
				{Key: "env", Value: "prod"},
				{Key: "region", Value: "us-east-1"},
			},
			wantErr: false,
		},
		{
			name:    "empty value",
			kvFlags: []string{"flag="},
			want:    []box.KV{{Key: "flag", Value: ""}},
			wantErr: false,
		},
		{
			name:    "value with equals",
			kvFlags: []string{"url=http://example.com?a=1&b=2"},
			want:    []box.KV{{Key: "url", Value: "http://example.com?a=1&b=2"}},
			wantErr: false,
		},
		{
			name:    "invalid format",
			kvFlags: []string{"invalid"},
			wantErr: true,
		},
		{
			name:    "empty key",
			kvFlags: []string{"=value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKVPairs(tt.kvFlags)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
