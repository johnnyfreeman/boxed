package io

import (
	"encoding/json"
	"fmt"
	"io"

	"boxed/internal/parser"
)

// JSONBox represents the JSON structure for defining a complete box.
// This allows tools that output JSON to easily generate boxed output without
// constructing complex shell command lines.
type JSONBox struct {
	Title       string            `json:"title"`
	Subtitle    string            `json:"subtitle"`
	KV          map[string]string `json:"kv"`
	Footer      string            `json:"footer"`
	Width       int               `json:"width"`
	BorderStyle string            `json:"border_style"`
}

// JSONReader parses box definitions from JSON input.
type JSONReader struct {
	reader io.Reader
}

// NewJSONReader creates a reader that parses JSON box definitions.
func NewJSONReader(r io.Reader) *JSONReader {
	return &JSONReader{reader: r}
}

// ReadBox parses a JSON box definition into parser.Options.
// The JSON format uses a map for KV pairs which is more natural in JSON
// than an array of "key=value" strings.
func (j *JSONReader) ReadBox() (parser.Options, error) {
	var jsonBox JSONBox

	decoder := json.NewDecoder(j.reader)
	if err := decoder.Decode(&jsonBox); err != nil {
		return parser.Options{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	opts := parser.Options{
		Title:       jsonBox.Title,
		Subtitle:    jsonBox.Subtitle,
		Footer:      jsonBox.Footer,
		Width:       jsonBox.Width,
		BorderStyle: jsonBox.BorderStyle,
	}

	for key, value := range jsonBox.KV {
		opts.KVFlags = append(opts.KVFlags, fmt.Sprintf("%s=%s", key, value))
	}

	return opts, nil
}
