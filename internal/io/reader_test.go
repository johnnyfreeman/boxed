package io

import (
	"strings"
	"testing"

	"boxed/internal/box"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStdinKVReader(t *testing.T) {
	input := strings.NewReader("")
	reader := NewStdinKVReader(input)

	require.NotNil(t, reader)
	assert.NotNil(t, reader.reader)
}

func TestStdinKVReader_ReadKVPairs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []box.KV
		wantErr bool
	}{
		{
			name:  "single KV pair",
			input: "key=value\n",
			want:  []box.KV{{Key: "key", Value: "value"}},
		},
		{
			name:  "multiple KV pairs",
			input: "env=prod\nregion=us-east-1\nversion=2.0\n",
			want: []box.KV{
				{Key: "env", Value: "prod"},
				{Key: "region", Value: "us-east-1"},
				{Key: "version", Value: "2.0"},
			},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "only empty lines",
			input: "\n\n\n",
			want:  nil,
		},
		{
			name:  "KV with empty value",
			input: "flag=\n",
			want:  []box.KV{{Key: "flag", Value: ""}},
		},
		{
			name:  "KV with multiple equals",
			input: "url=http://example.com?a=1&b=2\n",
			want:  []box.KV{{Key: "url", Value: "http://example.com?a=1&b=2"}},
		},
		{
			name:  "mixed with empty lines",
			input: "key1=value1\n\nkey2=value2\n",
			want: []box.KV{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
		},
		{
			name:  "whitespace handling",
			input: "  key=value  \n",
			want:  []box.KV{{Key: "key", Value: "value"}},
		},
		{
			name:    "invalid format (no equals)",
			input:   "invalid\n",
			wantErr: true,
		},
		{
			name:    "invalid format (empty key)",
			input:   "=value\n",
			wantErr: true,
		},
		{
			name:    "partial valid then invalid",
			input:   "valid=yes\ninvalid\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewStdinKVReader(strings.NewReader(tt.input))

			got, err := reader.ReadKVPairs()

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestStdinKVReader_Unicode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []box.KV
	}{
		{
			name:  "CJK characters",
			input: "環境=本番\n",
			want:  []box.KV{{Key: "環境", Value: "本番"}},
		},
		{
			name:  "emoji",
			input: "status=✅ complete\n",
			want:  []box.KV{{Key: "status", Value: "✅ complete"}},
		},
		{
			name:  "mixed unicode",
			input: "message=Déployé 🚀\n",
			want:  []box.KV{{Key: "message", Value: "Déployé 🚀"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewStdinKVReader(strings.NewReader(tt.input))

			got, err := reader.ReadKVPairs()

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
