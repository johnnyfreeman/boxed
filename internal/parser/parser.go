package parser

import (
	"strings"

	"boxed/internal/box"
	"boxed/internal/validate"
)

// Options represents the parsed CLI flags for box configuration.
// This struct acts as an intermediate representation between Cobra's
// flag parsing and our Box model, keeping CLI concerns separate from
// domain logic.
type Options struct {
	Title       string
	Subtitle    string
	KVFlags     []string
	Footer      string
	Width       int
	BorderStyle string
}

// ParseBox converts CLI arguments into a validated Box model.
// This is a pure function that performs all validation upfront (fail-fast)
// before constructing the box, ensuring that any Box instance that successfully
// returns from this function is guaranteed to be valid and renderable.
func ParseBox(boxType string, opts Options) (*box.Box, error) {
	if err := validate.BoxType(boxType); err != nil {
		return nil, err
	}

	if err := validate.BorderStyle(opts.BorderStyle); err != nil {
		return nil, err
	}

	kvPairs, err := parseKVPairs(opts.KVFlags)
	if err != nil {
		return nil, err
	}

	b := &box.Box{
		Type:        box.BoxType(boxType),
		Title:       opts.Title,
		Subtitle:    opts.Subtitle,
		KVPairs:     kvPairs,
		Footer:      opts.Footer,
		Width:       opts.Width,
		BorderStyle: opts.BorderStyle,
	}

	if err := validate.Box(b); err != nil {
		return nil, err
	}

	return b, nil
}

// parseKVPairs converts an array of "key=value" strings into KV structs.
// Each string is validated before parsing to ensure fail-fast behavior.
// Supports comma-separated pairs (e.g., "A=1,B=2,C=3") for convenience,
// but only splits on commas that appear before a new key=value pattern.
// This allows values to contain commas (e.g., "Status=1 staged, 2 modified").
// Splits on the first '=' only, allowing '=' characters in values
// (e.g., "url=http://example.com?a=1&b=2").
func parseKVPairs(kvFlags []string) ([]box.KV, error) {
	if len(kvFlags) == 0 {
		return nil, nil
	}

	kvPairs := make([]box.KV, 0, len(kvFlags))
	for _, kvFlag := range kvFlags {
		// Smart split: only split on commas followed by key=value pattern
		parts := smartSplitKV(kvFlag)

		for _, kv := range parts {
			kv = strings.TrimSpace(kv)

			if kv == "" {
				continue
			}

			if err := validate.KVPair(kv); err != nil {
				return nil, err
			}

			pair := strings.SplitN(kv, "=", 2)
			kvPairs = append(kvPairs, box.KV{
				Key:   pair[0],
				Value: pair[1],
			})
		}
	}

	return kvPairs, nil
}

// smartSplitKV splits on commas only when followed by a key=value pattern.
// This allows values to contain commas while still supporting comma-separated pairs.
// Examples:
//   - "A=1,B=2,C=3" -> ["A=1", "B=2", "C=3"]
//   - "Status=1 staged, 2 modified" -> ["Status=1 staged, 2 modified"]
//   - "A=1,2,3,B=4" -> ["A=1,2,3", "B=4"]
func smartSplitKV(s string) []string {
	var result []string
	var current strings.Builder

	i := 0
	for i < len(s) {
		if s[i] == ',' {
			// Look ahead to see if this comma is followed by a key=value pattern
			// We check if there's non-whitespace followed by an '=' sign
			j := i + 1
			// Skip whitespace after comma
			for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
				j++
			}

			// Look for the next '=' to see if this starts a new KV pair
			foundEquals := false
			for k := j; k < len(s) && k < j+50; k++ {
				if s[k] == '=' {
					foundEquals = true
					break
				}
				if s[k] == ',' {
					// Hit another comma before equals, so this is just a comma in a value
					break
				}
			}

			if foundEquals && j < len(s) {
				// This comma starts a new KV pair
				result = append(result, current.String())
				current.Reset()
				i++ // Skip the comma
				continue
			}
		}

		current.WriteByte(s[i])
		i++
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
